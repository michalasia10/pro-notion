package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"src/internal/modules/projects/domain"
	"src/internal/pkg/txctx"

	"github.com/google/uuid"
)

// ProjectRepository implements domain.Repository using PostgreSQL/GORM
type ProjectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new PostgreSQL project repository
func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) session(ctx context.Context) *gorm.DB {
	if tx := txctx.FromContext(ctx); tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

// Save persists a project
func (r *ProjectRepository) Save(ctx context.Context, project *domain.Project) error {
	record := toProjectRecord(*project)

	if err := r.session(ctx).Create(&record).Error; err != nil {
		if isUniqueViolation(err) {
			return domain.ErrProjectAlreadyExists
		}
		return err
	}

	// Update timestamps only (ID and PublicID are already set in domain)
	project.CreatedAt = record.CreatedAt
	project.UpdatedAt = record.UpdatedAt

	return nil
}

// FindByID retrieves a project by its internal ID
func (r *ProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	var record ProjectRecord

	err := r.session(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrProjectNotFound
		}
		return nil, err
	}

	project := toDomainProject(record)
	return &project, nil
}

// FindByPublicID retrieves a project by its public ID
func (r *ProjectRepository) FindByPublicID(ctx context.Context, publicID string) (*domain.Project, error) {
	var record ProjectRecord

	err := r.session(ctx).Where("public_id = ?", publicID).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrProjectNotFound
		}
		return nil, err
	}

	project := toDomainProject(record)
	return &project, nil
}

// FindByUserID retrieves all projects for a user
func (r *ProjectRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Project, error) {
	var records []ProjectRecord

	err := r.session(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&records).Error

	if err != nil {
		return nil, err
	}

	projects := make([]*domain.Project, 0, len(records))
	for _, record := range records {
		project := toDomainProject(record)
		projects = append(projects, &project)
	}

	return projects, nil
}

// FindByNotionDatabaseID retrieves a project by Notion database ID
func (r *ProjectRepository) FindByNotionDatabaseID(ctx context.Context, notionDatabaseID string) (*domain.Project, error) {
	var record ProjectRecord

	err := r.session(ctx).Where("notion_database_id = ?", notionDatabaseID).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrProjectNotFound
		}
		return nil, err
	}

	project := toDomainProject(record)
	return &project, nil
}

// Update updates an existing project
func (r *ProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	record := toProjectRecord(*project)

	err := r.session(ctx).Save(&record).Error
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrProjectAlreadyExists
		}
		return err
	}

	// Update the domain object with the updated timestamp
	project.UpdatedAt = record.UpdatedAt

	return nil
}

// Delete removes a project
func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.session(ctx).Where("id = ?", id).Delete(&ProjectRecord{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrProjectNotFound
	}

	return nil
}

func isUniqueViolation(err error) bool {
	return errors.Is(err, gorm.ErrDuplicatedKey)
}
