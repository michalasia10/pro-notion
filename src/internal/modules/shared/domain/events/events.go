package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	NotionWebhookReceivedTopic     = "notion.webhook.received"
	TaskPropertiesUpdatedTopic     = "task.properties.updated"
	TaskCreatedTopic               = "task.created"
	TaskDeletedTopic               = "task.deleted"
	DatabasePropertiesUpdatedTopic = "database.properties.updated"
)

type NotionWebhookReceived struct {
	Payload []byte
}

type TaskPropertiesUpdated struct {
	TaskID     uuid.UUID
	DatabaseID string
	PageID     string
	Changes    map[string]interface{}
	Timestamp  time.Time
}

func NewTaskPropertiesUpdated(taskID uuid.UUID, databaseID, pageID string, changes map[string]interface{}) *TaskPropertiesUpdated {
	return &TaskPropertiesUpdated{
		TaskID:     taskID,
		DatabaseID: databaseID,
		PageID:     pageID,
		Changes:    changes,
		Timestamp:  time.Now(),
	}
}

type TaskCreated struct {
	TaskID     uuid.UUID
	DatabaseID string
	PageID     string
	Properties map[string]interface{}
	Timestamp  time.Time
}

func NewTaskCreated(taskID uuid.UUID, databaseID, pageID string, properties map[string]interface{}) *TaskCreated {
	return &TaskCreated{
		TaskID:     taskID,
		DatabaseID: databaseID,
		PageID:     pageID,
		Properties: properties,
		Timestamp:  time.Now(),
	}
}

type TaskDeleted struct {
	TaskID     uuid.UUID
	DatabaseID string
	PageID     string
	Timestamp  time.Time
}

func NewTaskDeleted(taskID uuid.UUID, databaseID, pageID string) *TaskDeleted {
	return &TaskDeleted{
		TaskID:     taskID,
		DatabaseID: databaseID,
		PageID:     pageID,
		Timestamp:  time.Now(),
	}
}

type DatabasePropertiesUpdated struct {
	DatabaseID string
	PageID     string
	Changes    map[string]interface{}
	Timestamp  time.Time
}

func NewDatabasePropertiesUpdated(databaseID, pageID string, changes map[string]interface{}) *DatabasePropertiesUpdated {
	return &DatabasePropertiesUpdated{
		DatabaseID: databaseID,
		PageID:     pageID,
		Changes:    changes,
		Timestamp:  time.Now(),
	}
}
