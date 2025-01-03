package routes

import (
	"myapp/handler"
	partnerhandler "myapp/handler/partnerHandler"
	userhandler "myapp/handler/userHandler"
	"myapp/helpers/middlewares"

	"github.com/labstack/echo/v4"
)

func ApiRoutes(e *echo.Echo, handler handler.Handler) {

	public := e.Group("/api/v1/public")
	partner := e.Group("/api/partner/v0.1")
	partner.Use(middlewares.SignatureValidatorMiddleware(handler))
	userHandler := userhandler.NewUserHandler(handler)
	partnerHandler := partnerhandler.NewPartnerHandler(handler)

	private := e.Group("/api/v1/private")
	private.Use(middlewares.JWTMiddleware)

	//private.POST("/findbyid", middlewares.SuperAdminMiddleware(userHandler.FindUserByID))
	private.POST("/profile", userHandler.FindProfile)

	public.POST("/findbyid", userHandler.FindUserByID)

	userGroup := public.Group("/user")
	userGroup.POST("/register", userHandler.Register)
	userGroup.POST("/login", userHandler.Login)
	userGroup.POST("/delete", userHandler.DeleteUser)

	partnerGroup := public.Group("/partner")
	partnerGroup.POST("/create", partnerHandler.PartnerCreate)
	partnerGroup.POST("/assign-key", partnerHandler.PartnerAssignKey)

	partner.POST("/findbyid", userHandler.FindUserByIDWithSignature)
}
