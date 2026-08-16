package codereview

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// bedrockAnthropicReviewer is one review leg: a model plus the wire it is
// reached over. It owns prompt construction and response parsing; how the bytes
// get to bedrock-runtime is bedrockTransport's problem (see
// bedrock_transport.go), which is what lets the same leg run against local AWS
// credentials or through gx-cloud without a second implementation.
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

// friendlyModelName turns a Bedrock model ID into the name a person uses:
// us.anthropic.claude-haiku-4-5-20251001-v1:0 → "Claude Haiku 4.5",
// us.anthropic.claude-sonnet-4-6 → "Claude Sonnet 4.6". The report ledger and
// run details show this; the JSON keeps the exact ID, because "which snapshot"
// is a question that comes up the moment a review is slow or wrong.
func friendlyModelName(model string) string {
	short := shortBedrockModelName(model)
	if short == "" || strings.HasPrefix(short, "arn:") {
		return short
	}
	parts := strings.Split(short, "-")
	if len(parts) < 2 || parts[0] != "claude" {
		return short
	}
	family := strings.ToUpper(parts[1][:1]) + parts[1][1:]
	// Version is the run of numeric segments after the family, joined with a
	// dot; a trailing 8-digit date is a snapshot stamp, not part of the version.
	var version []string
	for _, seg := range parts[2:] {
		if len(seg) == 8 && isAllDigits(seg) {
			break
		}
		if isAllDigits(seg) {
			version = append(version, seg)
			continue
		}
		break
	}
	name := "Claude " + family
	if len(version) > 0 {
		name += " " + strings.Join(version, ".")
	}
	return name
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
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
// Leg A additionally honours the legacy GX_REVIEW_ANTHROPIC_MODEL so existing
// setups keep working.
func resolveBedrockReviewModels() (string, string) {
	modelA := normalizeBedrockModelID(firstNonEmpty(
		os.Getenv("GX_REVIEW_BEDROCK_MODEL_A"),
		os.Getenv("GX_REVIEW_ANTHROPIC_MODEL"),
		defaultBedrockReviewModelA,
	))
	modelB := normalizeBedrockModelID(firstNonEmpty(
		os.Getenv("GX_REVIEW_BEDROCK_MODEL_B"),
		defaultBedrockReviewModelB,
	))
	if bedrockLegDisabled(modelA) {
		modelA = ""
	}
	if bedrockLegDisabled(modelB) {
		modelB = ""
	}
	return modelA, modelB
}

// bedrockLegDisabled reports whether a leg was explicitly turned off, which is
// how a single-model review is requested: GX_REVIEW_BEDROCK_MODEL_B=off.
//
// A panel is two models so their misses are uncorrelated, which is worth its
// cost on a pull request. It is not always worth its latency: an interactive
// review a developer is waiting on is a different product than a PR summary
// nobody is watching, and the second leg roughly doubles neither quality nor
// wall clock but does add its slowest call to the critical path.
func bedrockLegDisabled(model string) bool {
	switch strings.ToLower(strings.TrimSpace(model)) {
	case "", "0", "off", "none", "disabled":
		return true
	}
	return false
}

// normalizeBedrockModelID upgrades a bare Anthropic model ID to its `us.`
// inference profile. bedrock-runtime rejects bare IDs for on-demand invocation,
// and an operator setting GX_REVIEW_BEDROCK_MODEL_A to a bare ID is easy to do,
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
	system := reviewDeveloperPrompt(brief)
	input := mustJSON(brief)
	// One retry, for malformed replies only — the same contract the
	// constraints judge runs under. A leg sometimes answers in prose with no
	// JSON at all ("I reviewed the change…"), which extractJSONObject cannot
	// recover; observed live on a binary-only diff, and it turned a healthy
	// panel into a DEGRADED run. The retry restates the JSON-only requirement
	// at the top of the system prompt, which is enough for a model that
	// drifted once. A truncated reply is not retried: the fix for that is the
	// output cap, not another attempt at the same cap.
	var lastParseErr error
	for attempt := 0; attempt < 2; attempt++ {
		prompt := system
		if attempt > 0 {
			prompt = reviewJSONOnlyReminder + "\n" + system
		}
		completion, err := r.completeJSON(ctx, prompt, input, defaultReviewMaxOutputTokens)
		if err != nil {
			return PRSummaryReview{}, err
		}
		output, parseErr := parseAIReviewOutput(completion.Text, brief)
		if parseErr == nil {
			return aiReviewOutputToPRSummaryReview(output), nil
		}
		if errors.Is(parseErr, errTruncatedButSalvaged) {
			// The cap cut the reply, but complete findings were recovered.
			// Return them; the truncation is reported through the summary's
			// Truncated field so the report can say the leg was cut short
			// rather than pretending it answered in full.
			review := aiReviewOutputToPRSummaryReview(output)
			review.Truncated = describeTruncatedCompletion("AI review", defaultReviewMaxOutputTokens, parseErr).Error()
			return review, nil
		}
		if completion.truncated() {
			return PRSummaryReview{}, describeTruncatedCompletion("AI review", defaultReviewMaxOutputTokens, parseErr)
		}
		lastParseErr = parseErr
	}
	return PRSummaryReview{}, lastParseErr
}

