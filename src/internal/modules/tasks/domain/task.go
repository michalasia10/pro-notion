package domain

import (
	"time"

	"github.com/google/uuid"

	shared "src/internal/modules/shared/domain"
)

// SyncStatus represents the synchronisation state between Notion and our system.
type SyncStatus string

const (
	SyncStatusUnknown   SyncStatus = "unknown"
	SyncStatusPending   SyncStatus = "pending"
	SyncStatusCompleted SyncStatus = "completed"
	SyncStatusFailed    SyncStatus = "failed"
)

// Task represents the canonical task aggregate tracked in our system.
type Task struct {
	ID               uuid.UUID
	PublicID         string
	ProjectID        uuid.UUID
	NotionPageID     string
	NotionDatabaseID string
	PropertiesHash   string
	SyncStatus       SyncStatus
	LastSyncedAt     *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

// NewTask creates a new task aggregate with the provided identifiers.
func NewTask(projectID uuid.UUID, notionDatabaseID, notionPageID string, now time.Time, idGenerator shared.IDGenerator) (*Task, error) {
	if projectID == uuid.Nil {
		return nil, ErrInvalidProjectID
	}
	if notionDatabaseID == "" {
		return nil, ErrInvalidNotionDatabaseID
	}
	if notionPageID == "" {
		return nil, ErrInvalidNotionPageID
	}
	if idGenerator == nil {
		return nil, ErrInvalidIDGenerator
	}

	return &Task{
		ID:               uuid.New(),
		PublicID:         idGenerator.NewID("task"),
		ProjectID:        projectID,
		NotionDatabaseID: notionDatabaseID,
		NotionPageID:     notionPageID,
		SyncStatus:       SyncStatusPending,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// MarkSynced sets the synchronisation metadata after a successful sync.
func (t *Task) MarkSynced(at time.Time, propertiesHash string) {
	t.LastSyncedAt = &at
	t.PropertiesHash = propertiesHash
	t.SyncStatus = SyncStatusCompleted
	t.UpdatedAt = at
}

// MarkSyncFailed updates the synchronisation status to failed.
func (t *Task) MarkSyncFailed(at time.Time) {
	t.SyncStatus = SyncStatusFailed
	t.UpdatedAt = at
}

// UpdateNotionIdentifiers allows adjusting Notion linkage if Notion sends new IDs.
func (t *Task) UpdateNotionIdentifiers(databaseID, pageID string, at time.Time) {
	if databaseID != "" {
		t.NotionDatabaseID = databaseID
	}
	if pageID != "" {
		t.NotionPageID = pageID
	}
	t.UpdatedAt = at
}
