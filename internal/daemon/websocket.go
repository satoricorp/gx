package daemon

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/satoricorp/gx/internal/providers"
	"github.com/satoricorp/gx/internal/storage"
)

type websocketFrame struct {
	Type          string `json:"type"`
	Payload       string `json:"payload,omitempty"`
	PayloadBase64 string `json:"payload_base64,omitempty"`
}

type websocketResult struct {
	frames []websocketFrame
	err    error
}

func (p *Proxy) serveWebSocket(w http.ResponseWriter, r *http.Request, sessionID string) {
	provider := detectProvider(r)
	endpoint := normalizeEndpoint(provider, r.URL.Path)
	upstream, ok := resolveUpstream(r, provider)
	if !ok {
		http.Error(w, "unsupported provider", http.StatusBadRequest)
		return
	}

	upstreamURL := upstream.URL("wss", r.URL.RawQuery)
	dialer := websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  30 * time.Second,
		EnableCompression: true,
		TLSClientConfig:   &tls.Config{MinVersion: tls.VersionTLS12},
	}
	upstreamConn, resp, err := dialer.DialContext(r.Context(), upstreamURL, websocketUpstreamHeaders(r.Header))
	if err != nil {
		requestID := uuid.NewString()
		now := time.Now().UnixMilli()
		_ = p.writeRequest(storage.Request{
			ID:             requestID,
			SessionID:      sessionID,
			CreatedAt:      now,
			Provider:       provider,
			Endpoint:       endpoint,
			Method:         r.Method,
			RequestBody:    []byte("[]"),
			RequestHeaders: redactHeaders(r.Header),
		})
		statusCode := http.StatusBadGateway
		responseHeaders := "{}"
		if resp != nil {
			statusCode = resp.StatusCode
			responseHeaders = redactHeaders(resp.Header)
		}
		errText := err.Error()
		_ = p.writeResponse(storage.Response{
			ID:              uuid.NewString(),
			RequestID:       requestID,
			CreatedAt:       now,
			CompletedAt:     time.Now().UnixMilli(),
			StatusCode:      statusCode,
			ResponseBody:    []byte("[]"),
			ResponseHeaders: responseHeaders,
			IsStreaming:     true,
			DurationMS:      time.Now().UnixMilli() - now,
			Error:           &errText,
		})
		http.Error(w, errText, http.StatusBadGateway)
		return
	}
	defer upstreamConn.Close()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer clientConn.Close()

	requestID := uuid.NewString()
	startedAt := time.Now().UnixMilli()
	requestHeaders := redactHeaders(r.Header)
	responseHeaders := "{}"
	var requestIDHeader *string
	if resp != nil {
		responseHeaders = redactHeaders(resp.Header)
		requestIDHeader = providerRequestID(resp.Header)
	}

	reqCh := make(chan websocketResult, 1)
	go func() {
		frames, err := proxyWebSocket(clientConn, upstreamConn)
		_ = upstreamConn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
		reqCh <- websocketResult{frames: frames, err: err}
	}()

	respFrames, respErr := proxyWebSocket(upstreamConn, clientConn)
	_ = clientConn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
	reqResult := <-reqCh

	requestBody, _ := json.Marshal(reqResult.frames)
	responseBody, _ := json.Marshal(respFrames)
	model := extractModelFromWebSocketFrames(reqResult.frames)
	summary := extractSummaryFromWebSocketFrames(respFrames)

	if err := p.writeRequest(storage.Request{
		ID:             requestID,
		SessionID:      sessionID,
		CreatedAt:      startedAt,
		Provider:       provider,
		Endpoint:       endpoint,
		Method:         r.Method,
		Model:          model,
		RequestBody:    requestBody,
		RequestHeaders: requestHeaders,
	}); err != nil {
		log.Println(err)
	}

	completedAt := time.Now().UnixMilli()
	var errText *string
	if err := firstWebSocketError(reqResult.err, respErr); err != nil {
		msg := err.Error()
		errText = &msg
	}
	if err := p.writeResponse(storage.Response{
		ID:                uuid.NewString(),
		RequestID:         requestID,
		CreatedAt:         startedAt,
		CompletedAt:       completedAt,
		StatusCode:        http.StatusSwitchingProtocols,
		ResponseBody:      responseBody,
		ResponseHeaders:   responseHeaders,
		IsStreaming:       true,
		DurationMS:        completedAt - startedAt,
		ProviderRequestID: requestIDHeader,
		InputTokens:       summary.Usage.InputTokens,
		OutputTokens:      summary.Usage.OutputTokens,
		CacheReadTokens:   summary.Usage.CacheReadTokens,
		CacheWriteTokens:  summary.Usage.CacheWriteTokens,
		FinishReason:      summary.FinishReason,
		Error:             errText,
	}); err != nil {
		log.Println(err)
	}
}

