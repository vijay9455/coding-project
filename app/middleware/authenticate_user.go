package middleware

import (
	"calendly/app/models"
	"calendly/app/repository"
	"calendly/lib/db"
	"calendly/lib/logger"
	"calendly/lib/web"
	"context"
	"errors"

	"gorm.io/gorm"
)

type currentUserCtx struct{}

const userIdHeader string = "calendly-user-id"

func AuthenticateUser(nextHandler web.Handle) web.Handle {
	return func(request *web.Request) (*web.JSONResponse, web.ErrorInterface) {
		// for now using user_id, actually we need to create and user some token or some other auth mechanist
		userID := request.Header.Get(userIdHeader)
		if userID == "" {
			return nil, web.ErrUnauthorised("user_id header missing")
		}

		user, err := repository.NewUserRepository().GetByID(request.Context(), db.Get(), userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, web.ErrUnauthorised("user not found")
			}

			logger.Error(request.Context(), "error while fetching user", map[string]any{"err": err})
			return nil, web.ErrInternalServerError("something went wrong.")
		}

		return nextHandler(request.WithContext(context.WithValue(request.Context(), currentUserCtx{}, user)))
	}
}

func GetCurrentUser(ctx context.Context) *models.User {
	currentUser := ctx.Value(currentUserCtx{})
	if currentUser != nil {
		return currentUser.(*models.User)
	}
	return nil
}
