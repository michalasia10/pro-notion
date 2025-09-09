package application

import (
	"context"
	"time"

	shared "src/internal/modules/shared/domain"
	"src/internal/modules/users/domain"
)

// UpdateNotionTokenRequest contains the data needed to update a user's Notion token
type UpdateNotionTokenRequest struct {
	UserID      string
	AccessToken string
	WorkspaceID string
	BotID       string
	TokenExpiry *time.Time
}

// UpdateNotionTokenResponse contains the updated user data
type UpdateNotionTokenResponse struct {
	User domain.User
}

// UpdateNotionTokenUseCase handles updating user's Notion token
type UpdateNotionTokenUseCase struct {
	repo  domain.UserRepository
	clock shared.Clock
	txMgr shared.TransactionManager
}

// NewUpdateNotionTokenUseCase creates a new UpdateNotionTokenUseCase
func NewUpdateNotionTokenUseCase(
	repo domain.UserRepository,
	clock shared.Clock,
	txMgr shared.TransactionManager,
) *UpdateNotionTokenUseCase {
	return &UpdateNotionTokenUseCase{
		repo:  repo,
		clock: clock,
		txMgr: txMgr,
	}
}

// Execute updates a user's Notion token information
func (uc *UpdateNotionTokenUseCase) Execute(ctx context.Context, req UpdateNotionTokenRequest) (UpdateNotionTokenResponse, error) {
	var response UpdateNotionTokenResponse

	err := uc.txMgr.WithinTransaction(ctx, func(ctx context.Context) error {
		// Get existing user
		user, err := uc.repo.GetByID(ctx, req.UserID)
		if err != nil {
			return err
		}

		// Update Notion token information
		user.SetNotionToken(req.AccessToken, req.WorkspaceID, req.BotID, req.TokenExpiry, uc.clock)

		// Save updated user
		updatedUser, err := uc.repo.Update(ctx, user)
		if err != nil {
			return err
		}

		response = UpdateNotionTokenResponse{User: updatedUser}
		return nil
	})

	if err != nil {
		return UpdateNotionTokenResponse{}, err
	}

	return response, nil
}
