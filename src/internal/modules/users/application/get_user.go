package application

import (
	"context"

	"src/internal/modules/users/domain"
)

// GetUserRequest contains the data needed to retrieve a user
type GetUserRequest struct {
	ID string
}

// GetUserResponse contains the retrieved user data
type GetUserResponse struct {
	User domain.User
}

// GetUserUseCase handles user retrieval business logic
type GetUserUseCase struct {
	repo domain.UserRepository
}

// NewGetUserUseCase creates a new GetUserUseCase
func NewGetUserUseCase(repo domain.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{repo: repo}
}

// Execute retrieves a user by ID
func (uc *GetUserUseCase) Execute(ctx context.Context, req GetUserRequest) (GetUserResponse, error) {
	if req.ID == "" {
		return GetUserResponse{}, domain.ErrInvalidUserID
	}

	user, err := uc.repo.GetByID(ctx, req.ID)
	if err != nil {
		return GetUserResponse{}, err
	}

	return GetUserResponse{User: user}, nil
}

// GetUserByEmailRequest contains the data needed to retrieve a user by email
type GetUserByEmailRequest struct {
	Email string
}

// GetUserByEmailUseCase handles user retrieval by email business logic
type GetUserByEmailUseCase struct {
	repo domain.UserRepository
}

// NewGetUserByEmailUseCase creates a new GetUserByEmailUseCase
func NewGetUserByEmailUseCase(repo domain.UserRepository) *GetUserByEmailUseCase {
	return &GetUserByEmailUseCase{repo: repo}
}

// Execute retrieves a user by email
func (uc *GetUserByEmailUseCase) Execute(ctx context.Context, req GetUserByEmailRequest) (GetUserResponse, error) {
	if req.Email == "" {
		return GetUserResponse{}, domain.ErrInvalidEmail
	}

	user, err := uc.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return GetUserResponse{}, err
	}

	return GetUserResponse{User: user}, nil
}
