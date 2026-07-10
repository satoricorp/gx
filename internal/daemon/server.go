package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"

	cursoringest "github.com/satoricorp/gx/internal/ingest/cursor"
	gxservice "github.com/satoricorp/gx/internal/service"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/version"
)

const idleTimeout = 10 * time.Minute
const ambientIdleTimeout = 10 * time.Minute
const cursorSyncInterval = 30 * time.Second

type Options struct {
	Ambient bool
}

type Server struct {
	store           *storage.Store
	registry        *MemorySessionRegistry
	proxy           *Proxy
	control         *http.Server
	controlLn       net.Listener
	ambient         *http.Server
	ambientLn       net.Listener
	ambientResolver *AmbientSessionResolver

	mu        sync.Mutex
	sessions  map[string]*sessionRuntime
	lastIdle  time.Time
	keepAlive bool
}

type sessionRuntime struct {
	sessionID string
	port      int
	server    *http.Server
	listener  net.Listener
}

type createSessionRequest struct {
	Command   string `json:"command"`
	Cwd       string `json:"cwd"`
	GXVersion string `json:"gx_version"`
}

type endSessionRequest struct {
	ExitCode int `json:"exit_code"`
}

type createSessionResponse struct {
	SessionID string `json:"session_id"`
	Port      int    `json:"port"`
}

func Run(ctx context.Context, options ...Options) error {
	var opts Options
	if len(options) > 0 {
		opts = options[0]
	}
	db, err := storage.Open(ctx)
	if err != nil {
		return err
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		return err
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Println(err)
		}
	}()

	registry := NewMemorySessionRegistry()
	srv := &Server{
		store:     store,
		registry:  registry,
		proxy:     NewProxy(registry, store),
		sessions:  map[string]*sessionRuntime{},
		lastIdle:  time.Now(),
		keepAlive: opts.Ambient,
	}

	if err := srv.serveControl(); err != nil {
		return err
	}
	if opts.Ambient {
		if err := gxservice.NewManager().ApplyEnv(ctx); err != nil {
			return err
		}
		if err := srv.serveAmbientProxy(); err != nil {
			return err
		}
		defer srv.ambientLn.Close()
	}
	defer os.Remove(srv.portFile())
	defer os.Remove(srv.pidFile())

	go backfillUsage(ctx, store)
	go srv.watchIdle()
	if opts.Ambient {
		go srv.watchAmbientIdle()
		go srv.watchCursorIngest(ctx)
	}
	if err := srv.control.Serve(srv.controlLn); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func backfillUsage(ctx context.Context, store *storage.Store) {
	result, err := store.BackfillUsage(ctx)
	if err != nil {
		log.Printf("usage backfill: %v", err)
		return
	}
	if result.ResponsesUpdated > 0 || result.SessionsRefreshed > 0 {
		log.Printf("usage backfill: updated %d responses, refreshed %d sessions", result.ResponsesUpdated, result.SessionsRefreshed)
	}
}

func (s *Server) serveAmbientProxy() error {
	ln, err := net.Listen("tcp", gxservice.ProxyAddress())
	if err != nil {
		return fmt.Errorf("listen ambient proxy: %w", err)
	}
	s.ambientLn = ln
	s.ambientResolver = NewAmbientSessionResolver(s.store)
	s.ambient = &http.Server{Handler: NewProxyWithResolver(s.ambientResolver, s.store)}
	go func() {
		if err := s.ambient.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Println(err)
		}
	}()
	return nil
}