func isWebSocketUpgrade(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Upgrade"), "websocket") &&
		strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
}

func websocketUpstreamHeaders(src http.Header) http.Header {
	dst := http.Header{}
	for key, values := range src {
		lower := strings.ToLower(key)
		if lower == "connection" || lower == "upgrade" || strings.HasPrefix(lower, "sec-websocket-") {
			continue
		}
		for _, value := range values {
			dst.Add(key, value)
		}
	}
	return dst
}

func proxyWebSocket(src, dst *websocket.Conn) ([]websocketFrame, error) {
	frames := make([]websocketFrame, 0, 8)
	for {
		messageType, payload, err := src.ReadMessage()
		if err != nil {
			if isBenignWebSocketClosure(err) {
				return frames, nil
			}
			return frames, err
		}
		if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
			continue
		}
		frames = append(frames, loggedWebSocketFrame(messageType, payload))
		if err := dst.WriteMessage(messageType, payload); err != nil {
			return frames, err
		}
	}
}

func isBenignWebSocketClosure(err error) bool {
	if err == nil {
		return true
	}
	if err == io.EOF || strings.Contains(err.Error(), "unexpected EOF") {
		return true
	}
	return websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
}

func loggedWebSocketFrame(messageType int, payload []byte) websocketFrame {
	frame := websocketFrame{Type: websocketMessageType(messageType)}
	if messageType == websocket.TextMessage {
		frame.Payload = string(payload)
		return frame
	}
	frame.PayloadBase64 = base64.StdEncoding.EncodeToString(payload)
	return frame
}

func websocketMessageType(messageType int) string {
	switch messageType {
	case websocket.TextMessage:
		return "text"
	case websocket.BinaryMessage:
		return "binary"
	default:
		return "unknown"
	}
}

func extractModelFromWebSocketFrames(frames []websocketFrame) *string {
	for _, frame := range frames {
		if frame.Type != "text" || frame.Payload == "" {
			continue
		}
		if model := extractModel([]byte(frame.Payload)); model != nil {
			return model
		}
	}
	return nil
}

func extractSummaryFromWebSocketFrames(frames []websocketFrame) providers.Summary {
	var summary providers.Summary
	for _, frame := range frames {
		if frame.Type != "text" || frame.Payload == "" {
			continue
		}

		var payload map[string]any
		if err := json.Unmarshal([]byte(frame.Payload), &payload); err != nil {
			continue
		}

		if response, _ := payload["response"].(map[string]any); response != nil {
			mergeWebSocketSummary(&summary, response)
		}
		mergeWebSocketSummary(&summary, payload)
	}
	return summary
}

func mergeWebSocketSummary(summary *providers.Summary, payload map[string]any) {
	if usage, _ := payload["usage"].(map[string]any); usage != nil {
		if value := usageInt(usage, "input_tokens", "prompt_tokens"); value != nil {
			summary.Usage.InputTokens = value
		}
		if value := usageInt(usage, "output_tokens", "completion_tokens"); value != nil {
			summary.Usage.OutputTokens = value
		}
		if value := usageInt(usage, "cache_read_input_tokens"); value != nil {
			summary.Usage.CacheReadTokens = value
		}
		if value := usageInt(usage, "cache_creation_input_tokens"); value != nil {
			summary.Usage.CacheWriteTokens = value
		}
	}

	if finishReason := payloadString(payload, "finish_reason", "stop_reason"); finishReason != "" {
		summary.FinishReason = &finishReason
	}
}

func usageInt(payload map[string]any, keys ...string) *int {
	for _, key := range keys {
		switch value := payload[key].(type) {
		case float64:
			v := int(value)
			return &v
		case int:
			v := value
			return &v
		}
	}
	return nil
}

func payloadString(payload map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := payload[key].(string); ok {
			return value
		}
	}
	return ""
}

func firstWebSocketError(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}
