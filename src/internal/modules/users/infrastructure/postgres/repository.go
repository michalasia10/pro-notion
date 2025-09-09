package postgres

import (
	"context"

	"gorm.io/gorm"

	"src/internal/modules/users/domain"
)

// UserRepository implements domain.UserRepository using PostgreSQL/GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	record := toUserRecord(user)

	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.User{}, err
	}

	return toDomainUser(record), nil
}

// GetByID retrieves a user by ID from the database
func (r *UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	var record UserRecord

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return toDomainUser(record), nil
}

// GetByEmail retrieves a user by email from the database
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var record UserRecord

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return toDomainUser(record), nil
}

// Update updates an existing user in the database
func (r *UserRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	record := toUserRecord(user)

	err := r.db.WithContext(ctx).Save(&record).Error
	if err != nil {
		return domain.User{}, err
	}

	return toDomainUser(record), nil
}

// Delete removes a user from the database
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&UserRecord{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// List retrieves users with pagination from the database
func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]domain.User, error) {
	var records []UserRecord

	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&records).Error

	if err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(records))
	for _, record := range records {
		users = append(users, toDomainUser(record))
	}

	return users, nil
}
