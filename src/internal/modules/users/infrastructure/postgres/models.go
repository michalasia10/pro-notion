package postgres

import (
	"time"

	"src/internal/modules/users/domain"

	"github.com/google/uuid"
)

// UserRecord represents the user table structure in PostgreSQL
type UserRecord struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid();index"` // Internal UUID for DB relations and ordering
	PublicID  string    `gorm:"uniqueIndex;type:varchar(255);index"`                  // Public ID with prefix for API
	Email     string    `gorm:"uniqueIndex;not null;type:varchar(255)"`
	Name      string    `gorm:"not null;type:varchar(255)"`
	CreatedAt time.Time `gorm:"not null;index"`
	UpdatedAt time.Time `gorm:"not null"`

	// Notion integration fields
	NotionAccessToken string     `gorm:"type:text"`
	NotionWorkspaceID string     `gorm:"type:varchar(255);index"`
	NotionBotID       string     `gorm:"type:varchar(255)"`
	NotionTokenExpiry *time.Time `gorm:""`
}

// TableName specifies the table name for GORM
func (UserRecord) TableName() string {
	return "users"
}

// toDomainUser converts a UserRecord to a domain User
func toDomainUser(record UserRecord) domain.User {
	// Parse the ID string back to UUID
	id, _ := uuid.Parse(record.ID)

	return domain.User{
		ID:        id,              // Internal UUID
		PublicID:  record.PublicID, // Public ID with prefix
		Email:     record.Email,
		Name:      record.Name,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,

		NotionAccessToken: record.NotionAccessToken,
		NotionWorkspaceID: record.NotionWorkspaceID,
		NotionBotID:       record.NotionBotID,
		NotionTokenExpiry: record.NotionTokenExpiry,
	}
}

// toUserRecord converts a domain User to a UserRecord
func toUserRecord(user domain.User) UserRecord {
	return UserRecord{
		ID:        user.ID.String(), // Convert UUID to string for DB storage
		PublicID:  user.PublicID,    // Public ID with prefix
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,

		NotionAccessToken: user.NotionAccessToken,
		NotionWorkspaceID: user.NotionWorkspaceID,
		NotionBotID:       user.NotionBotID,
		NotionTokenExpiry: user.NotionTokenExpiry,
	}
}
