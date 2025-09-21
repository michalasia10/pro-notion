package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrUserNotFound       = errors.New("user not found")
	ErrNotionTokenMissing = errors.New("notion token missing")
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID
	PublicID  string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time

	// Notion integration
	NotionAccessToken string
	NotionWorkspaceID string
	NotionBotID       string
	NotionTokenExpiry *time.Time
}

// NewUser creates a new user with validation
func NewUser(email, name string, idGen IDGenerator, clock Clock) (User, error) {
	if email == "" {
		return User{}, ErrInvalidEmail
	}

	now := clock.Now()

	return User{
		ID:        uuid.New(),          // Internal UUID for DB relations and ordering
		PublicID:  idGen.NewID("user"), // Public ID with prefix for API
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// SetNotionToken updates the user's Notion access token and related information
func (u *User) SetNotionToken(accessToken, workspaceID, botID string, expiry *time.Time, clock Clock) {
	u.NotionAccessToken = accessToken
	u.NotionWorkspaceID = workspaceID
	u.NotionBotID = botID
	u.NotionTokenExpiry = expiry
	u.UpdatedAt = clock.Now()
}

// HasValidNotionToken checks if the user has a valid (non-expired) Notion token
func (u *User) HasValidNotionToken(clock Clock) bool {
	if u.NotionAccessToken == "" {
		return false
	}

	if u.NotionTokenExpiry != nil && clock.Now().After(*u.NotionTokenExpiry) {
		return false
	}

	return true
}

// Clock interface for dependency injection
type Clock interface {
	Now() time.Time
}

// IDGenerator provides unique ID generation for domain entities
type IDGenerator interface {
	NewID(prefix string) string
}
