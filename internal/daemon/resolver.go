package daemon

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/version"
)

type RequestSession struct {
	ID      string
	Capture bool
}

type RequestSessionResolver interface {
	ResolveRequestSession(r *http.Request) (RequestSession, error)
	Active() int
}

type PortSessionResolver struct {
	registry SessionRegistry
}

func NewPortSessionResolver(registry SessionRegistry) PortSessionResolver {
	return PortSessionResolver{registry: registry}
}

func (r PortSessionResolver) ResolveRequestSession(req *http.Request) (RequestSession, error) {
	port, ok := localPort(req)
	if !ok {
		return RequestSession{}, fmt.Errorf("unknown session")
	}
	sessionID, ok := r.registry.Resolve(port)
	if !ok {
		return RequestSession{}, fmt.Errorf("unknown session")
	}
	return RequestSession{ID: sessionID, Capture: true}, nil
}

func (r PortSessionResolver) Active() int {
	if active, ok := r.registry.(interface{ Active() int }); ok {
		return active.Active()
	}
	return 0
}

type AmbientSessionResolver struct {
	store          *storage.Store
	resolveProcess func(*http.Request) ClientProcess

	mu       sync.Mutex
	sessions map[string]*ambientSession
}

type ambientSession struct {
	ID       string
	LastSeen time.Time
}

type ClientProcess struct {
	PID     int
	PPID    int
	Name    string
	Command string
	Cwd     string
}

func NewAmbientSessionResolver(store *storage.Store) *AmbientSessionResolver {
	return &AmbientSessionResolver{
		store:          store,
		resolveProcess: resolveClientProcess,
		sessions:       map[string]*ambientSession{},
	}
}

func (r *AmbientSessionResolver) ResolveRequestSession(req *http.Request) (RequestSession, error) {
	proc := r.resolveProcess(req)
	if !isKnownAgent(proc.Name, proc.Command) {
		return RequestSession{Capture: false}, nil
	}
	key := ambientKey(proc, req)
	now := time.Now()
	cwd := firstNonEmpty(proc.Cwd, ".")
	repoRoot := detectRepoRoot(cwd)

	r.mu.Lock()
	if existing := r.sessions[key]; existing != nil {
		existing.LastSeen = now
		r.mu.Unlock()
		_ = r.store.TouchSession(context.Background(), existing.ID, now.UnixMilli())
		return RequestSession{ID: existing.ID, Capture: true}, nil
	}
	sessionID := uuid.NewString()
	r.sessions[key] = &ambientSession{ID: sessionID, LastSeen: now}
	r.mu.Unlock()

	source := "ambient"
	lastSeen := now.UnixMilli()
	var pid *int
	var ppid *int
	if proc.PID > 0 {
		pid = &proc.PID
	}
	if proc.PPID > 0 {
		ppid = &proc.PPID
	}
	processName := emptyStringNil(proc.Name)
	if err := r.store.WriteSession(context.Background(), storage.Session{
		ID:          sessionID,
		CreatedAt:   lastSeen,
		Command:     firstNonEmpty(proc.Command, proc.Name, "ambient"),
		Cwd:         cwd,
		ClientPID:   pid,
		GXVersion:   version.Current(),
		Source:      &source,
		ProcessName: processName,
		ParentPID:   ppid,
		LastSeenAt:  &lastSeen,
		RepoRoot:    emptyStringNil(repoRoot),
	}); err != nil {
		r.mu.Lock()
		delete(r.sessions, key)
		r.mu.Unlock()
		return RequestSession{}, err
	}
	return RequestSession{ID: sessionID, Capture: true}, nil
}

func (r *AmbientSessionResolver) Active() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.sessions)
}

func (r *AmbientSessionResolver) EndIdle(ctx context.Context, idleFor time.Duration) {
	cutoff := time.Now().Add(-idleFor)
	var ended []string
	r.mu.Lock()
	for key, session := range r.sessions {
		if session.LastSeen.Before(cutoff) {
			ended = append(ended, session.ID)
			delete(r.sessions, key)
		}
	}
	r.mu.Unlock()
	for _, sessionID := range ended {
		_ = r.store.EndSession(ctx, sessionID, time.Now().UnixMilli(), 0)
	}
}

func resolveClientProcess(req *http.Request) ClientProcess {
	local, localOK := localPort(req)
	remote, remoteOK := remotePort(req)
	if !localOK || !remoteOK {
		return ClientProcess{}
	}
	out, err := exec.Command("lsof", "-nP", "-iTCP", "-sTCP:ESTABLISHED").Output()
	if err != nil {
		return ClientProcess{}
	}
	needle := fmt.Sprintf("127.0.0.1:%d->127.0.0.1:%d", remote, local)
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, needle) {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		pid, _ := strconv.Atoi(fields[1])
		proc := ClientProcess{PID: pid, Name: fields[0]}
		if pid > 0 {
			proc.Command = psValue(pid, "command")
			proc.PPID, _ = strconv.Atoi(strings.TrimSpace(psValue(pid, "ppid")))
			proc.Cwd = processCwd(pid)
		}
		return proc
	}
	return ClientProcess{}
}

func localPort(req *http.Request) (int, bool) {
	addr, ok := req.Context().Value(http.LocalAddrContextKey).(net.Addr)
	if !ok || addr == nil {
		return 0, false
	}
	tcp, ok := addr.(*net.TCPAddr)
	if !ok {
		return 0, false
	}
	return tcp.Port, true
}

func remotePort(req *http.Request) (int, bool) {
	host, portText, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil || (host != "127.0.0.1" && host != "::1" && host != "localhost") {
		return 0, false
	}
	port, err := strconv.Atoi(portText)
	return port, err == nil
}

func psValue(pid int, field string) string {
	out, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", field+"=").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func processCwd(pid int) string {
	out, err := exec.Command("lsof", "-a", "-p", strconv.Itoa(pid), "-d", "cwd", "-Fn").Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "n") {
			return strings.TrimPrefix(line, "n")
		}
	}
	return ""
}

func isKnownAgent(name, command string) bool {
	haystack := strings.ToLower(name + " " + command)
	for _, agent := range []string{"codex", "claude"} {
		if strings.Contains(haystack, agent) {
			return true
		}
	}
	return false
}

func ambientKey(proc ClientProcess, req *http.Request) string {
	if proc.PID > 0 {
		return "pid:" + strconv.Itoa(proc.PID)
	}
	return "remote:" + req.RemoteAddr
}

func detectRepoRoot(cwd string) string {
	for dir := cwd; dir != "" && dir != string(filepath.Separator); dir = filepath.Dir(dir) {
		if exists(filepath.Join(dir, ".git")) || exists(filepath.Join(dir, ".jj")) {
			return dir
		}
	}
	return ""
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func emptyStringNil(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
