package endpoints

import (
	"calendly/app/middleware"
	"calendly/app/services"
	"calendly/lib/logger"
	"calendly/lib/web"
)

func UpdateAvailability(request *web.Request) (*web.JSONResponse, web.ErrorInterface) {
	currentUser := middleware.GetCurrentUser(request.Context())
	if currentUser == nil {
		logger.Error(request.Context(), "should not have reached here. current_user missing in authenticated endpoint", nil)
		return nil, web.ErrInternalServerError("something went wrong")
	}

	var params services.AvailabilityUpdateParams
	err := request.ValidateBodyToStruct(&params)
	if err != nil {
		logger.Error(request.Context(), "error validation", map[string]any{"error": err})
		return nil, web.ErrValidationFailed(err.Error())
	}

	if params.MarkUnAvailable != nil && !*params.MarkUnAvailable && len(params.Availabilities) <= 0 {
		logger.Error(request.Context(), "invalid input, mark unavailable is false, but no availabilities sent", nil)
		return nil, web.ErrValidationFailed("availabilities is required when mark_unavailable is false")
	}

	updateSvc := services.NewAvailabilityUpdateService()
	availabilities, err := updateSvc.Update(request.Context(), currentUser, &params)
	if err != nil {
		logger.Error(request.Context(), "error while updating availabilities", map[string]any{"error": err, "params": params})
		return nil, web.ErrInternalServerError("something went wrong")
	}

	response := web.JSONResponse(map[string]any{"availabilities": buildAvailabilities(request.Context(), availabilities)})
	return &response, nil
}
