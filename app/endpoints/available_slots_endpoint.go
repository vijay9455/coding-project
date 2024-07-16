package endpoints

import (
	"calendly/app/middleware"
	"calendly/app/services"
	"calendly/lib/logger"
	"calendly/lib/web"
)

func AvailableSlots(request *web.Request) (*web.JSONResponse, web.ErrorInterface) {
	currentUser := middleware.GetCurrentUser(request.Context())
	if currentUser == nil {
		logger.Error(request.Context(), "should not have reached here. current_user missing in authenticated endpoint", nil)
		return nil, web.ErrInternalServerError("something went wrong")
	}

	var params services.AvailableSlotsParams
	err := request.ValidateQueryParamsToStruct(&params)
	if err != nil {
		return nil, web.ErrValidationFailed(err.Error())
	}

	svc := services.NewAvailableSlotFetcherService()
	availabilityData, err := svc.Fetch(request.Context(), currentUser, *params.StartDate, *params.EndDate)
	if err != nil {
		logger.Error(request.Context(), "error while fetching available slots", map[string]any{"error": err})
		return nil, web.ErrInternalServerError("something went wrong")
	}

	response := web.JSONResponse(map[string]any{"availability_data": availabilityData})
	return &response, nil
}
