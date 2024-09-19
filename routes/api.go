package routes

import (
	"myapp/handler"
	userhandler "myapp/handler/userHandler"

	"github.com/labstack/echo/v4"
)

func ApiRoutes(e *echo.Echo, handler handler.Handler) {

	public := e.Group("/api/v1/public")
	userHandler := userhandler.NewUserHandler(handler)

	userGroup := public.Group("/user")
	userGroup.POST("/findbyid", userHandler.FindUserByID)
	userGroup.POST("/register", userHandler.Register)
}
