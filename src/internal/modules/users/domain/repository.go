package domain

import "context"

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user User) (User, error)

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id string) (User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (User, error)

	// Update updates an existing user
	Update(ctx context.Context, user User) (User, error)

	// Delete removes a user by ID
	Delete(ctx context.Context, id string) error

	// List retrieves all users with pagination
	List(ctx context.Context, offset, limit int) ([]User, error)
}