// reviewJSONOnlyReminder is prepended to the system prompt on the one retry
// after a reply that carried no decodable JSON. It is deliberately blunt: the
// first prompt already says "Return JSON only", and the model ignored it.
const reviewJSONOnlyReminder = "Your previous reply was not JSON and could not be used. Reply with the JSON object only: no prose before it, no markdown headings, no summary after it. If there is nothing to report, reply exactly {\"recommendations\":[]}."

// completeJSON is the single model call. Both reviewer legs and the judge go
// through it, so the request shape and the response parsing live in one place;
// the transport underneath decides whether that request is signed locally or
// posted to gx-cloud.
func (r *bedrockAnthropicReviewer) completeJSON(ctx context.Context, system string, input string, maxOutputTokens int) (bedrockCompletion, error) {
	if r == nil || r.transport == nil {
		return bedrockCompletion{}, fmt.Errorf("Bedrock reviewer has no transport")
	}
	completion, err := r.transport.complete(ctx, r.model, system, input, maxOutputTokens)
	dumpReviewExchange(r.model, system, input, completion, err)
	return completion, err
}

// dumpReviewExchange writes what a leg actually sent and received to the
// directory named by GX_REVIEW_DUMP_DIR.
//
// A review that returns no findings is indistinguishable, from the outside,
// between "the model read the change and found nothing", "the prompt talked it
// out of reporting", and "the reply did not parse". Those need completely
// different fixes, and the only way to tell them apart is to read the exchange.
// Off unless the variable is set; it writes prompts verbatim, so it is a
// debugging tool and not something to leave on.
func dumpReviewExchange(model, system, input string, completion bedrockCompletion, callErr error) {
	dir := strings.TrimSpace(os.Getenv("GX_REVIEW_DUMP_DIR"))
	if dir == "" {
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}
	safeModel := strings.NewReplacer("/", "_", ":", "_", ".", "-").Replace(model)
	stamp := time.Now().UTC().Format("150405.000")
	body := strings.Join([]string{
		"=== MODEL: " + model,
		"=== SYSTEM PROMPT (" + strconv.Itoa(len(system)) + " bytes):",
		system,
		"=== INPUT BRIEF (" + strconv.Itoa(len(input)) + " bytes):",
		input,
		"=== STOP REASON: " + completion.StopReason,
		"=== ERROR: " + fmt.Sprint(callErr),
		"=== REPLY (" + strconv.Itoa(len(completion.Text)) + " bytes):",
		completion.Text,
	}, "\n")
	_ = os.WriteFile(filepath.Join(dir, safeModel+"-"+stamp+".txt"), []byte(body), 0o644)
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
		return fmt.Errorf("Bedrock model %q cannot be invoked on demand; use the inference profile ID %q instead (set GX_REVIEW_BEDROCK_MODEL_A/_B)", model, normalizeBedrockModelID(model))
	case strings.Contains(lower, "not available for this account"):
		return fmt.Errorf("Bedrock model %q is not enabled for this AWS account in %s; request access in the Bedrock console or point GX_REVIEW_BEDROCK_MODEL_A/_B at a model you can call", model, region)
	case strings.Contains(lower, "model identifier is invalid") || strings.Contains(lower, "could not resolve the model"):
		return fmt.Errorf("Bedrock model %q does not exist in region %s; set AWS_REGION to a region where it is offered (gx defaults to %s)", model, region, defaultBedrockRegion)
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
