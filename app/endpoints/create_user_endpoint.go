package endpoints

import (
	"calendly/app/services"
	"calendly/lib/logger"
	"calendly/lib/web"
)

func CreateUser(request *web.Request) (*web.JSONResponse, web.ErrorInterface) {
	var params services.UserCreateParams
	err := request.ValidateBodyToStruct(&params)
	logger.Info(request.Context(), "data", map[string]any{"err": err, "params": params})
	if err != nil {
		logger.Error(request.Context(), "error validation", map[string]any{"error": err})
		return nil, web.ErrValidationFailed(err.Error())
	}

	userSvc := services.NewUserService()
	user, err := userSvc.Create(request.Context(), &params)
	if err != nil {
		return nil, web.ErrInternalServerError(err.Error())
	}

	logger.Info(request.Context(), "user", map[string]any{"user": user})
	response := web.JSONResponse(map[string]any{"id": user.ID, "email": user.Email})
	return &response, nil
}
