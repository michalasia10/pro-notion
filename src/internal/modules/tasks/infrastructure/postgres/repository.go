package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"src/internal/database"
	taskdomain "src/internal/modules/tasks/domain"
)

// Repository implements taskdomain.Repository using GORM.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a tasks repository backed by the shared GORM instance.
func NewRepository(db *gorm.DB) *Repository {
	if db == nil {
		db = database.GormDB()
	}
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, task *taskdomain.Task) error {
	record := FromDomain(task)
	return r.db.WithContext(ctx).Create(&record).Error
}

func (r *Repository) Update(ctx context.Context, task *taskdomain.Task) error {
	record := FromDomain(task)
	return r.db.WithContext(ctx).
		Model(&TaskRecord{}).
		Where("id = ?", task.ID).
		Updates(map[string]interface{}{
			"public_id":          record.PublicID,
			"project_id":         record.ProjectID,
			"notion_page_id":     record.NotionPageID,
			"notion_database_id": record.NotionDatabaseID,
			"properties_hash":    record.PropertiesHash,
			"sync_status":        record.SyncStatus,
			"last_synced_at":     record.LastSyncedAt,
			"updated_at":         record.UpdatedAt,
			"deleted_at":         record.DeletedAt,
		}).Error
}

func (r *Repository) Upsert(ctx context.Context, task *taskdomain.Task) error {
	record := FromDomain(task)
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "notion_database_id"}, {Name: "notion_page_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"project_id":      record.ProjectID,
				"properties_hash": record.PropertiesHash,
				"sync_status":     record.SyncStatus,
				"last_synced_at":  record.LastSyncedAt,
				"updated_at":      record.UpdatedAt,
				"deleted_at":      record.DeletedAt,
			}),
		}).
		Create(&record).Error
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*taskdomain.Task, error) {
	var record TaskRecord
	if err := r.db.WithContext(ctx).First(&record, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, taskdomain.ErrTaskNotFound
		}
		return nil, err
	}
	return record.ToDomain()
}

func (r *Repository) FindByPublicID(ctx context.Context, publicID string) (*taskdomain.Task, error) {
	var record TaskRecord
	if err := r.db.WithContext(ctx).First(&record, "public_id = ?", publicID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, taskdomain.ErrTaskNotFound
		}
		return nil, err
	}
	return record.ToDomain()
}

func (r *Repository) FindByNotionIDs(ctx context.Context, notionDatabaseID, notionPageID string) (*taskdomain.Task, error) {
	var record TaskRecord
	if err := r.db.WithContext(ctx).
		First(&record, "notion_database_id = ? AND notion_page_id = ?", notionDatabaseID, notionPageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, taskdomain.ErrTaskNotFound
		}
		return nil, err
	}
	return record.ToDomain()
}

func (r *Repository) WithTransaction(ctx context.Context, fn func(ctx context.Context, repo taskdomain.Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &Repository{db: tx}
		return fn(ctx, txRepo)
	})
}
