package repository

import (
	"calendly/app/models"
	"context"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(context.Context, *gorm.DB, *models.User) error
}

func NewUserRepository() UserRepositoryInterface {
	return &userRepository{}
}

type userRepository struct{}

func (repo *userRepository) Create(ctx context.Context, db *gorm.DB, user *models.User) error {
	return db.WithContext(ctx).Create(user).Error
}
