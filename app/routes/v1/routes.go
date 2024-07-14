package v1

import (
	"calendly/app/endpoints"
	"calendly/app/middleware"

	"github.com/julienschmidt/httprouter"
)

func Init(router *httprouter.Router) {
	router.POST("/api/v1/users", middleware.ServeEndpoint(middleware.V1, endpoints.CreateUser))
}
