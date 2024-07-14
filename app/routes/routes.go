package routes

import (
	"calendly/app/middleware"
	v1 "calendly/app/routes/v1"
	"calendly/lib/web"

	"github.com/julienschmidt/httprouter"
)

func Init(router *httprouter.Router) {
	router.GET("/health-check", middleware.ServeEndpoint(middleware.V1, func(_ *web.Request) (*web.JSONResponse, web.ErrorInterface) {
		response := web.JSONResponse(map[string]any{"message": "Hello world! Health check success."})
		return &response, nil
	}))

	v1.Init(router)
}
