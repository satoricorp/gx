package cloud

import (
	"encoding/json"
	"errors"
	"strings"
)

// PaymentRequiredError is returned when Totality Cloud rejects an AI request with
// HTTP 402 (expired trial / no active plan).
type PaymentRequiredError struct {
	Message string
}

func (e *PaymentRequiredError) Error() string {
	if e == nil || strings.TrimSpace(e.Message) == "" {
		return "Totality Cloud AI needs an active plan; the free trial for this org has ended."
	}
	return e.Message
}

// PaymentRequiredMessage formats the JSON body of a Totality Cloud 402 response
// into a user-facing message that names both remedies: upgrading the org, or
// bringing your own model key via ANTHROPIC_API_KEY / OPENAI_API_KEY.
func PaymentRequiredMessage(body []byte) string {
	var payload struct {
		Message    string `json:"message"`
		UpgradeURL string `json:"upgrade_url"`
	}
	_ = json.Unmarshal(body, &payload)
	message := strings.TrimSpace(payload.Message)
	if message == "" {
		message = "Totality Cloud AI needs an active plan; the free trial for this org has ended. Upgrade to keep using Totality Cloud AI, or set your own model key via ANTHROPIC_API_KEY or OPENAI_API_KEY."
	}
	if url := strings.TrimSpace(payload.UpgradeURL); url != "" && !strings.Contains(message, url) {
		message += " Upgrade: " + url
	}
	return message
}

// NewPaymentRequiredError builds a typed 402 error from a response body.
func NewPaymentRequiredError(body []byte) error {
	return &PaymentRequiredError{Message: PaymentRequiredMessage(body)}
}

// IsPaymentRequired reports whether err (or any wrapped cause) is a cloud
// payment/trial gate failure.
func IsPaymentRequired(err error) bool {
	var payment *PaymentRequiredError
	return errors.As(err, &payment)
}
