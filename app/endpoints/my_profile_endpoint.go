package endpoints

import (
	"calendly/app/middleware"
	"calendly/app/repository"
	"calendly/lib/db"
	"calendly/lib/logger"
	"calendly/lib/web"
)

func MyProfile(request *web.Request) (*web.JSONResponse, web.ErrorInterface) {
	currentUser := middleware.GetCurrentUser(request.Context())

	if currentUser == nil {
		logger.Error(request.Context(), "should not have reached here. current_user missing in authenticated endpoint", nil)
		return nil, web.ErrInternalServerError("something went wrong")
	}

	uaRepo := repository.NewUserAvailabilityRepository()
	availabilities, err := uaRepo.GetByUserID(request.Context(), db.Get(), currentUser.ID)
	if err != nil {
		logger.Error(request.Context(), "error while fetching availabilities", map[string]any{"error": err})
		return nil, web.ErrInternalServerError("something went wrong")
	}

	response := web.JSONResponse(buildUser(request.Context(), currentUser, availabilities))
	return &response, nil
}
