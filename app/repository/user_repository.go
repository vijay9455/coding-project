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
	err := db.Transaction(func(tx *gorm.DB) error {
		availabilities := user.Availabilities
		user.Availabilities = nil
		if err := db.WithContext(ctx).Create(user).Error; err != nil {
			return err
		}

		return repo.createAvailabilities(ctx, tx, user, availabilities)
	})

	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "unique_idx_users_email") {
		return newUniqueConstrainError(err, "unique_idx_users_email", "email")
	}

	return err
}

func (repo *userRepository) createAvailabilities(ctx context.Context, db *gorm.DB, user *models.User, availabilities []*models.UserAvailability) error {
	if len(availabilities) > 0 {
		for _, availability := range availabilities {
			availability.UserID = &user.ID
		}

		if err := db.Create(availabilities).Error; err != nil {
			return err
		}

		user.Availabilities = availabilities
	}

	return nil
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
