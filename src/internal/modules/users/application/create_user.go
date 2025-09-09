package application

import (
	"context"

	shared "src/internal/modules/shared/domain"
	"src/internal/modules/users/domain"
)

// CreateUserRequest contains the data needed to create a new user
type CreateUserRequest struct {
	Email string
	Name  string
}

// CreateUserResponse contains the created user data
type CreateUserResponse struct {
	User domain.User
}

// CreateUserUseCase handles user creation business logic
type CreateUserUseCase struct {
	repo  domain.UserRepository
	idGen shared.IDGenerator
	clock shared.Clock
	txMgr shared.TransactionManager
}

// NewCreateUserUseCase creates a new CreateUserUseCase
func NewCreateUserUseCase(
	repo domain.UserRepository,
	idGen shared.IDGenerator,
	clock shared.Clock,
	txMgr shared.TransactionManager,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		repo:  repo,
		idGen: idGen,
		clock: clock,
		txMgr: txMgr,
	}
}

// Execute creates a new user
func (uc *CreateUserUseCase) Execute(ctx context.Context, req CreateUserRequest) (CreateUserResponse, error) {
	var response CreateUserResponse

	err := uc.txMgr.WithinTransaction(ctx, func(ctx context.Context) error {
		// Check if user already exists
		_, err := uc.repo.GetByEmail(ctx, req.Email)
		if err == nil {
			return domain.ErrInvalidEmail // User already exists
		}
		if err != domain.ErrUserNotFound {
			return err
		}

		// Generate new ID
		userID := uc.idGen.NewID("usr")

		// Create domain entity
		user, err := domain.NewUser(userID, req.Email, req.Name, uc.clock)
		if err != nil {
			return err
		}

		// Save to repository
		createdUser, err := uc.repo.Create(ctx, user)
		if err != nil {
			return err
		}

		response = CreateUserResponse{User: createdUser}
		return nil
	})

	if err != nil {
		return CreateUserResponse{}, err
	}

	return response, nil
}
