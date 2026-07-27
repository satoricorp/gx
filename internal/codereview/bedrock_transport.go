package codereview

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/cloud"
)

// bedrockTransport is how a review leg's bytes reach bedrock-runtime.
//
// Both legs and the judge share one request shape and one response parser (see
// bedrockAnthropicReviewer.completeJSON); only the wire differs. Keeping that
// difference behind this interface is what lets the same reviewer run against a
// developer's own AWS credentials and against GX Cloud without a second
// implementation of prompt construction, parsing, or failure reporting.
type bedrockTransport interface {
	// complete performs one model call and returns the model's text.
	complete(ctx context.Context, model, system, input string, maxOutputTokens int) (string, error)
	// detail is the human phrase naming this wire, reported on every review so
	// that "slow" and "denied" stay answerable after the fact.
	detail() string
}

const (
	bedrockTransportKindDirect = "direct"
	bedrockTransportKindCloud  = "cloud"

	// bedrockAnthropicVersion is the Anthropic API version bedrock-runtime
	// expects in an InvokeModel body. GX Cloud accepts and ignores it.
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

// resolveBedrockTransportPlan picks the wire: local AWS credentials if they are
// present, GX Cloud otherwise.
//
// That order is deliberate. A developer who has exported AWS credentials is
// asking for their own account and their own quota, and silently routing them
// through GX Cloud would spend GX's budget and hide their misconfiguration.
// Everyone else — the common case, a user with no AWS account at all — reaches
// the same models through GX Cloud.
//
// When neither is available the error names the local fix, because that is the
// one the caller can act on without an account. It is returned rather than
// swallowed: a review with no reviewer must read as "nothing reviewed this",
// never as a clean review.
func resolveBedrockTransportPlan() (bedrockTransportPlan, error) {
	creds, credsErr := bedrockCredentialsFromEnv()
	if credsErr == nil {
		return bedrockTransportPlan{Kind: bedrockTransportKindDirect, creds: creds}, nil
	}
	client := cloud.NewBedrockClient()
	if client == nil {
		return bedrockTransportPlan{}, credsErr
	}
	if _, err := cloud.CloudAPIToken(); err != nil {
		// Cloud is reachable but this machine is not signed in. Say both, so
		// the reader can pick whichever is cheaper for them to fix.
		return bedrockTransportPlan{}, fmt.Errorf("%v; or sign in with `gx auth login` to review through GX Cloud (%v)", credsErr, err)
	}
	return bedrockTransportPlan{Kind: bedrockTransportKindCloud, client: client, cloudURL: cloud.CloudURL()}, nil
}

// bedrockTransportShortName is the label fragment for one wire, short enough to
// sit inside a finding's reviewer attribution.
func bedrockTransportShortName(kind string) string {
	switch kind {
	case bedrockTransportKindDirect:
		return "AWS"
	case bedrockTransportKindCloud:
		return "GX Cloud"
	default:
		return ""
	}
}

// bedrockRequestBody is the Anthropic messages shape bedrock-runtime's
// InvokeModel takes, and the shape GX Cloud's /gx/bedrock/fight normalizes.
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

func (t *directBedrockTransport) complete(ctx context.Context, model, system, input string, maxOutputTokens int) (string, error) {
	if t == nil {
		return "", fmt.Errorf("Bedrock reviewer has no transport")
	}
	body, err := json.Marshal(bedrockRequestPayload(system, input, maxOutputTokens))
	if err != nil {
		return "", fmt.Errorf("marshal Bedrock request: %w", err)
	}
	req, err := t.newSignedRequest(ctx, model, body)
	if err != nil {
		return "", err
	}
	client := t.client
	if client == nil {
		client = &http.Client{Timeout: bedrockDirectTimeout}
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("call Bedrock model %s in %s: %w", model, t.region, err)
	}
	defer resp.Body.Close()
	raw, readErr := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// describeBedrockFailure turns the exception name and prose into the
		// knob to turn; see its comment in ai.go.
		return "", describeBedrockFailure(resp.StatusCode, resp.Status, raw, model, t.region)
	}
	if readErr != nil {
		return "", fmt.Errorf("read Bedrock response for %s: %w", model, readErr)
	}
	var decoded bedrockResponseBody
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return "", fmt.Errorf("decode Bedrock response for %s: %w", model, err)
	}
	text := decoded.text()
	if text == "" {
		return "", fmt.Errorf("Bedrock model %s returned no content", model)
	}
	return text, nil
}

func (t *directBedrockTransport) newSignedRequest(ctx context.Context, model string, body []byte) (*http.Request, error) {
	model = strings.TrimSpace(model)
	if model == "" {
		return nil, fmt.Errorf("Bedrock request has no model")
	}
	host := "bedrock-runtime." + t.region + ".amazonaws.com"
	// The request line carries the path escaped once; the canonical request
	// carries it escaped twice. See bedrockCanonicalPath for why they differ.
	rawPath := bedrockRequestPath(model)
	endpoint := &url.URL{
		Scheme:  "https",
		Host:    host,
		Path:    "/model/" + model + "/invoke",
		RawPath: rawPath,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create Bedrock request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if err := t.signV4(req, host, bedrockCanonicalPath(model), body, time.Now().UTC()); err != nil {
		return nil, err
	}
	return req, nil
}

// bedrockRequestPath is the path that goes on the wire: each path segment
// percent-encoded once, which is what an HTTP request line requires.
func bedrockRequestPath(model string) string {
	return "/model/" + awsEscapePathSegment(model) + "/invoke"
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
func bedrockCanonicalPath(model string) string {
	return "/model/" + awsEscapePathSegment(awsEscapePathSegment(model)) + "/invoke"
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

// cloudBedrockTransport posts the same request to GX Cloud, which holds the AWS
// credentials. It is the path for everyone without an AWS account, which is
// almost everyone.
type cloudBedrockTransport struct {
	client *cloud.Client
	url    string
}

func (t *cloudBedrockTransport) detail() string {
	if strings.TrimSpace(t.url) == "" {
		return "GX Cloud"
	}
	return "GX Cloud (" + t.url + ")"
}

func (t *cloudBedrockTransport) complete(ctx context.Context, model, system, input string, maxOutputTokens int) (string, error) {
	if t == nil || t.client == nil {
		return "", fmt.Errorf("Bedrock reviewer has no transport")
	}
	payload := bedrockRequestPayload(system, input, maxOutputTokens)
	resp, err := t.client.BedrockFight(ctx, cloud.BedrockFightRequest{
		Model:            model,
		AnthropicVersion: payload.AnthropicVersion,
		System:           payload.System,
		MaxTokens:        payload.MaxTokens,
		Messages:         []cloud.BedrockMessage{{Role: "user", Content: input}},
	})
	if err != nil {
		return "", fmt.Errorf("GX Cloud Bedrock call for %s failed: %w", model, err)
	}
	text := strings.TrimSpace(resp.Text())
	if text == "" {
		return "", fmt.Errorf("GX Cloud returned no content for Bedrock model %s", model)
	}
	return text, nil
}
