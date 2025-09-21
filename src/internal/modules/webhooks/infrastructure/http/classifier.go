package http

import (
	"encoding/json"

	"src/internal/modules/webhooks/domain"
)

// PayloadClassifier implements WebhookRequestClassifier
type PayloadClassifier struct{}

// NewPayloadClassifier creates a new payload classifier
func NewPayloadClassifier() *PayloadClassifier {
	return &PayloadClassifier{}
}

// ClassifyRequest determines the type of webhook request based on payload content
func (c *PayloadClassifier) ClassifyRequest(payload domain.WebhookPayload) (domain.WebhookEventType, error) {
	var verificationRequest domain.VerificationRequest

	if err := json.Unmarshal(payload, &verificationRequest); err != nil {
		// If we can't unmarshal as verification request, treat as regular event
		return domain.WebhookEventTypeRegular, nil
	}

	if verificationRequest.IsVerification() {
		return domain.WebhookEventTypeVerification, nil
	}

	return domain.WebhookEventTypeRegular, nil
}