func (s *Server) serveControl() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealthz)
	mux.HandleFunc("/sessions", s.handleCreateSession)
	mux.HandleFunc("/sessions/", s.handleSession)

	ln, err := net.Listen("tcp", gxservice.ControlAddress())
	if err != nil {
		return fmt.Errorf("listen control plane: %w", err)
	}
	s.controlLn = ln
	s.control = &http.Server{Handler: mux}

	if err := os.MkdirAll(filepath.Dir(s.portFile()), 0o755); err != nil {
		return fmt.Errorf("create gx dir: %w", err)
	}
	if err := os.WriteFile(s.portFile(), []byte(strconv.Itoa(ln.Addr().(*net.TCPAddr).Port)), 0o644); err != nil {
		return fmt.Errorf("write daemon port: %w", err)
	}
	if err := os.WriteFile(s.pidFile(), []byte(strconv.Itoa(os.Getpid())), 0o644); err != nil {
		return fmt.Errorf("write daemon pid: %w", err)
	}
	return nil
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req createSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	sessionID := uuid.NewString()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		http.Error(w, "allocate port", http.StatusInternalServerError)
		return
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := s.registry.Register(sessionID, port); err != nil {
		_ = ln.Close()
		http.Error(w, "register session", http.StatusInternalServerError)
		return
	}

	child := &http.Server{Handler: s.proxy}
	runtime := &sessionRuntime{
		sessionID: sessionID,
		port:      port,
		server:    child,
		listener:  ln,
	}
	s.mu.Lock()
	s.sessions[sessionID] = runtime
	s.mu.Unlock()

	if err := s.store.WriteSession(r.Context(), storage.Session{
		ID:        sessionID,
		CreatedAt: time.Now().UnixMilli(),
		Command:   req.Command,
		Cwd:       req.Cwd,
		GXVersion: firstNonEmpty(req.GXVersion, version.Current()),
	}); err != nil {
		_ = s.registry.Release(port)
		_ = ln.Close()
		http.Error(w, "store session", http.StatusInternalServerError)
		return
	}

	go func() {
		if err := child.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	s.markActive()
	writeJSON(w, createSessionResponse{SessionID: sessionID, Port: port})
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || filepath.Base(r.URL.Path) != "end" {
		http.NotFound(w, r)
		return
	}
	sessionID := filepath.Base(filepath.Dir(r.URL.Path))
	var req endSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := s.endSession(r.Context(), sessionID, req.ExitCode); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"ok": true})
}

func (s *Server) endSession(ctx context.Context, sessionID string, exitCode int) error {
	s.mu.Lock()
	runtime := s.sessions[sessionID]
	delete(s.sessions, sessionID)
	s.mu.Unlock()
	if runtime == nil {
		return fmt.Errorf("unknown session %s", sessionID)
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	_ = runtime.server.Shutdown(shutdownCtx)
	_ = runtime.listener.Close()
	_ = s.registry.Release(runtime.port)
	if err := s.store.EndSession(ctx, sessionID, time.Now().UnixMilli(), exitCode); err != nil {
		return err
	}
	if s.registry.Active() == 0 {
		s.markIdle()
	}
	return nil
}

func (s *Server) watchIdle() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		if s.keepAlive {
			continue
		}
		s.mu.Lock()
		idleSince := s.lastIdle
		active := len(s.sessions)
		s.mu.Unlock()
		if active == 0 && !idleSince.IsZero() && time.Since(idleSince) >= idleTimeout {
			_ = s.control.Shutdown(context.Background())
			return
		}
	}
}

func (s *Server) watchAmbientIdle() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		if s.ambientResolver != nil {
			s.ambientResolver.EndIdle(context.Background(), ambientIdleTimeout)
		}
	}
}

func (s *Server) watchCursorIngest(ctx context.Context) {
	s.runCursorIngest(ctx)
	ticker := time.NewTicker(cursorSyncInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.runCursorIngest(ctx)
		}
	}
}

func (s *Server) runCursorIngest(ctx context.Context) {
	if _, err := cursoringest.Sync(ctx, s.store); err != nil {
		log.Printf("cursor ingest: %v", err)
	}
}

func (s *Server) markActive() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastIdle = time.Time{}
}

func (s *Server) markIdle() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastIdle = time.Now()
}

func (s *Server) portFile() string {
	dir, _ := storage.DefaultDir()
	return filepath.Join(dir, "daemon.port")
}

func (s *Server) pidFile() string {
	dir, _ := storage.DefaultDir()
	return filepath.Join(dir, "daemon.pid")
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
