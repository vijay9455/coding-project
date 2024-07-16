package v1

import (
	"calendly/app/endpoints"
	"calendly/app/middleware"

	"github.com/julienschmidt/httprouter"
)

func Init(router *httprouter.Router) {
	router.POST("/api/v1/users", middleware.ServeEndpoint(middleware.V1, endpoints.CreateUser))
	router.GET("/api/v1/my_profile", middleware.ServeEndpoint(middleware.V1, middleware.AuthenticateUser(endpoints.MyProfile)))
	router.PUT("/api/v1/users/availabilities", middleware.ServeEndpoint(middleware.V1, middleware.AuthenticateUser(endpoints.UpdateAvailability)))
	router.GET("/api/v1/available_slots", middleware.ServeEndpoint(middleware.V1, middleware.AuthenticateUser(endpoints.AvailableSlots)))
}
