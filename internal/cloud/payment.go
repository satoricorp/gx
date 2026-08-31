package cloud

import (
	"encoding/json"
	"errors"
	"strings"
)

// PaymentRequiredError is returned when gx Cloud rejects an AI request with
// HTTP 402: the account's free runs are spent and there is no subscription.
type PaymentRequiredError struct {
	Message string
}

func (e *PaymentRequiredError) Error() string {
	if e == nil || strings.TrimSpace(e.Message) == "" {
		return "gx Cloud AI needs a subscription; your free runs are used up."
	}
	return e.Message
}

// PaymentRequiredMessage formats the JSON body of a gx Cloud 402 response
// into a user-facing message that names the remedy: subscribing at the
// checkout URL the server sends.
func PaymentRequiredMessage(body []byte) string {
	var payload struct {
		Message     string `json:"message"`
		CheckoutURL string `json:"checkout_url"`
		UpgradeURL  string `json:"upgrade_url"`
	}
	_ = json.Unmarshal(body, &payload)
	message := strings.TrimSpace(payload.Message)
	if message == "" {
		message = "gx Cloud AI needs a subscription; your free runs are used up."
	}
	url := strings.TrimSpace(payload.CheckoutURL)
	if url == "" {
		url = strings.TrimSpace(payload.UpgradeURL)
	}
	if url != "" && !strings.Contains(message, url) {
		message += " Subscribe: " + url
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
