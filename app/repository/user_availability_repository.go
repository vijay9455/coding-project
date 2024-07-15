package repository

import (
	"calendly/app/models"
	"context"

	"gorm.io/gorm"
)

type UserAvailabilityRepositoryInterface interface {
	BulkCreate(context.Context, *gorm.DB, []*models.UserAvailability) error
	GetByUserID(ctx context.Context, db *gorm.DB, userID string) ([]*models.UserAvailability, error)
}

func NewUserAvailabilityRepository() UserAvailabilityRepositoryInterface {
	return &userAvailabilityRepository{}
}

type userAvailabilityRepository struct{}

func (uaRepo *userAvailabilityRepository) BulkCreate(ctx context.Context, db *gorm.DB, availabilities []*models.UserAvailability) error {
	return db.WithContext(ctx).Create(availabilities).Error
}

func (uaRepo *userAvailabilityRepository) GetByUserID(ctx context.Context, db *gorm.DB, userID string) ([]*models.UserAvailability, error) {
	var availabilities []*models.UserAvailability
	if err := db.WithContext(ctx).Find(&availabilities, "user_id=?", userID).Error; err != nil {
		return nil, err
	}

	return availabilities, nil
}
