package partnerhandler

import (
	"log"
	"myapp/constants"
	"myapp/handler"
	"myapp/helpers"
	"myapp/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PartnerHandler struct {
	handler handler.Handler
}

func NewPartnerHandler(handler handler.Handler) PartnerHandler {
	return PartnerHandler{
		handler: handler,
	}
}

func (h PartnerHandler) PartnerCreate(ctx echo.Context) error {
	var result models.Response

	req := new(models.PartnerCreateRequest)
	if err := helpers.ValidateStruct(ctx, req); err != nil {
		log.Printf("Error Failed to validate request: %v", err)
		result = helpers.ResponseJSON(false, constants.VALIDATION_ERROR_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}
	partnerID, err := h.handler.PartnerService.PartnerCreate(*req)
	if err != nil {
		log.Printf("Error CreatePartner: %v", err)
		result = helpers.ResponseJSON(false, constants.INTERNAL_SERVER_ERROR, err.Error(), nil)
		return ctx.JSON(http.StatusInternalServerError, result)
	}
	result = helpers.ResponseJSON(true, constants.SUCCESS_CODE, constants.EMPTY_VALUE, partnerID)
	return ctx.JSON(http.StatusCreated, result)
}

func (h PartnerHandler) PartnerAssignKey(ctx echo.Context) error {
	var result models.Response

	req := new(models.PartnerAssignKeyRequest)
	if err := helpers.ValidateStruct(ctx, req); err != nil {
		log.Printf("Error Failed to validate request: %v", err)
		result = helpers.ResponseJSON(false, constants.VALIDATION_ERROR_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}
	id, err := h.handler.PartnerService.PartnerAssignKey(*req)
	if err != nil {
		log.Printf("Error AssignKey: %v", err)
		result = helpers.ResponseJSON(false, constants.INTERNAL_SERVER_ERROR, err.Error(), nil)
		return ctx.JSON(http.StatusInternalServerError, result)
	}
	result = helpers.ResponseJSON(true, constants.SUCCESS_CODE, constants.EMPTY_VALUE, id)
	return ctx.JSON(http.StatusOK, result)
}
