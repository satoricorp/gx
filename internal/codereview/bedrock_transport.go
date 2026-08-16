package codereview

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/cloud"
)

// bedrockTransport is how a review leg's bytes reach bedrock-runtime.
//
// Both legs and the judge share one request shape and one response parser (see
// bedrockAnthropicReviewer.completeJSON); only the wire differs. Keeping that
// difference behind this interface is what lets the same reviewer run against a
// developer's own AWS credentials and against gx Cloud without a second
// implementation of prompt construction, parsing, or failure reporting.
type bedrockTransport interface {
	// complete performs one model call and returns the model's reply.
	complete(ctx context.Context, model, system, input string, maxOutputTokens int) (bedrockCompletion, error)
	// detail is the human phrase naming this wire, reported on every review so
	// that "slow" and "denied" stay answerable after the fact.
	detail() string
}

// bedrockCompletion is one model reply: the text, and why the model stopped.
//
// The stop reason is carried because a reply cut off at the output cap and a
// reply the model formatted badly are the same symptom at the parse site — JSON
// that will not decode — and completely different problems. Both wires report
// it, so the distinction survives whichever one a user is on.
type bedrockCompletion struct {
	Text string
	// StopReason is Bedrock's stop_reason: "end_turn", "max_tokens", and so on.
	// Empty when the reply carried none.
	StopReason string
}

// truncated reports whether the model ran out of output budget mid-reply.
func (c bedrockCompletion) truncated() bool {
	return strings.EqualFold(strings.TrimSpace(c.StopReason), "max_tokens")
}

const (
	bedrockTransportKindDirect = "direct"
	bedrockTransportKindCloud  = "cloud"

	// bedrockAnthropicVersion is the Anthropic API version bedrock-runtime
	// expects in an InvokeModel body. gx Cloud accepts and ignores it.
	bedrockAnthropicVersion = "bedrock-2023-05-31"

	// bedrockDirectTimeout bounds one direct model call. Flagship models
	// reading a large brief routinely take minutes.
	bedrockDirectTimeout = 5 * time.Minute
)

// bedrockTransportPlan is the resolved answer to "how will this review reach
// Bedrock?", decided once so that the reviewer legs and the judge cannot end up
// on different wires with no way to tell from the output.
type bedrockTransportPlan struct {
	Kind  string
	creds bedrockCredentials
	// client is set for the cloud kind. It is resolved during planning so a
	// missing login is an error at panel-construction time, where it can be
	// reported, rather than a 401 per leg mid-review.
	client   *cloud.Client
	cloudURL string
}

// newTransport builds one transport for one leg. Each leg gets its own so the
// two reviewers and the judge do not share an http.Client's connection budget.
func (p bedrockTransportPlan) newTransport() bedrockTransport {
	if p.Kind == bedrockTransportKindCloud {
		return &cloudBedrockTransport{client: p.client, url: p.cloudURL}
	}
	return newDirectBedrockTransport(p.creds)
}

// bedrockDirectEnvVar opts a caller out of gx Cloud and onto their own AWS
// account. It must be set deliberately; ambient AWS credentials never select it.
const bedrockDirectEnvVar = "GX_REVIEW_BEDROCK_DIRECT"

