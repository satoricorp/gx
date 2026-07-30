package codereview

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// bedrockAnthropicReviewer is one review leg: a model plus the wire it is
// reached over. It owns prompt construction and response parsing; how the bytes
// get to bedrock-runtime is bedrockTransport's problem (see
// bedrock_transport.go), which is what lets the same leg run against local AWS
// credentials or through totality-cloud without a second implementation.
type bedrockAnthropicReviewer struct {
	model     string
	transport bedrockTransport
}

// bedrockCredentials is what it takes to sign a bedrock-runtime request, and
// the one place that decides whether Bedrock is configured at all.
type bedrockCredentials struct {
	accessKey    string
	secretKey    string
	sessionToken string
	region       string
}

func newBedrockReviewer(transport bedrockTransport, model string) *bedrockAnthropicReviewer {
	return &bedrockAnthropicReviewer{model: model, transport: transport}
}

// bedrockLegLabel names a leg by its model and its transport so a merged
// finding's "Reviewer" evidence says which model actually raised it and over
// which wire. "Bedrock A" alone is a slot, not an attribution; and the same
// model reached two different ways is two different latency and failure stories.
func bedrockLegLabel(slot, model, transportKind string) string {
	short := shortBedrockModelName(model)
	transport := bedrockTransportShortName(transportKind)
	switch {
	case short == "" && transport == "":
		return slot
	case short == "":
		return slot + " (via " + transport + ")"
	case transport == "":
		return slot + " (" + short + ")"
	default:
		return slot + " (" + short + " via " + transport + ")"
	}
}

func shortBedrockModelName(model string) string {
	model = strings.TrimSpace(model)
	if model == "" || strings.HasPrefix(model, "arn:") {
		return model
	}
	if index := strings.LastIndex(model, "anthropic."); index >= 0 {
		model = model[index+len("anthropic."):]
	}
	if index := strings.Index(model, "-v1:"); index > 0 {
		model = model[:index]
	}
	return strings.TrimSuffix(model, "-v1")
}

// bedrockCredentialsFromEnv resolves AWS credentials and region, or says which
// of the two is missing. Returning an error rather than nil is the point: "not
// configured" and "configured wrong" used to be the same message.
func bedrockCredentialsFromEnv() (bedrockCredentials, error) {
	accessKey := strings.TrimSpace(os.Getenv("AWS_ACCESS_KEY_ID"))
	secretKey := strings.TrimSpace(os.Getenv("AWS_SECRET_ACCESS_KEY"))
	if accessKey == "" || secretKey == "" {
		return bedrockCredentials{}, fmt.Errorf("Bedrock reviewer has no AWS credentials: set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY (plus AWS_SESSION_TOKEN for temporary credentials)")
	}
	return bedrockCredentials{
		accessKey:    accessKey,
		secretKey:    secretKey,
		sessionToken: strings.TrimSpace(os.Getenv("AWS_SESSION_TOKEN")),
		region:       bedrockRegionFromEnv(),
	}, nil
}

func bedrockRegionFromEnv() string {
	return strings.TrimSpace(firstNonEmpty(
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_DEFAULT_REGION"),
		defaultBedrockRegion,
	))
}

// resolveBedrockReviewModels picks the model for each leg.
//
// Precedence per leg is env > default. Which model reviews is the operator's
// call: a REVIEW.md line used to win over both, which let the repository under
// review choose the reviewer that judged it — including choosing a weaker one.
// Leg A additionally honours the legacy TOTALITY_REVIEW_ANTHROPIC_MODEL so existing
// setups keep working.
func resolveBedrockReviewModels() (string, string) {
	modelA := normalizeBedrockModelID(firstNonEmpty(
		os.Getenv("TOTALITY_REVIEW_BEDROCK_MODEL_A"),
		os.Getenv("TOTALITY_REVIEW_ANTHROPIC_MODEL"),
		defaultBedrockReviewModelA,
	))
	modelB := normalizeBedrockModelID(firstNonEmpty(
		os.Getenv("TOTALITY_REVIEW_BEDROCK_MODEL_B"),
		defaultBedrockReviewModelB,
	))
	return modelA, modelB
}

// normalizeBedrockModelID upgrades a bare Anthropic model ID to its `us.`
// inference profile. bedrock-runtime rejects bare IDs for on-demand invocation,
// and an operator setting TOTALITY_REVIEW_BEDROCK_MODEL_A to a bare ID is easy to do,
// so an unnormalized ID is a
// guaranteed ValidationException rather than a preference.
//
// ARNs, already-prefixed profiles (us./eu./apac./global./us-gov.), and anything
// that is not an Anthropic model ID are left exactly as written — those are
// deliberate choices a caller made, including provisioned throughput.
func normalizeBedrockModelID(model string) string {
	model = strings.TrimSpace(model)
	if model == "" || strings.HasPrefix(model, "arn:") {
		return model
	}
	for _, prefix := range []string{"us.", "eu.", "apac.", "global.", "us-gov."} {
		if strings.HasPrefix(model, prefix) {
			return model
		}
	}
	if strings.HasPrefix(model, "anthropic.") {
		return "us." + model
	}
	if strings.HasPrefix(model, "claude") {
		return "us.anthropic." + model
	}
	return model
}

