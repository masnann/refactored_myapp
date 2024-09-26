package routes

import (
	"myapp/handler"
	userhandler "myapp/handler/userHandler"
	"myapp/helpers/middlewares"

	"github.com/labstack/echo/v4"
)

func ApiRoutes(e *echo.Echo, handler handler.Handler) {

	public := e.Group("/api/v1/public")
	userHandler := userhandler.NewUserHandler(handler)

	private := e.Group("/api/v1/private")
	private.Use(middlewares.JWTMiddleware)

	private.POST("/findbyid", middlewares.SuperAdminMiddleware(userHandler.FindUserByID))

	userGroup := public.Group("/user")
	userGroup.POST("/register", userHandler.Register)
	userGroup.POST("/login", userHandler.Login)
	userGroup.POST("/delete", userHandler.DeleteUser)
}
