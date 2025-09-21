package domain

import "fmt"

// WebhookPayload represents the raw payload received from a webhook
type WebhookPayload []byte

// String returns string representation of the payload
func (p WebhookPayload) String() string {
	return string(p)
}

// VerificationRequest represents a webhook verification request from Notion
type VerificationRequest struct {
	VerificationToken string `json:"verification_token"`
}

// IsVerification checks if this is a verification request
func (r VerificationRequest) IsVerification() bool {
	return r.VerificationToken != ""
}

// WebhookEvent represents a processed webhook event
type WebhookEvent struct {
	ID      string
	Payload WebhookPayload
	Type    WebhookEventType
}

// WebhookEventType represents the type of webhook event
type WebhookEventType string

const (
	WebhookEventTypeVerification WebhookEventType = "verification"
	WebhookEventTypeRegular      WebhookEventType = "regular"
)

// NewWebhookEvent creates a new webhook event
func NewWebhookEvent(id string, payload WebhookPayload, eventType WebhookEventType) *WebhookEvent {
	return &WebhookEvent{
		ID:      id,
		Payload: payload,
		Type:    eventType,
	}
}

// WebhookSignature represents a webhook signature for validation
type WebhookSignature struct {
	HeaderValue string
	Secret      string
}

// NewWebhookSignature creates a new webhook signature
func NewWebhookSignature(headerValue, secret string) *WebhookSignature {
	return &WebhookSignature{
		HeaderValue: headerValue,
		Secret:      secret,
	}
}

// WebhookValidator handles webhook signature validation
type WebhookValidator interface {
	ValidateSignature(signature *WebhookSignature, payload WebhookPayload) error
}

// WebhookEventPublisher handles publishing webhook events
type WebhookEventPublisher interface {
	PublishEvent(event *WebhookEvent) error
}

// WebhookRequestClassifier classifies webhook requests
type WebhookRequestClassifier interface {
	ClassifyRequest(payload WebhookPayload) (WebhookEventType, error)
}

// WebhookProcessingError represents an error during webhook processing
type WebhookProcessingError struct {
	Message string
	Code    string
}

func (e WebhookProcessingError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Common webhook processing errors
var (
	ErrInvalidSignature = WebhookProcessingError{Code: "INVALID_SIGNATURE", Message: "Invalid webhook signature"}
	ErrMissingSignature = WebhookProcessingError{Code: "MISSING_SIGNATURE", Message: "Missing webhook signature header"}
	ErrInvalidPayload   = WebhookProcessingError{Code: "INVALID_PAYLOAD", Message: "Invalid webhook payload"}
	ErrProcessingFailed = WebhookProcessingError{Code: "PROCESSING_FAILED", Message: "Failed to process webhook"}
)
