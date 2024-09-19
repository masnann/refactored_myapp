package userhandler

import (
	"log"
	"myapp/constants"
	"myapp/handler"
	"myapp/helpers"
	"myapp/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	handler handler.Handler
}

func NewUserHandler(handler handler.Handler) UserHandler {
	return UserHandler{
		handler: handler,
	}
}

func (h UserHandler) FindUserByID(ctx echo.Context) error {
	var result models.Response

	req := new(models.RequestID)
	if err := helpers.ValidateStruct(ctx, req); err != nil {
		log.Printf("Error Failed to validate request: %v", err)
		result = helpers.ResponseJSON(false, constants.VALIDATION_ERROR_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	user, err := h.handler.UserService.FindUserByID(*req)
	if err != nil {
		result = helpers.ResponseJSON(false, constants.INTERNAL_SERVER_ERROR, err.Error(), nil)
		return ctx.JSON(http.StatusInternalServerError, result)
	}

	result = helpers.ResponseJSON(true, constants.SUCCESS_CODE, constants.EMPTY_VALUE, user)
	return ctx.JSON(http.StatusOK, result)
}
