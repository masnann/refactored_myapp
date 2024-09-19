package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupCustomLogger(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `${time_custom} Status=${status}, method=${method}, uri=${uri}, user_agent="${user_agent}"` + "\n",
		CustomTimeFormat: "2006/01/02 15:04:05",
	}))
}