// resolveBedrockTransportPlan picks the wire: gx Cloud, unless the caller has
// explicitly asked for their own AWS account.
//
// Cloud is THE path, not a fallback. The previous order — local AWS credentials
// win if present — meant the wire was chosen by whatever happened to be exported
// in the caller's shell. Anyone with `AWS_ACCESS_KEY_ID` set for unrelated work
// (most developers, most CI) silently reviewed on their own account and quota
// while believing they were using the product, and the only way to reach Cloud
// was to unset credentials they needed for something else. Worse, the two wires
// fail differently, so "it works on my machine" tracked the developer's ambient
// environment rather than anything about gx.
//
// Direct AWS is still reachable, but only by asking:
// GX_REVIEW_BEDROCK_DIRECT=1 plus the usual AWS_* variables. That is for
// working on the reviewer itself — benchmarking models, testing a region — not
// for ordinary review.
//
// Errors here are returned, never swallowed: a review with no reviewer must
// read as "nothing reviewed this", never as a clean review.
func resolveBedrockTransportPlan() (bedrockTransportPlan, error) {
	if bedrockDirectRequested() {
		creds, credsErr := bedrockCredentialsFromEnv()
		if credsErr != nil {
			// They asked for direct explicitly, so do not quietly reroute to
			// Cloud — that would hide the misconfiguration they need to see.
			return bedrockTransportPlan{}, fmt.Errorf("%s is set, so gx is using your own AWS account, but %v", bedrockDirectEnvVar, credsErr)
		}
		return bedrockTransportPlan{Kind: bedrockTransportKindDirect, creds: creds}, nil
	}
	client := cloud.NewBedrockClient()
	if client == nil {
		return bedrockTransportPlan{}, fmt.Errorf("gx Cloud is not configured for this build (no cloud URL), so there is no reviewer; set %s=1 with AWS_* credentials to review on your own AWS account instead", bedrockDirectEnvVar)
	}
	if _, err := cloud.CloudAPIToken(); err != nil {
		return bedrockTransportPlan{}, fmt.Errorf("not signed in to gx Cloud: run `gx auth login` (%v)", err)
	}
	return bedrockTransportPlan{Kind: bedrockTransportKindCloud, client: client, cloudURL: cloud.CloudURL()}, nil
}

func bedrockDirectRequested() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(bedrockDirectEnvVar))) {
	case "1", "true", "yes", "on":
		return true
	}
	return false
}

// bedrockTransportShortName is the label fragment for one wire, short enough to
// sit inside a finding's reviewer attribution.
func bedrockTransportShortName(kind string) string {
	switch kind {
	case bedrockTransportKindDirect:
		return "AWS"
	case bedrockTransportKindCloud:
		return "gx Cloud"
	case "openai":
		return "OpenAI"
	default:
		return ""
	}
}

// bedrockRequestBody is the Anthropic messages shape bedrock-runtime's
// InvokeModel takes, and the shape gx Cloud's /gx/bedrock/fight normalizes.
type bedrockRequestBody struct {
	AnthropicVersion string           `json:"anthropic_version"`
	MaxTokens        int              `json:"max_tokens"`
	System           string           `json:"system,omitempty"`
	Messages         []bedrockMessage `json:"messages"`
}

type bedrockMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type bedrockResponseBody struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
}

// text joins the reply's text blocks. A model that answers in several blocks
// has still answered once.
func (r bedrockResponseBody) text() string {
	var b strings.Builder
	for _, block := range r.Content {
		if block.Type != "" && block.Type != "text" {
			continue
		}
		b.WriteString(block.Text)
	}
	return strings.TrimSpace(b.String())
}

func bedrockRequestPayload(system, input string, maxOutputTokens int) bedrockRequestBody {
	if maxOutputTokens <= 0 {
		maxOutputTokens = defaultReviewMaxOutputTokens
	}
	return bedrockRequestBody{
		AnthropicVersion: bedrockAnthropicVersion,
		MaxTokens:        maxOutputTokens,
		System:           system,
		Messages:         []bedrockMessage{{Role: "user", Content: input}},
	}
}

// bedrockModelSpeaksAnthropic reports whether a model takes the Anthropic
// messages body on InvokeModel. Anthropic model IDs — bare, profile-prefixed,
// or inside an ARN — all carry "anthropic."; every other vendor on Bedrock
// (Amazon Nova, Meta Llama, Mistral, Writer, ...) has its own native body, and
// the one shape they all share is the Converse API.
func bedrockModelSpeaksAnthropic(model string) bool {
	return strings.Contains(model, "anthropic.")
}

