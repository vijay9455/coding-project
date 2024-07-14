package services

import (
	"calendly/app/models"
	"calendly/app/repository"
	"calendly/lib/db"
	"context"
)

type UserCreateServiceInterface interface {
	Create(context.Context, *UserCreateParams) (*models.User, error)
}

type userCreateService struct {
	userRepo repository.UserRepositoryInterface
}

type UserCreateParams struct {
	FirstName *string `json:"first_name" validate:"required"`
	LastName  *string `json:"last_name" validate:"required"`
	Email     *string `json:"email" validate:"required,email"`
}

func NewUserService() UserCreateServiceInterface {
	return &userCreateService{userRepo: repository.NewUserRepository()}
}

func (svc *userCreateService) Create(ctx context.Context, params *UserCreateParams) (*models.User, error) {
	user := svc.buildUserModel(params)
	if err := svc.userRepo.Create(ctx, db.Get(), user); err != nil {
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
