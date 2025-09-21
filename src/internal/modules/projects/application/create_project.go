package application

import (
	"context"

	"src/internal/modules/projects/domain"
	shared "src/internal/modules/shared/domain"

	"github.com/google/uuid"
)

// CreateProjectRequest contains the data needed to create a new project
type CreateProjectRequest struct {
	UserID              uuid.UUID
	NotionDatabaseID    string
	NotionWebhookSecret string
}

// CreateProjectResponse contains the created project data
type CreateProjectResponse struct {
	Project domain.Project
}

// CreateProjectUseCase handles project creation business logic
type CreateProjectUseCase struct {
	repo  domain.Repository
	idGen shared.IDGenerator
	clock shared.Clock
	txMgr shared.TransactionManager
}

// NewCreateProjectUseCase creates a new CreateProjectUseCase
func NewCreateProjectUseCase(
	repo domain.Repository,
	idGen shared.IDGenerator,
	clock shared.Clock,
	txMgr shared.TransactionManager,
) *CreateProjectUseCase {
	return &CreateProjectUseCase{
		repo:  repo,
		idGen: idGen,
		clock: clock,
		txMgr: txMgr,
	}
}

// Execute creates a new project
func (uc *CreateProjectUseCase) Execute(ctx context.Context, req CreateProjectRequest) (CreateProjectResponse, error) {
	var response CreateProjectResponse

	err := uc.txMgr.WithinTransaction(ctx, func(ctx context.Context) error {
		// Check if project already exists for this Notion database
		_, err := uc.repo.FindByNotionDatabaseID(ctx, req.NotionDatabaseID)
		if err == nil {
			return domain.ErrProjectNotFound // Project already exists for this database
		}
		if err != domain.ErrProjectNotFound {
			return err
		}

		// Create domain entity (ID and PublicID will be generated inside NewProject)
		project, err := domain.NewProject(req.UserID, req.NotionDatabaseID, req.NotionWebhookSecret, uc.idGen, uc.clock)
		if err != nil {
			return err
		}

		// Save to repository
		err = uc.repo.Save(ctx, &project)
		if err != nil {
			return err
		}

		response = CreateProjectResponse{Project: project}
		return nil
	})

	if err != nil {
		return CreateProjectResponse{}, err
	}

	return response, nil
}