// bedrockActionForModel is the bedrock-runtime endpoint suffix for one model:
// Anthropic models keep the InvokeModel wire gx has always spoken; everyone
// else goes through Converse, whose request and reply shapes are model-agnostic.
func bedrockActionForModel(model string) string {
	if bedrockModelSpeaksAnthropic(model) {
		return "invoke"
	}
	return "converse"
}

// converseRequestBody is the Converse API's model-agnostic request shape. It
// exists so a reviewer leg can be a non-Anthropic model without gx learning
// each vendor's native InvokeModel dialect.
type converseRequestBody struct {
	System          []converseText          `json:"system,omitempty"`
	Messages        []converseMessage       `json:"messages"`
	InferenceConfig converseInferenceConfig `json:"inferenceConfig"`
}

type converseText struct {
	Text string `json:"text"`
}

type converseMessage struct {
	Role    string         `json:"role"`
	Content []converseText `json:"content"`
}

type converseInferenceConfig struct {
	MaxTokens int `json:"maxTokens"`
}

// converseResponseBody is the Converse reply. Content blocks other than text —
// a reasoning model's reasoningContent, for instance — decode with an empty
// Text and are skipped by text(), which is the behaviour we want: the leg's
// answer is its answer, not its scratch work.
type converseResponseBody struct {
	Output struct {
		Message struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"message"`
	} `json:"output"`
	StopReason string `json:"stopReason"`
}

func (r converseResponseBody) text() string {
	var b strings.Builder
	for _, block := range r.Output.Message.Content {
		b.WriteString(block.Text)
	}
	return strings.TrimSpace(b.String())
}

func converseRequestPayload(system, input string, maxOutputTokens int) converseRequestBody {
	if maxOutputTokens <= 0 {
		maxOutputTokens = defaultReviewMaxOutputTokens
	}
	body := converseRequestBody{
		Messages:        []converseMessage{{Role: "user", Content: []converseText{{Text: input}}}},
		InferenceConfig: converseInferenceConfig{MaxTokens: maxOutputTokens},
	}
	if strings.TrimSpace(system) != "" {
		body.System = []converseText{{Text: system}}
	}
	return body
}

// directBedrockTransport signs and posts InvokeModel itself.
//
// gx depends on no AWS SDK, so the SigV4 signing below is the whole of it. That
// is a deliberate trade: the signing is thirty lines and fully covered by the
// canned-response tests, against an SDK that would pull in dozens of modules
// for one endpoint.
type directBedrockTransport struct {
	region       string
	accessKey    string
	secretKey    string
	sessionToken string
	client       *http.Client
}

func newDirectBedrockTransport(creds bedrockCredentials) *directBedrockTransport {
	region := strings.TrimSpace(creds.region)
	if region == "" {
		region = defaultBedrockRegion
	}
	return &directBedrockTransport{
		region:       region,
		accessKey:    creds.accessKey,
		secretKey:    creds.secretKey,
		sessionToken: creds.sessionToken,
		client:       &http.Client{Timeout: bedrockDirectTimeout},
	}
}

func (t *directBedrockTransport) detail() string {
	return "direct AWS credentials (" + t.region + ")"
}

// bedrockRetryBackoffs paces retries of transient bedrock-runtime failures.
// A model call is minutes of work; waiting seconds to save one is always the
// right trade. A var rather than a const so tests can shrink the waits.
var bedrockRetryBackoffs = []time.Duration{2 * time.Second, 8 * time.Second, 20 * time.Second}

func (t *directBedrockTransport) complete(ctx context.Context, model, system, input string, maxOutputTokens int) (bedrockCompletion, error) {
	if t == nil {
		return bedrockCompletion{}, fmt.Errorf("Bedrock reviewer has no transport")
	}
	var payload any = bedrockRequestPayload(system, input, maxOutputTokens)
	if !bedrockModelSpeaksAnthropic(model) {
		payload = converseRequestPayload(system, input, maxOutputTokens)
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return bedrockCompletion{}, fmt.Errorf("marshal Bedrock request: %w", err)
	}
	// Throttles (429) and server errors (5xx) are transient: without a retry,
	// one blip permanently costs a shard's leg — every file in that shard goes
	// unread by that model and the review degrades. Each attempt re-signs, so
	// X-Amz-Date stays fresh.
	var lastErr error
	for attempt := 0; ; attempt++ {
		completion, retryable, attemptErr := t.completeOnce(ctx, model, body)
		if attemptErr == nil {
			return completion, nil
		}
		lastErr = attemptErr
		if !retryable || attempt >= len(bedrockRetryBackoffs) {
			return bedrockCompletion{}, lastErr
		}
		select {
		case <-ctx.Done():
			return bedrockCompletion{}, lastErr
		case <-time.After(bedrockRetryBackoffs[attempt]):
		}
	}
}

// completeOnce is a single signed InvokeModel attempt. retryable reports
// whether the failure is worth another attempt: throttles and server-side
// errors are; validation errors, access denials, and parse failures are not.
func (t *directBedrockTransport) completeOnce(ctx context.Context, model string, body []byte) (bedrockCompletion, bool, error) {
	req, err := t.newSignedRequest(ctx, model, body)
	if err != nil {
		return bedrockCompletion{}, false, err
	}
	client := t.client
	if client == nil {
		client = &http.Client{Timeout: bedrockDirectTimeout}
	}
	resp, err := client.Do(req)
	if err != nil {
		// Transport-level failures (connection reset, timeout) are as
		// transient as a 503, unless the context itself is done.
		return bedrockCompletion{}, ctx.Err() == nil, fmt.Errorf("call Bedrock model %s in %s: %w", model, t.region, err)
	}
	defer resp.Body.Close()
	raw, readErr := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		retryable := resp.StatusCode == http.StatusTooManyRequests ||
			resp.StatusCode >= 500 ||
			strings.Contains(strings.ToLower(string(raw)), "throttl")
		// describeBedrockFailure turns the exception name and prose into the
		// knob to turn; see its comment in ai.go.
		return bedrockCompletion{}, retryable, describeBedrockFailure(resp.StatusCode, resp.Status, raw, model, t.region)
	}
	if readErr != nil {
		return bedrockCompletion{}, true, fmt.Errorf("read Bedrock response for %s: %w", model, readErr)
	}
	var text, stopReason string
	if bedrockModelSpeaksAnthropic(model) {
		var decoded bedrockResponseBody
		if err := json.Unmarshal(raw, &decoded); err != nil {
			return bedrockCompletion{}, false, fmt.Errorf("decode Bedrock response for %s: %w", model, err)
		}
		text, stopReason = decoded.text(), decoded.StopReason
	} else {
		var decoded converseResponseBody
		if err := json.Unmarshal(raw, &decoded); err != nil {
			return bedrockCompletion{}, false, fmt.Errorf("decode Bedrock Converse response for %s: %w", model, err)
		}
		text, stopReason = decoded.text(), decoded.StopReason
	}
	if text == "" {
		return bedrockCompletion{}, false, fmt.Errorf("Bedrock model %s returned no content", model)
	}
	return bedrockCompletion{Text: text, StopReason: stopReason}, false, nil
}

func (t *directBedrockTransport) newSignedRequest(ctx context.Context, model string, body []byte) (*http.Request, error) {
	model = strings.TrimSpace(model)
	if model == "" {
		return nil, fmt.Errorf("Bedrock request has no model")
	}
	host := "bedrock-runtime." + t.region + ".amazonaws.com"
	action := bedrockActionForModel(model)
	// The request line carries the path escaped once; the canonical request
	// carries it escaped twice. See bedrockCanonicalPath for why they differ.
	rawPath := bedrockRequestPath(model, action)
	endpoint := &url.URL{
		Scheme:  "https",
		Host:    host,
		Path:    "/model/" + model + "/" + action,
		RawPath: rawPath,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create Bedrock request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if err := t.signV4(req, host, bedrockCanonicalPath(model, action), body, time.Now().UTC()); err != nil {
		return nil, err
	}
	return req, nil
}

// bedrockRequestPath is the path that goes on the wire: each path segment
// percent-encoded once, which is what an HTTP request line requires.
func bedrockRequestPath(model, action string) string {
	return "/model/" + awsEscapePathSegment(model) + "/" + action
}

// bedrockCanonicalPath is the path that goes into the SigV4 canonical request,
// which is the wire path percent-encoded a SECOND time.
//
// This is not symmetry-breaking for its own sake, it is the SigV4 rule: every
// service except S3 requires each path segment to be URI-encoded twice when
// building the canonical request, while the request line is encoded once. For
// most model IDs the two are identical, because nothing in them needs escaping —
// which is exactly why signing both the same way looked correct. The A leg's
// ID is one of those. The B leg's inference profile ID ends in ":0", so its wire
// path carries "%3A" and its canonical path must carry "%253A"; signing the wire
// path made every single B call fail with SignatureDoesNotMatch.
//
// Measured, not inferred: instrumenting multiAIReviewer across a three-shard
// live review showed bedrock-a returning 5 findings per shard and bedrock-b
// returning 0 with this error on all three. The panel was one model wearing two
// names, and it was silent because ReviewForSummary discards a leg's error
// whenever the other leg answers — see legFailureLog in ai.go for that half.
//
// AWS itself supplied the expected canonical string in that error response, and
// TestBedrockCanonicalRequestMatchesAWSExpectation pins this function's output
// against it verbatim.
func bedrockCanonicalPath(model, action string) string {
	return "/model/" + awsEscapePathSegment(awsEscapePathSegment(model)) + "/" + action
}

// signV4 applies AWS Signature Version 4 to a bedrock-runtime request.
func (t *directBedrockTransport) signV4(req *http.Request, host, canonicalURI string, body []byte, now time.Time) error {
	if strings.TrimSpace(t.accessKey) == "" || strings.TrimSpace(t.secretKey) == "" {
		return fmt.Errorf("Bedrock reviewer has no AWS credentials: set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY")
	}
	const service = "bedrock"
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")
	payloadHash := hexSHA256(body)

	req.Header.Set("X-Amz-Date", amzDate)
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)
	if strings.TrimSpace(t.sessionToken) != "" {
		req.Header.Set("X-Amz-Security-Token", t.sessionToken)
	}

	signed := []string{"content-type", "host", "x-amz-content-sha256", "x-amz-date"}
	values := map[string]string{
		"content-type":         req.Header.Get("Content-Type"),
		"host":                 host,
		"x-amz-content-sha256": payloadHash,
		"x-amz-date":           amzDate,
	}
	if token := req.Header.Get("X-Amz-Security-Token"); token != "" {
		signed = append(signed, "x-amz-security-token")
		values["x-amz-security-token"] = token
	}

	signedHeaders := strings.Join(signed, ";")
	canonicalRequest := bedrockCanonicalRequest(req.Method, canonicalURI, signed, values, payloadHash)

	scope := strings.Join([]string{dateStamp, t.region, service, "aws4_request"}, "/")
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		scope,
		hexSHA256([]byte(canonicalRequest)),
	}, "\n")

	key := hmacSHA256([]byte("AWS4"+t.secretKey), dateStamp)
	key = hmacSHA256(key, t.region)
	key = hmacSHA256(key, service)
	key = hmacSHA256(key, "aws4_request")
	signature := hex.EncodeToString(hmacSHA256(key, stringToSign))

	req.Header.Set("Authorization", fmt.Sprintf(
		"AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		t.accessKey, scope, signedHeaders, signature,
	))
	return nil
}

// bedrockCanonicalRequest assembles SigV4's "task 1" string.
//
// It is a separate function only so a test can compare it, character for
// character, with the canonical string AWS reports back in a
// SignatureDoesNotMatch body. That comparison is the only external oracle this
// hand-rolled signer has.
func bedrockCanonicalRequest(method, canonicalURI string, signedHeaders []string, values map[string]string, payloadHash string) string {
	var canonicalHeaders strings.Builder
	for _, name := range signedHeaders {
		canonicalHeaders.WriteString(name)
		canonicalHeaders.WriteString(":")
		canonicalHeaders.WriteString(strings.TrimSpace(values[name]))
		canonicalHeaders.WriteString("\n")
	}
	return strings.Join([]string{
		method,
		canonicalURI,
		"", // no query string
		canonicalHeaders.String(),
		strings.Join(signedHeaders, ";"),
		payloadHash,
	}, "\n")
}

// awsEscapePathSegment percent-encodes everything outside RFC 3986's unreserved
// set, which is what SigV4 canonicalization expects for a non-S3 service.
func awsEscapePathSegment(segment string) string {
	var b strings.Builder
	for i := 0; i < len(segment); i++ {
		c := segment[i]
		switch {
		case c >= 'A' && c <= 'Z', c >= 'a' && c <= 'z', c >= '0' && c <= '9',
			c == '-', c == '_', c == '.', c == '~':
			b.WriteByte(c)
		default:
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

func hexSHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func hmacSHA256(key []byte, data string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(data))
	return mac.Sum(nil)
}

// cloudBedrockTransport posts the same request to gx Cloud, which holds the AWS
// credentials. It is the path for everyone without an AWS account, which is
// almost everyone.
type cloudBedrockTransport struct {
	client *cloud.Client
	url    string
}

func (t *cloudBedrockTransport) detail() string {
	if strings.TrimSpace(t.url) == "" {
		return "gx Cloud"
	}
	return "gx Cloud (" + t.url + ")"
}

func (t *cloudBedrockTransport) complete(ctx context.Context, model, system, input string, maxOutputTokens int) (bedrockCompletion, error) {
	if t == nil || t.client == nil {
		return bedrockCompletion{}, fmt.Errorf("Bedrock reviewer has no transport")
	}
	// gx Cloud's /gx/bedrock/fight is a passthrough to Bedrock Converse: it
	// takes this Anthropic-shaped request and normalizes it into a
	// ConverseCommand for whatever model is named, so non-Anthropic IDs
	// (zai.*, nvidia.*, us.openai.*) work here exactly as they do on the
	// direct wire. Whether a given model is ALLOWED is the server's call —
	// it answers model_not_allowed with the list — and that error is more
	// useful than a client-side refusal that pointed at the wrong layer.
	payload := bedrockRequestPayload(system, input, maxOutputTokens)
	resp, err := t.client.BedrockFight(ctx, cloud.BedrockFightRequest{
		Model:            model,
		AnthropicVersion: payload.AnthropicVersion,
		System:           payload.System,
		MaxTokens:        payload.MaxTokens,
		Messages:         []cloud.BedrockMessage{{Role: "user", Content: input}},
	})
	if err != nil {
		var notAllowed *cloud.ModelNotAllowedError
		if errors.As(err, &notAllowed) {
			// Already a complete sentence naming the model and the list.
			return bedrockCompletion{}, err
		}
		return bedrockCompletion{}, fmt.Errorf("gx Cloud Bedrock call for %s failed: %w", model, err)
	}
	text := strings.TrimSpace(resp.Text())
	if text == "" {
		return bedrockCompletion{}, fmt.Errorf("gx Cloud returned no content for Bedrock model %s", model)
	}
	return bedrockCompletion{Text: text, StopReason: resp.StopReason}, nil
}
