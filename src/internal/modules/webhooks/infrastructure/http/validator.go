package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"src/internal/modules/webhooks/domain"
)

// HMACSHA256Validator implements WebhookValidator using HMAC-SHA256
type HMACSHA256Validator struct{}

// NewHMACSHA256Validator creates a new HMAC-SHA256 validator
func NewHMACSHA256Validator() *HMACSHA256Validator {
	return &HMACSHA256Validator{}
}

// ValidateSignature validates webhook signature using HMAC-SHA256
func (v *HMACSHA256Validator) ValidateSignature(signature *domain.WebhookSignature, payload domain.WebhookPayload) error {
	if signature.HeaderValue == "" {
		return domain.ErrMissingSignature
	}

	// Expected format: "sha256=<signature>"
	if len(signature.HeaderValue) < 7 || signature.HeaderValue[:7] != "sha256=" {
		return domain.WebhookProcessingError{
			Code:    "INVALID_SIGNATURE_FORMAT",
			Message: "Invalid signature format",
		}
	}

	expectedSignatureHex := signature.HeaderValue[7:] // Remove "sha256=" prefix

	// Compute our signature
	h := hmac.New(sha256.New, []byte(signature.Secret))
	h.Write(payload)
	computedSignature := h.Sum(nil)
	computedSignatureHex := hex.EncodeToString(computedSignature)

	// Compare signatures using constant-time comparison for security
	if !hmac.Equal([]byte(computedSignatureHex), []byte(expectedSignatureHex)) {
		return domain.ErrInvalidSignature
	}

	return nil
}
