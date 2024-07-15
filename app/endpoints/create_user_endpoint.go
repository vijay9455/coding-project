package endpoints

import (
	"calendly/app/models"
	"calendly/app/repository"
	"calendly/app/services"
	"calendly/lib/logger"
	"calendly/lib/web"
	"context"
	"errors"
	"fmt"
)

func CreateUser(request *web.Request) (*web.JSONResponse, web.ErrorInterface) {
	var params services.UserCreateParams
	err := request.ValidateBodyToStruct(&params)
	if err != nil {
		logger.Error(request.Context(), "error validation", map[string]any{"error": err})
		return nil, web.ErrValidationFailed(err.Error())
	}

	userSvc := services.NewUserService()
	user, err := userSvc.Create(request.Context(), &params)
	if err != nil {
		var uniqueConstrainErr *repository.UniqueConstrainError
		if errors.As(err, &uniqueConstrainErr) {
			return nil, web.ErrValidationFailed(fmt.Sprintf("%s is already taken", uniqueConstrainErr.Columns()))
		}

		return nil, web.ErrInternalServerError(err.Error())
	}

	response := web.JSONResponse(buildUser(request.Context(), user, user.Availabilities))
	return &response, nil
}

func buildUser(ctx context.Context, user *models.User, availabilities []*models.UserAvailability) map[string]any {
	return map[string]any{
		"user": map[string]any{
			"id":             user.ID,
			"email":          user.Email,
			"availabilities": buildAvailabilities(ctx, availabilities),
		},
	}
}

func buildAvailabilities(_ context.Context, availabilities []*models.UserAvailability) []map[string]any {
	availabilitiesList := make([]map[string]any, len(availabilities))
	for idx, availability := range availabilities {
		availabilitiesList[idx] = map[string]any{
			"start_time":  availability.StartTime,
			"end_time":    availability.EndTime,
			"day_of_week": availability.DayOfWeek,
		}
	}
	return availabilitiesList
}
