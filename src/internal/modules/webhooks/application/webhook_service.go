package application

import (
	"context"
	"encoding/json"

	shared "src/internal/modules/shared/domain"
	"src/internal/modules/webhooks/domain"
)

// WebhookProcessingRequest represents a request to process a webhook
type WebhookProcessingRequest struct {
	Payload        domain.WebhookPayload
	Signature      *domain.WebhookSignature
	RequestContext context.Context
}

// WebhookProcessingResponse represents the response from webhook processing
type WebhookProcessingResponse struct {
	Event   *domain.WebhookEvent
	Message string
	Success bool
}

// WebhookService handles webhook processing use cases
type WebhookService struct {
	validator   domain.WebhookValidator
	classifier  domain.WebhookRequestClassifier
	publisher   domain.WebhookEventPublisher
	idGenerator shared.IDGenerator
}

// NewWebhookService creates a new webhook service
func NewWebhookService(
	validator domain.WebhookValidator,
	classifier domain.WebhookRequestClassifier,
	publisher domain.WebhookEventPublisher,
	idGenerator shared.IDGenerator,
) *WebhookService {
	return &WebhookService{
		validator:   validator,
		classifier:  classifier,
		publisher:   publisher,
		idGenerator: idGenerator,
	}
}

// ProcessWebhook processes a webhook request
func (s *WebhookService) ProcessWebhook(ctx context.Context, req WebhookProcessingRequest) (WebhookProcessingResponse, error) {
	// Validate signature
	if err := s.validator.ValidateSignature(req.Signature, req.Payload); err != nil {
		return WebhookProcessingResponse{Success: false}, err
	}

	// Classify the request
	eventType, err := s.classifier.ClassifyRequest(req.Payload)
	if err != nil {
		return WebhookProcessingResponse{Success: false}, domain.WebhookProcessingError{
			Code:    "CLASSIFICATION_FAILED",
			Message: err.Error(),
		}
	}

	// Create event
	event := domain.NewWebhookEvent(
		s.idGenerator.NewID("webhook"),
		req.Payload,
		eventType,
	)

	// For verification events, don't publish to avoid unnecessary processing
	if eventType != domain.WebhookEventTypeVerification {
		if err := s.publisher.PublishEvent(event); err != nil {
			return WebhookProcessingResponse{Success: false}, err
		}
	}

	message := "Webhook verification received"
	if eventType == domain.WebhookEventTypeRegular {
		message = "Webhook event processed successfully"
	}

	return WebhookProcessingResponse{
		Event:   event,
		Message: message,
		Success: true,
	}, nil
}

// VerificationTokenResponse represents a verification token response
type VerificationTokenResponse struct {
	Token   string
	Message string
}

// ExtractVerificationToken extracts verification token from payload
func (s *WebhookService) ExtractVerificationToken(payload domain.WebhookPayload) (*VerificationTokenResponse, error) {
	var req domain.VerificationRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}

	return &VerificationTokenResponse{
		Token:   req.VerificationToken,
		Message: "Verification token received. Please use this token in your Notion integration settings.",
	}, nil
}
