package services

import (
	"calendly/app/models"
	"calendly/app/repository"
	"calendly/lib/db"
	"calendly/lib/lock"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type AvailabilityUpdateInterface interface {
	Update(ctx context.Context, user *models.User, params *AvailabilityUpdateParams) ([]*models.UserAvailability, error)
}

type availabilityUpdateService struct {
	uaRepo repository.UserAvailabilityRepositoryInterface
}

type AvailabilityUpdateParams struct {
	DayOfWeek       int64 `json:"day_of_week" validate:"required,gte=1,lte=7"`
	MarkUnAvailable *bool `json:"mark_unavailable" validate:"required"`
	Availabilities  []*struct {
		StartTime *models.TimeOnly `json:"start_time" validate:"required"`
		EndTime   *models.TimeOnly `json:"end_time" validate:"required"`
	} `json:"availabilities"`
}

func NewAvailabilityUpdateService() AvailabilityUpdateInterface {
	return &availabilityUpdateService{
		uaRepo: repository.NewUserAvailabilityRepository(),
	}
}

func (svc *availabilityUpdateService) Update(ctx context.Context, user *models.User, params *AvailabilityUpdateParams) ([]*models.UserAvailability, error) {
	var availabilities []*models.UserAvailability
	err := lock.WithLock(ctx, db.Get(), fmt.Sprintf("%s:weekofDay:%d", user.ID, params.DayOfWeek), func(db *gorm.DB) error {
		return db.Transaction(func(tx *gorm.DB) error {
			var err error
			availabilities, err = svc.uaRepo.Delete(ctx, tx, user.ID, params.DayOfWeek)
			if err != nil {
				return err
			}

			if params.MarkUnAvailable != nil && *params.MarkUnAvailable {
				return nil
			}

			availabilities, err = svc.createAvailabilities(ctx, tx, user.ID, params)
			return err
		})
	})

	return availabilities, err
}

func (svc *availabilityUpdateService) createAvailabilities(ctx context.Context, db *gorm.DB, userID string, params *AvailabilityUpdateParams) ([]*models.UserAvailability, error) {
	availabilities := make([]*models.UserAvailability, len(params.Availabilities))
	for idx, availability := range params.Availabilities {
		availabilities[idx] = &models.UserAvailability{
			StartTime: *availability.StartTime,
			EndTime:   *availability.EndTime,
			DayOfWeek: params.DayOfWeek,
			UserID:    &userID,
		}
	}

	if err := svc.uaRepo.BulkCreate(ctx, db, availabilities); err != nil {
		return nil, err
	}

	return availabilities, nil
}
