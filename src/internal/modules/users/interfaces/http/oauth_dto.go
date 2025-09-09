package http

import (
	"src/internal/modules/users/domain"
)

// NotionAuthURLResponseDTO represents the response with authorization URL
type NotionAuthURLResponseDTO struct {
	AuthorizationURL string `json:"authorization_url"`
	State            string `json:"state,omitempty"`
}

// NotionCallbackResponseDTO represents the response after successful OAuth callback
type NotionCallbackResponseDTO struct {
	User     UserResponseDTO `json:"user"`
	Message  string          `json:"message"`
	Success  bool            `json:"success"`
	JWTToken string          `json:"jwt_token,omitempty"` // For future JWT implementation
}

// toNotionCallbackResponseDTO converts domain data to callback response DTO
func toNotionCallbackResponseDTO(user domain.User) NotionCallbackResponseDTO {
	return NotionCallbackResponseDTO{
		User:    toUserResponseDTO(user),
		Message: "Successfully connected to Notion",
		Success: true,
	}
}
