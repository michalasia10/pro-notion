package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"src/internal/config"
	"src/internal/modules/webhooks/domain"
	"src/internal/pkg/httpx"
)

type contextKey string

const (
	webhookBodyKey contextKey = "webhook_body"
	signatureVerifiedKey contextKey = "webhook_signature_verified"
)

// WebhookMiddleware validates webhook signature and stores raw body in context
type WebhookMiddleware struct {
	validator domain.WebhookValidator
}

// NewWebhookMiddleware creates a new webhook middleware
func NewWebhookMiddleware(validator domain.WebhookValidator) *WebhookMiddleware {
	return &WebhookMiddleware{
		validator: validator,
	}
}

// Handler returns the HTTP handler for webhook validation
func (m *WebhookMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Get()
		if cfg.Notion.WebhookSecret == "" {
			httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Webhook secret not configured",
			})
			return
		}

		// Read the raw body
		body, err := m.readRequestBody(r)
		if err != nil {
			httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Failed to read request body",
			})
			return
		}

		if isVerificationRequest(body) {
			ctx := context.WithValue(r.Context(), signatureVerifiedKey, true)
			ctx = context.WithValue(ctx, webhookBodyKey, body)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Create signature for validation
		signature := domain.NewWebhookSignature(
			r.Header.Get("X-Notion-Signature"),
			cfg.Notion.WebhookSecret,
		)

		if err := m.validator.ValidateSignature(signature, body); err != nil {
			if os.Getenv("APP_ENV") == "local" {
				log.Printf("webhook signature validation failed: header=%q error=%v", signature.HeaderValue, err)
			}
			statusCode := http.StatusUnauthorized
			if err == domain.ErrInvalidPayload {
				statusCode = http.StatusBadRequest
			}

			httpx.WriteJSON(w, statusCode, map[string]string{
				"error": err.Error(),
			})
			return
		}

		// Signature is valid, store raw body in context
		ctx := context.WithValue(r.Context(), signatureVerifiedKey, true)
		ctx = context.WithValue(ctx, webhookBodyKey, body)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// readRequestBody reads the request body and restores it
func (m *WebhookMiddleware) readRequestBody(r *http.Request) (domain.WebhookPayload, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// Restore the body for subsequent reads
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return body, nil
}

// GetWebhookBody retrieves the raw webhook body from request context
func GetWebhookBody(ctx context.Context) (domain.WebhookPayload, error) {
	body, ok := ctx.Value(webhookBodyKey).(domain.WebhookPayload)
	if !ok {
		return nil, domain.WebhookProcessingError{
			Code:    "CONTEXT_ERROR",
			Message: "webhook body not found in context",
		}
	}
	return body, nil
}

func IsSignatureVerified(ctx context.Context) bool {
	verified, ok := ctx.Value(signatureVerifiedKey).(bool)
	return ok && verified
}

func isVerificationRequest(payload domain.WebhookPayload) bool {
	var req domain.VerificationRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		return false
	}
	return req.IsVerification()
}
