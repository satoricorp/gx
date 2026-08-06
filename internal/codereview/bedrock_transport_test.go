package codereview

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

// The B reviewer leg's inference profile ID. Its trailing ":0" is the only
// character in either default model ID that needs percent-encoding, and it is
// the whole reason the signing bug was invisible on the A leg.
const bedrockColonModel = "us.anthropic.claude-opus-4-5-20251101-v1:0"

// TestBedrockCanonicalRequestMatchesAWSExpectation pins the canonical request
// against the one AWS itself computed.
//
// The expected string below is not invented: bedrock-runtime returns the
// canonical string it derived from the received request inside every
// SignatureDoesNotMatch body, and this is that string, verbatim, from a live
// call to the B leg. Signing the wire path instead of the double-encoded
// canonical path made every B-leg call fail this comparison at "%3A" vs
// "%253A", so the two-model review panel silently ran one model.
func TestBedrockCanonicalRequestMatchesAWSExpectation(t *testing.T) {
	const (
		amzDate     = "20260727T195250Z"
		payloadHash = "6e915c5e44cb371331d2077821eea63ae108ecc4279e723d2434706406ae1ca5"
		host        = "bedrock-runtime.us-west-2.amazonaws.com"
	)
	signed := []string{"content-type", "host", "x-amz-content-sha256", "x-amz-date"}
	values := map[string]string{
		"content-type":         "application/json",
		"host":                 host,
		"x-amz-content-sha256": payloadHash,
		"x-amz-date":           amzDate,
	}

	want := strings.Join([]string{
		"POST",
		"/model/us.anthropic.claude-opus-4-5-20251101-v1%253A0/invoke",
		"",
		"content-type:application/json",
		"host:" + host,
		"x-amz-content-sha256:" + payloadHash,
		"x-amz-date:" + amzDate,
		"",
		"content-type;host;x-amz-content-sha256;x-amz-date",
		payloadHash,
	}, "\n")

	got := bedrockCanonicalRequest(http.MethodPost, bedrockCanonicalPath(bedrockColonModel), signed, values, payloadHash)
	if got != want {
		t.Fatalf("canonical request mismatch\n got: %q\nwant: %q", got, want)
	}
}

// TestBedrockPathEncodingLayers states the asymmetry directly: the wire path is
// encoded once, the canonical path twice, and for a model ID that needs no
// escaping the two coincide — which is why the A leg never noticed.
func TestBedrockPathEncodingLayers(t *testing.T) {
	if got, want := bedrockRequestPath(bedrockColonModel), "/model/us.anthropic.claude-opus-4-5-20251101-v1%3A0/invoke"; got != want {
		t.Errorf("bedrockRequestPath() = %q, want %q", got, want)
	}
	if got, want := bedrockCanonicalPath(bedrockColonModel), "/model/us.anthropic.claude-opus-4-5-20251101-v1%253A0/invoke"; got != want {
		t.Errorf("bedrockCanonicalPath() = %q, want %q", got, want)
	}
	// A literal, not defaultBedrockReviewModelA: this case is about an ID with
	// nothing to escape, and tying it to whichever model is currently the
	// default made it fail the moment the default changed to one whose ID ends
	// in ":0" — which is the escaping case the assertions above already cover.
	plain := "us.anthropic.claude-opus-4-6-v1"
	if bedrockRequestPath(plain) != bedrockCanonicalPath(plain) {
		t.Errorf("model %q needs no escaping, so both paths must agree; got %q and %q",
			plain, bedrockRequestPath(plain), bedrockCanonicalPath(plain))
	}
}

// TestSignedRequestSendsSingleEncodedPath guards the other half of the pair: the
// bytes on the wire must stay single-encoded, or Bedrock 404s on a model ID it
// would otherwise recognise.
func TestSignedRequestSendsSingleEncodedPath(t *testing.T) {
	transport := newDirectBedrockTransport(bedrockCredentials{
		accessKey: "AKIAEXAMPLE",
		secretKey: "secret",
		region:    "us-west-2",
	})
	req, err := transport.newSignedRequest(context.Background(), bedrockColonModel, []byte(`{}`))
	if err != nil {
		t.Fatalf("newSignedRequest() error = %v", err)
	}
	if got, want := req.URL.EscapedPath(), "/model/us.anthropic.claude-opus-4-5-20251101-v1%3A0/invoke"; got != want {
		t.Errorf("request path = %q, want %q", got, want)
	}
	if auth := req.Header.Get("Authorization"); !strings.HasPrefix(auth, "AWS4-HMAC-SHA256 Credential=AKIAEXAMPLE/") {
		t.Errorf("Authorization header = %q, want an AWS4-HMAC-SHA256 credential", auth)
	}
}
