package middleware

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"src/internal/config"
	"src/internal/pkg/httpx"
)

const (
	notionWebhookBodyKey contextKey = "notion_webhook_body"
)

// NotionWebhookMiddleware validates webhook signature and stores raw body in context
func NotionWebhookMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Get()
		if cfg.Notion.WebhookSecret == "" {
			httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Webhook secret not configured",
			})
			return
		}

		// Read the raw body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Failed to read request body",
			})
			return
		}

		// Restore the body for subsequent reads
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// Get signature from header
		signatureHeader := r.Header.Get("X-Notion-Signature")
		if signatureHeader == "" {
			httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Missing X-Notion-Signature header",
			})
			return
		}

		// Expected format: "sha256=<signature>"
		if len(signatureHeader) < 7 || signatureHeader[:7] != "sha256=" {
			httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Invalid signature format",
			})
			return
		}

		expectedSignatureHex := signatureHeader[7:] // Remove "sha256=" prefix

		// Compute our signature
		h := hmac.New(sha256.New, []byte(cfg.Notion.WebhookSecret))
		h.Write(body)
		computedSignature := h.Sum(nil)
		computedSignatureHex := hex.EncodeToString(computedSignature)

		// Compare signatures using constant-time comparison for security
		if !hmac.Equal([]byte(computedSignatureHex), []byte(expectedSignatureHex)) {
			httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Invalid signature",
			})
			return
		}

		// Signature is valid, store raw body in context
		ctx := setNotionWebhookBody(r.Context(), body)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// setNotionWebhookBody stores the raw webhook body in request context
func setNotionWebhookBody(ctx context.Context, body []byte) context.Context {
	return context.WithValue(ctx, notionWebhookBodyKey, body)
}

// GetNotionWebhookBody retrieves the raw webhook body from request context
func GetNotionWebhookBody(ctx context.Context) ([]byte, error) {
	body, ok := ctx.Value(notionWebhookBodyKey).([]byte)
	if !ok {
		return nil, fmt.Errorf("webhook body not found in context")
	}
	return body, nil
}
