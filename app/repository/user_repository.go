package repository

import (
	"calendly/app/models"
	"context"
	"strings"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(context.Context, *gorm.DB, *models.User) error
	GetByID(context.Context, *gorm.DB, string) (*models.User, error)
	GetByEmail(context.Context, *gorm.DB, string) (*models.User, error)
}

func NewUserRepository() UserRepositoryInterface {
	return &userRepository{}
}

type userRepository struct{}

func (repo *userRepository) Create(ctx context.Context, db *gorm.DB, user *models.User) error {
	err := db.WithContext(ctx).Create(user).Error
	if err != nil && strings.Contains(err.Error(), "unique_idx_users_email") {
		return newUniqueConstrainError(err, "unique_idx_users_email", "email")
	}

	return err
}

func (repo *userRepository) GetByID(ctx context.Context, db *gorm.DB, userID string) (*models.User, error) {
	var user *models.User
	err := db.WithContext(ctx).Where("id=?", userID).First(&user).Error
	return user, err
}

func (repo *userRepository) GetByEmail(ctx context.Context, db *gorm.DB, email string) (*models.User, error) {
	var user *models.User
	err := db.WithContext(ctx).Where("email=?", email).First(&user).Error
	return user, err
}
