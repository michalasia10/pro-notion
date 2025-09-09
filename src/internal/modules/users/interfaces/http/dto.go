package http

import (
	"time"

	"src/internal/modules/users/domain"
)

// CreateUserRequestDTO represents the request payload for creating a user
type CreateUserRequestDTO struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=1,max=255"`
}

// UserResponseDTO represents the response payload for user operations
type UserResponseDTO struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Notion integration - only include if token exists
	HasNotionIntegration bool       `json:"has_notion_integration"`
	NotionWorkspaceID    string     `json:"notion_workspace_id,omitempty"`
	NotionTokenExpiry    *time.Time `json:"notion_token_expiry,omitempty"`
}

// UsersListResponseDTO represents the response payload for listing users
type UsersListResponseDTO struct {
	Users []UserResponseDTO `json:"users"`
	Count int               `json:"count"`
}

// toUserResponseDTO converts a domain User to UserResponseDTO
func toUserResponseDTO(user domain.User) UserResponseDTO {
	return UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,

		HasNotionIntegration: user.NotionAccessToken != "",
		NotionWorkspaceID:    user.NotionWorkspaceID,
		NotionTokenExpiry:    user.NotionTokenExpiry,
	}
}

// toUserResponseDTOs converts a slice of domain Users to UserResponseDTOs
func toUserResponseDTOs(users []domain.User) []UserResponseDTO {
	dtos := make([]UserResponseDTO, 0, len(users))
	for _, user := range users {
		dtos = append(dtos, toUserResponseDTO(user))
	}
	return dtos
}
