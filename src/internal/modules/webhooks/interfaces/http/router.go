package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi/v5"

	shared "src/internal/modules/shared/domain"
	sharedEvents "src/internal/modules/shared/domain/events"
	"src/internal/pkg/httpx"
	"src/internal/pkg/middleware"
)

// NewRouter creates a new HTTP router for the webhooks module
func NewRouter(publisher message.Publisher) chi.Router {
	r := chi.NewRouter()

	r.Post("/notion", httpx.Endpoint(func(req *http.Request) (int, any, error) {
		// Get the validated webhook body from middleware
		payload, err := middleware.GetNotionWebhookBody(req.Context())
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		// Check if this is a verification request
		var verificationRequest struct {
			VerificationToken string `json:"verification_token"`
		}

		if err := json.Unmarshal(payload, &verificationRequest); err == nil && verificationRequest.VerificationToken != "" {
			// This is a verification request from Notion
			log.Printf("Received Notion webhook verification token: %s", verificationRequest.VerificationToken)

			// Return success response - the token should be used in the Notion dashboard
			return http.StatusOK, map[string]string{
				"message": "Verification token received. Please use this token in your Notion integration settings.",
				"token":   verificationRequest.VerificationToken,
			}, nil
		}

		// This is a regular webhook event - publish to Watermill
		event := sharedEvents.NotionWebhookReceived{
			Payload: payload,
		}

		eventBytes, err := json.Marshal(event)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		msg := message.NewMessage(shared.NewUUIDGenerator().NewID("webhook"), eventBytes)

		if err := publisher.Publish(sharedEvents.NotionWebhookReceivedTopic, msg); err != nil {
			log.Printf("Failed to publish webhook event: %v", err)
			return http.StatusInternalServerError, nil, err
		}

		log.Printf("Successfully published webhook event to topic: %s", sharedEvents.NotionWebhookReceivedTopic)

		return http.StatusOK, map[string]string{
			"message": "Webhook event processed successfully",
		}, nil
	}))

	return r
}