func (r *bedrockAnthropicReviewer) Available() bool { return r != nil && r.transport != nil }

func (r *bedrockAnthropicReviewer) Review(ctx context.Context, brief ReviewBrief) ([]Finding, error) {
	summary, err := r.ReviewForSummary(ctx, brief)
	return summary.Findings, err
}

func (r *bedrockAnthropicReviewer) ReviewForSummary(ctx context.Context, brief ReviewBrief) (PRSummaryReview, error) {
	brief = compactReviewBriefForAI(brief)
	completion, err := r.completeJSON(ctx, reviewDeveloperPrompt(brief), mustJSON(brief), defaultReviewMaxOutputTokens)
	if err != nil {
		return PRSummaryReview{}, err
	}
	output, err := parseAIReviewOutput(completion.Text, brief)
	if err != nil {
		if completion.truncated() {
			return PRSummaryReview{}, describeTruncatedCompletion("AI review", defaultReviewMaxOutputTokens, err)
		}
		return PRSummaryReview{}, err
	}
	return aiReviewOutputToPRSummaryReview(output), nil
}

// completeJSON is the single model call. Both reviewer legs and the judge go
// through it, so the request shape and the response parsing live in one place;
// the transport underneath decides whether that request is signed locally or
// posted to totality-cloud.
func (r *bedrockAnthropicReviewer) completeJSON(ctx context.Context, system string, input string, maxOutputTokens int) (bedrockCompletion, error) {
	if r == nil || r.transport == nil {
		return bedrockCompletion{}, fmt.Errorf("Bedrock reviewer has no transport")
	}
	return r.transport.complete(ctx, r.model, system, input, maxOutputTokens)
}

// describeBedrockFailure turns a bedrock-runtime error into a message that says
// what to change.
//
// These four failures are indistinguishable in the raw response to anyone who
// has not read the Bedrock error catalogue — all of them arrive as an exception
// name plus prose — and previously all four surfaced as some flavor of "not
// configured". The mapping is from responses observed against bedrock-runtime
// in us-west-2:
//
//	denied model    AccessDeniedException      "<id> is not available for this account"
//	bare model ID   ValidationException        "with on-demand throughput isn't supported"
//	wrong region    ValidationException        "The provided model identifier is invalid."
//	bad credentials UnrecognizedClientException "The security token ... is invalid"
//
// Note that a bare model ID is a ValidationException, not an AccessDenied: the
// account is allowed to use the model, it just cannot be reached that way.
func describeBedrockFailure(statusCode int, status string, body []byte, model, region string) error {
	detail := strings.TrimSpace(string(body))
	lower := strings.ToLower(detail)
	switch {
	case strings.Contains(lower, "on-demand throughput") || strings.Contains(lower, "inference profile"):
		return fmt.Errorf("Bedrock model %q cannot be invoked on demand; use the inference profile ID %q instead (set TOTALITY_REVIEW_BEDROCK_MODEL_A/_B)", model, normalizeBedrockModelID(model))
	case strings.Contains(lower, "not available for this account"):
		return fmt.Errorf("Bedrock model %q is not enabled for this AWS account in %s; request access in the Bedrock console or point TOTALITY_REVIEW_BEDROCK_MODEL_A/_B at a model you can call", model, region)
	case strings.Contains(lower, "model identifier is invalid") || strings.Contains(lower, "could not resolve the model"):
		return fmt.Errorf("Bedrock model %q does not exist in region %s; set AWS_REGION to a region where it is offered (tx defaults to %s)", model, region, defaultBedrockRegion)
	case strings.Contains(lower, "security token") || strings.Contains(lower, "signature") || strings.Contains(lower, "unrecognizedclient") || statusCode == http.StatusForbidden || statusCode == http.StatusUnauthorized:
		return fmt.Errorf("AWS rejected the Bedrock credentials for %s in %s: check AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY, and AWS_SESSION_TOKEN if they are temporary (%s)", model, region, detail)
	case statusCode == http.StatusTooManyRequests || strings.Contains(lower, "throttl"):
		return fmt.Errorf("Bedrock throttled %s in %s; retry or lower review concurrency (%s)", model, region, detail)
	}
	if detail != "" {
		return fmt.Errorf("Bedrock model %s in %s returned %s: %s", model, region, status, detail)
	}
	return fmt.Errorf("Bedrock model %s in %s returned %s", model, region, status)
}
