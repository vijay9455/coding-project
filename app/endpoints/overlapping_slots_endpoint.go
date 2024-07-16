package endpoints

import (
	"calendly/app/middleware"
	"calendly/app/repository"
	"calendly/app/services"
	"calendly/lib/db"
	"calendly/lib/logger"
	"calendly/lib/web"
	"errors"

	"gorm.io/gorm"
)

func OverlappingSlots(request *web.Request) (*web.JSONResponse, web.ErrorInterface) {
	currentUser := middleware.GetCurrentUser(request.Context())
	if currentUser == nil {
		logger.Error(request.Context(), "should not have reached here. current_user missing in authenticated endpoint", nil)
		return nil, web.ErrInternalServerError("something went wrong")
	}

	var params services.OverlappingSlotsParams
	err := request.ValidateQueryParamsToStruct(&params)
	if err != nil {
		return nil, web.ErrValidationFailed(err.Error())
	}

	otherUser, err := repository.NewUserRepository().GetByEmail(request.Context(), db.Get(), *params.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, web.ErrValidationFailed("unable to find user with given email")
		}
		return nil, web.ErrInternalServerError("something went wrong")
	}

	slots, err := services.NewOverlappingSlotFetcher().Fetch(request.Context(), currentUser, otherUser, *params.Date)
	if err != nil {
		return nil, web.ErrInternalServerError("something went wrong")
	}

	response := web.JSONResponse(map[string]any{"slots": slots})
	return &response, nil
}
