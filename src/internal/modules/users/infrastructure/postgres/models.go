package postgres

import (
	"time"

	"src/internal/modules/users/domain"
)

// UserRecord represents the user table structure in PostgreSQL
type UserRecord struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)"`
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
	return domain.User{
		ID:        record.ID,
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
		ID:        user.ID,
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
