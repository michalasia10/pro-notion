package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"src/internal/config"
	"src/internal/modules/webhooks/application"
	"src/internal/modules/webhooks/domain"
	webhookInfra "src/internal/modules/webhooks/infrastructure/http"
	"src/internal/pkg/httpx"

	shared "src/internal/modules/shared/domain"

	"github.com/ThreeDotsLabs/watermill/message"
)

// WebhookHandler handles HTTP requests for webhooks
type WebhookHandler struct {
	webhookService application.WebhookService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(service application.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: service,
	}
}

// HandleWebhook processes webhook requests
func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Get validated payload from middleware
	payload, err := webhookInfra.GetWebhookBody(r.Context())
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Get signature for processing
	cfg := config.Get()
	signature := domain.NewWebhookSignature(
		r.Header.Get("X-Notion-Signature"),
		cfg.Notion.WebhookSecret,
	)

	// Process webhook through application service
	req := application.WebhookProcessingRequest{
		Payload:        payload,
		Signature:      signature,
		RequestContext: r.Context(),
	}

	response, err := h.webhookService.ProcessWebhook(r.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if domainErr, ok := err.(domain.WebhookProcessingError); ok {
			switch domainErr.Code {
			case "INVALID_SIGNATURE", "MISSING_SIGNATURE":
				statusCode = http.StatusUnauthorized
			case "INVALID_PAYLOAD":
				statusCode = http.StatusBadRequest
			}
		}

		httpx.WriteJSON(w, statusCode, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Handle verification response differently
	if response.Event.Type == domain.WebhookEventTypeVerification {
		verificationResp, err := h.webhookService.ExtractVerificationToken(payload)
		if err != nil {
			httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid verification request",
			})
			return
		}

		httpx.WriteJSON(w, http.StatusOK, map[string]string{
			"message": verificationResp.Message,
			"token":   verificationResp.Token,
		})
		return
	}

	// Regular event response
	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"message": response.Message,
	})
}

// NewRouter creates a new HTTP router for the webhooks module
func NewRouter(publisher message.Publisher) chi.Router {
	r := chi.NewRouter()

	// Initialize dependencies
	validator := webhookInfra.NewHMACSHA256Validator()
	classifier := webhookInfra.NewPayloadClassifier()
	eventPublisher := webhookInfra.NewWatermillEventPublisher(publisher, log.Default())
	idGenerator := shared.NewUUIDGenerator()

	// Initialize application service
	webhookService := application.NewWebhookService(
		validator,
		classifier,
		eventPublisher,
		idGenerator,
	)

	// Initialize middleware
	middleware := webhookInfra.NewWebhookMiddleware(validator)

	// Initialize handler
	handler := NewWebhookHandler(*webhookService)

	// Setup routes with middleware
	r.With(middleware.Handler).Post("/notion", handler.HandleWebhook)

	return r
}
