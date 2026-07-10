package daemon

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/satoricorp/gx/internal/inference"
	"github.com/satoricorp/gx/internal/providers"
	"github.com/satoricorp/gx/internal/storage"
)

type Proxy struct {
	resolver   RequestSessionResolver
	store      storage.StorageWriter
	httpClient *http.Client
}

func NewProxy(registry SessionRegistry, store storage.StorageWriter) *Proxy {
	return NewProxyWithResolver(NewPortSessionResolver(registry), store)
}

func NewProxyWithResolver(resolver RequestSessionResolver, store storage.StorageWriter) *Proxy {
	return &Proxy{
		resolver: resolver,
		store:    store,
		httpClient: &http.Client{
			Transport: &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				ForceAttemptHTTP2:     true,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 30 * time.Second,
				TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12},
			},
		},
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := p.resolver.ResolveRequestSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if isWebSocketUpgrade(r) {
		if !session.Capture {
			http.Error(w, "websocket capture requires a known session", http.StatusBadGateway)
			return
		}
		p.serveWebSocket(w, r, session.ID)
		return
	}

	provider := detectProvider(r)
	endpoint := normalizeEndpoint(provider, r.URL.Path)
	upstream, ok := resolveUpstream(r, provider)
	if !ok {
		http.Error(w, "unsupported provider", http.StatusBadRequest)
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read request body", http.StatusBadRequest)
		return
	}
	_ = r.Body.Close()

	requestID := uuid.NewString()
	now := time.Now().UnixMilli()
	model := extractModel(requestBody)
	if session.Capture {
		if err := p.writeRequest(storage.Request{
			ID:             requestID,
			SessionID:      session.ID,
			CreatedAt:      now,
			Provider:       provider,
			Endpoint:       endpoint,
			Method:         r.Method,
			Model:          model,
			RequestBody:    requestBody,
			RequestHeaders: redactHeaders(r.Header),
		}); err != nil {
			http.Error(w, "store request", http.StatusInternalServerError)
			return
		}
	}

	upstreamURL := upstream.URL("https", r.URL.RawQuery)
	req, err := http.NewRequestWithContext(r.Context(), r.Method, upstreamURL, bytes.NewReader(requestBody))
	if err != nil {
		http.Error(w, "build upstream request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header.Clone()
	req.Host = upstream.host
	applyInferenceAuth(req, provider)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		errMsg := err.Error()
		if session.Capture {
			p.writeErrorResponse(requestID, now, &errMsg)
		}
		http.Error(w, errMsg, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	responseStarted := time.Now().UnixMilli()

	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	responseBody, streaming, err := copyResponse(w, resp)
	if err != nil {
		errMsg := err.Error()
		if session.Capture {
			p.writeErrorResponse(requestID, now, &errMsg)
		}
		return
	}
	if !session.Capture {
		return
	}

	storedBody := responseBody
	summary, sumErr := providers.SummarizeJSON(provider, responseBody)
	if sumErr != nil {
		summary = providers.Summary{}
	}
	if streaming || (sumErr != nil && looksLikeSSE(responseBody)) {
		if handler := providers.HandlerFor(provider, endpoint); handler != nil {
			if assembled, assembledSummary, err := handler.Assemble(bytes.NewReader(responseBody)); err == nil {
				storedBody = assembled
				summary = assembledSummary
			}
		}
	}

	requestIDHeader := providerRequestID(resp.Header)
	var errText *string
	if resp.StatusCode >= 400 {
		msg := string(responseBody)
		errText = &msg
	}

	completedAt := time.Now().UnixMilli()
	if err := p.writeResponse(storage.Response{
		ID:                uuid.NewString(),
		RequestID:         requestID,
		CreatedAt:         responseStarted,
		CompletedAt:       completedAt,
		StatusCode:        resp.StatusCode,
		ResponseBody:      storedBody,
		ResponseHeaders:   redactHeaders(resp.Header),
		IsStreaming:       streaming,
		DurationMS:        completedAt - responseStarted,
		ProviderRequestID: requestIDHeader,
		FinishReason:      summary.FinishReason,
		InputTokens:       summary.Usage.InputTokens,
		OutputTokens:      summary.Usage.OutputTokens,
		CacheReadTokens:   summary.Usage.CacheReadTokens,
		CacheWriteTokens:  summary.Usage.CacheWriteTokens,
		Error:             errText,
	}); err != nil {
		log.Println(err)
	}
}

func looksLikeSSE(body []byte) bool {
	body = bytes.TrimSpace(body)
	return bytes.HasPrefix(body, []byte("data:")) || bytes.Contains(body, []byte("\ndata:"))
}

func (p *Proxy) sessionForRequest(r *http.Request) (string, bool) {
	session, err := p.resolver.ResolveRequestSession(r)
	return session.ID, err == nil && session.Capture
}

func (p *Proxy) writeErrorResponse(requestID string, startedAt int64, errText *string) {
	completedAt := time.Now().UnixMilli()
	_ = p.writeResponse(storage.Response{
		ID:              uuid.NewString(),
		RequestID:       requestID,
		CreatedAt:       startedAt,
		CompletedAt:     completedAt,
		StatusCode:      http.StatusBadGateway,
		ResponseBody:    []byte{},
		ResponseHeaders: "{}",
		IsStreaming:     false,
		DurationMS:      completedAt - startedAt,
		Error:           errText,
	})
}

func (p *Proxy) writeRequest(req storage.Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.store.WriteRequest(ctx, req)
}

func (p *Proxy) writeResponse(resp storage.Response) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.store.WriteResponse(ctx, resp)
}

func detectProvider(r *http.Request) string {
	switch {
	case r.Header.Get("x-api-key") != "", r.Header.Get("anthropic-api-key") != "":
		return "anthropic"
	case strings.HasPrefix(strings.ToLower(r.Header.Get("authorization")), "bearer "):
		return "openai"
	case strings.HasPrefix(r.URL.Path, "/v1/messages"):
		return "anthropic"
	default:
		return "openai"
	}
}

func applyInferenceAuth(req *http.Request, provider string) {
	creds, ok := inference.Resolve()
	if !ok || creds.Provider != provider {
		return
	}
	switch provider {
	case "anthropic":
		if req.Header.Get("x-api-key") != "" || req.Header.Get("anthropic-api-key") != "" {
			return
		}
		req.Header.Set("x-api-key", creds.APIKey)
	case "openai":
		if req.Header.Get("Authorization") != "" || isChatGPTAuth(req) {
			return
		}
		req.Header.Set("Authorization", "Bearer "+creds.APIKey)
	}
}

func normalizeEndpoint(provider, path string) string {
	switch provider {
	case "anthropic":
		if strings.HasPrefix(path, "/v1/messages") {
			return "/v1/messages"
		}
	case "openai":
		switch {
		case strings.HasPrefix(path, "/chat/completions"), strings.HasPrefix(path, "/v1/chat/completions"):
			return "/v1/chat/completions"
		case strings.HasPrefix(path, "/responses"), strings.HasPrefix(path, "/v1/responses"):
			return "/v1/responses"
		case strings.HasPrefix(path, "/models"), strings.HasPrefix(path, "/v1/models"):
			return "/v1/models"
		}
	}
	return path
}

func upstreamPath(provider, path string) string {
	if provider != "openai" {
		return path
	}
	switch {
	case strings.HasPrefix(path, "/v1/"):
		return path
	case strings.HasPrefix(path, "/chat/completions"), strings.HasPrefix(path, "/responses"), strings.HasPrefix(path, "/models"):
		return "/v1" + path
	default:
		return path
	}
}

func providerRequestID(headers http.Header) *string {
	for _, key := range []string{"x-request-id", "request-id"} {
		if value := headers.Get(key); value != "" {
			return &value
		}
	}
	return nil
}

func extractModel(body []byte) *string {
	type payload struct {
		Model string `json:"model"`
	}
	var p payload
	if err := json.Unmarshal(body, &p); err != nil || p.Model == "" {
		return nil
	}
	return &p.Model
}

func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func copyResponse(w http.ResponseWriter, resp *http.Response) ([]byte, bool, error) {
	streaming := strings.Contains(strings.ToLower(resp.Header.Get("content-type")), "text/event-stream")
	var buf bytes.Buffer
	if !streaming {
		_, err := io.Copy(io.MultiWriter(w, &buf), resp.Body)
		return buf.Bytes(), false, err
	}
	writer := &flushWriter{dst: w, buf: &buf}
	_, err := io.Copy(writer, resp.Body)
	return buf.Bytes(), true, err
}

type flushWriter struct {
	dst http.ResponseWriter
	buf *bytes.Buffer
}

func (w *flushWriter) Write(p []byte) (int, error) {
	if _, err := w.buf.Write(p); err != nil {
		return 0, err
	}
	n, err := w.dst.Write(p)
	if flusher, ok := w.dst.(http.Flusher); ok {
		flusher.Flush()
	}
	return n, err
}
