package services

import (
	"calendly/app/models"
	"calendly/app/repository"
	"calendly/lib/db"
	"context"
	"time"

	"gorm.io/gorm"
)

var (
	defaultDayOfWeeks = []int64{1, 2, 3, 4, 5}

	defaultStartTime = time.Date(0, 0, 0, 9, 0, 0, 0, time.UTC)
	defaultEndTime   = time.Date(0, 0, 0, 17, 0, 0, 0, time.UTC)
)

type UserCreateServiceInterface interface {
	Create(context.Context, *UserCreateParams) (*models.User, error)
}

type userCreateService struct {
	userRepo repository.UserRepositoryInterface
	uaRepo   repository.UserAvailabilityRepositoryInterface
}

type UserCreateParams struct {
	FirstName *string `json:"first_name" validate:"required"`
	LastName  *string `json:"last_name" validate:"required"`
	Email     *string `json:"email" validate:"required,email"`
}

func NewUserService() UserCreateServiceInterface {
	return &userCreateService{
		userRepo: repository.NewUserRepository(),
		uaRepo:   repository.NewUserAvailabilityRepository(),
	}
}

func (svc *userCreateService) Create(ctx context.Context, params *UserCreateParams) (*models.User, error) {
	user := svc.buildUserModel(params)
	err := db.Get().Transaction(func(tx *gorm.DB) error {
		if err := svc.userRepo.Create(ctx, tx, user); err != nil {
			return err
		}

		availabilities := svc.buildDefaultAvailabilities(user)
		if err := svc.uaRepo.BulkCreate(ctx, tx, availabilities); err != nil {
			return err
		}

		user.Availabilities = availabilities
		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *userCreateService) buildUserModel(params *UserCreateParams) *models.User {
	return &models.User{
		FirstName: *params.FirstName,
		LastName:  *params.LastName,
		Email:     *params.Email,
	}
}

func (svc *userCreateService) buildDefaultAvailabilities(user *models.User) []*models.UserAvailability {
	// by default as part of creating user create availability for 9 AM - 5 PM on every week day
	var availabilities []*models.UserAvailability
	for _, dayOfWeek := range defaultDayOfWeeks {
		availabilities = append(availabilities, &models.UserAvailability{
			StartTime: models.TimeOnly{Time: defaultStartTime},
			EndTime:   models.TimeOnly{Time: defaultEndTime},
			DayOfWeek: dayOfWeek,
			UserID:    &user.ID,
		})
	}
	return availabilities
}
