package dto

import (
	"src/internal/modules/shared/domain/events"

	"github.com/google/uuid"
)

// DTOs for JSON serialization/deserialization

// TaskPropertiesUpdatedDTO represents the JSON serializable version of TaskPropertiesUpdated
type TaskPropertiesUpdatedDTO struct {
	TaskID     uuid.UUID              `json:"task_id"`
	DatabaseID string                 `json:"database_id"`
	PageID     string                 `json:"page_id"`
	Changes    map[string]interface{} `json:"changes"`
	Timestamp  int64                  `json:"timestamp"`
}

// ToDomain converts DTO to domain event
func (dto *TaskPropertiesUpdatedDTO) ToDomain() *events.TaskPropertiesUpdated {
	return &events.TaskPropertiesUpdated{
		TaskID:     dto.TaskID,
		DatabaseID: dto.DatabaseID,
		PageID:     dto.PageID,
		Changes:    dto.Changes,
		// Note: timestamp is set to current time in domain constructor
	}
}

// FromDomain converts domain event to DTO
func TaskPropertiesUpdatedDTOFromDomain(event *events.TaskPropertiesUpdated) *TaskPropertiesUpdatedDTO {
	return &TaskPropertiesUpdatedDTO{
		TaskID:     event.TaskID,
		DatabaseID: event.DatabaseID,
		PageID:     event.PageID,
		Changes:    event.Changes,
		Timestamp:  event.Timestamp.Unix(),
	}
}

// TaskCreatedDTO represents the JSON serializable version of TaskCreated
type TaskCreatedDTO struct {
	TaskID     uuid.UUID              `json:"task_id"`
	DatabaseID string                 `json:"database_id"`
	PageID     string                 `json:"page_id"`
	Properties map[string]interface{} `json:"properties"`
	Timestamp  int64                  `json:"timestamp"`
}

// ToDomain converts DTO to domain event
func (dto *TaskCreatedDTO) ToDomain() *events.TaskCreated {
	return &events.TaskCreated{
		TaskID:     dto.TaskID,
		DatabaseID: dto.DatabaseID,
		PageID:     dto.PageID,
		Properties: dto.Properties,
		// Note: timestamp is set to current time in domain constructor
	}
}

// FromDomain converts domain event to DTO
func TaskCreatedDTOFromDomain(event *events.TaskCreated) *TaskCreatedDTO {
	return &TaskCreatedDTO{
		TaskID:     event.TaskID,
		DatabaseID: event.DatabaseID,
		PageID:     event.PageID,
		Properties: event.Properties,
		Timestamp:  event.Timestamp.Unix(),
	}
}

// TaskDeletedDTO represents the JSON serializable version of TaskDeleted
type TaskDeletedDTO struct {
	TaskID     uuid.UUID `json:"task_id"`
	DatabaseID string    `json:"database_id"`
	PageID     string    `json:"page_id"`
	Timestamp  int64     `json:"timestamp"`
}

// ToDomain converts DTO to domain event
func (dto *TaskDeletedDTO) ToDomain() *events.TaskDeleted {
	return &events.TaskDeleted{
		TaskID:     dto.TaskID,
		DatabaseID: dto.DatabaseID,
		PageID:     dto.PageID,
		// Note: timestamp is set to current time in domain constructor
	}
}

// FromDomain converts domain event to DTO
func TaskDeletedDTOFromDomain(event *events.TaskDeleted) *TaskDeletedDTO {
	return &TaskDeletedDTO{
		TaskID:     event.TaskID,
		DatabaseID: event.DatabaseID,
		PageID:     event.PageID,
		Timestamp:  event.Timestamp.Unix(),
	}
}

// DatabasePropertiesUpdatedDTO represents the JSON serializable version of DatabasePropertiesUpdated
type DatabasePropertiesUpdatedDTO struct {
	DatabaseID string                 `json:"database_id"`
	PageID     string                 `json:"page_id"`
	Changes    map[string]interface{} `json:"changes"`
	Timestamp  int64                  `json:"timestamp"`
}

// ToDomain converts DTO to domain event
func (dto *DatabasePropertiesUpdatedDTO) ToDomain() *events.DatabasePropertiesUpdated {
	return &events.DatabasePropertiesUpdated{
		DatabaseID: dto.DatabaseID,
		PageID:     dto.PageID,
		Changes:    dto.Changes,
		// Note: timestamp is set to current time in domain constructor
	}
}

// FromDomain converts domain event to DTO
func DatabasePropertiesUpdatedDTOFromDomain(event *events.DatabasePropertiesUpdated) *DatabasePropertiesUpdatedDTO {
	return &DatabasePropertiesUpdatedDTO{
		DatabaseID: event.DatabaseID,
		PageID:     event.PageID,
		Changes:    event.Changes,
		Timestamp:  event.Timestamp.Unix(),
	}
}

