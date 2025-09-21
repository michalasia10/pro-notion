package http

import (
	"time"

	"src/internal/modules/projects/domain"
)

// CreateProjectRequestDTO represents the request payload for creating a project
type CreateProjectRequestDTO struct {
	NotionDatabaseID    string `json:"notion_database_id" validate:"required"`
	NotionWebhookSecret string `json:"notion_webhook_secret" validate:"required"`
}

// ProjectResponseDTO represents the response payload for project operations
type ProjectResponseDTO struct {
	ID                  string    `json:"id"`
	UserID              string    `json:"user_id"`
	NotionDatabaseID    string    `json:"notion_database_id"`
	NotionWebhookSecret string    `json:"notion_webhook_secret,omitempty"` // Hide in responses
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ProjectsListResponseDTO represents the response payload for listing projects
type ProjectsListResponseDTO struct {
	Projects []ProjectResponseDTO `json:"projects"`
	Count    int                  `json:"count"`
}

// toProjectResponseDTO converts a domain Project to ProjectResponseDTO
func toProjectResponseDTO(project domain.Project) ProjectResponseDTO {
	return ProjectResponseDTO{
		ID:               project.PublicID, // Use PublicID for API responses
		UserID:           project.UserID.String(),
		NotionDatabaseID: project.NotionDatabaseID,
		// NotionWebhookSecret is omitted for security
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}
}

// toProjectResponseDTOs converts a slice of domain Projects to ProjectResponseDTOs
func toProjectResponseDTOs(projects []*domain.Project) []ProjectResponseDTO {
	dtos := make([]ProjectResponseDTO, 0, len(projects))
	for _, project := range projects {
		dtos = append(dtos, toProjectResponseDTO(*project))
	}
	return dtos
}
