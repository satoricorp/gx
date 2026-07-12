package cloud

import (
	"encoding/json"
	"errors"
	"strings"
)

// PaymentRequiredError is returned when GX Cloud rejects an AI request with
// HTTP 402 (expired trial / no active plan).
type PaymentRequiredError struct {
	Message string
}

func (e *PaymentRequiredError) Error() string {
	if e == nil || strings.TrimSpace(e.Message) == "" {
		return "GX Cloud AI needs an active plan; the free trial for this org has ended."
	}
	return e.Message
}

// PaymentRequiredMessage formats the JSON body of a GX Cloud 402 response
// into a user-facing message that names both remedies: upgrading the org, or
// bringing your own model key with gx set key.
func PaymentRequiredMessage(body []byte) string {
	var payload struct {
		Message    string `json:"message"`
		UpgradeURL string `json:"upgrade_url"`
	}
	_ = json.Unmarshal(body, &payload)
	message := strings.TrimSpace(payload.Message)
	if message == "" {
		message = "GX Cloud AI needs an active plan; the free trial for this org has ended. Upgrade to keep using GX Cloud AI, or set your own model key with `gx set key`."
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
