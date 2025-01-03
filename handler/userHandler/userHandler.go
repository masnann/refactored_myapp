package userhandler

import (
	"encoding/json"
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

	// Ambil data request ID
	req := new(models.RequestID)
	if err := ctx.Bind(req); err != nil {
		log.Printf("Error Failed to bind request: %v", err)
		result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	// Ambil tanda tangan dari header (misalnya X-Signature)
	signature := ctx.Request().Header.Get("X-Signature")
	if signature == "" {
		result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, "Signature is required", nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	// Serialize request data secara konsisten
	dataBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error serializing request: %v", err)
		result = helpers.ResponseJSON(false, constants.INTERNAL_SERVER_ERROR, "Failed to serialize request data", nil)
		return ctx.JSON(http.StatusInternalServerError, result)
	}

	// Ambil public key
	publicKey, err := helpers.LoadPublicKey("public_key.pem")
	if err != nil {
		log.Printf("Error loading public key: %v", err)
		result = helpers.ResponseJSON(false, constants.INTERNAL_SERVER_ERROR, "Failed to load public key", nil)
		return ctx.JSON(http.StatusInternalServerError, result)
	}

	// Verifikasi tanda tangan menggunakan public key
	valid, err := helpers.VerifySignature(string(dataBytes), signature, publicKey)
	if err != nil {
		log.Printf("Error verifying signature: %v", err)
		result = helpers.ResponseJSON(false, constants.FORBIDDEN_CODE, "Invalid signature", nil)
		return ctx.JSON(http.StatusForbidden, result)
	}

	if !valid {
		result = helpers.ResponseJSON(false, constants.FORBIDDEN_CODE, "Invalid signature", nil)
		return ctx.JSON(http.StatusForbidden, result)
	}

	// Lanjutkan untuk mencari user berdasarkan ID
	user, err := h.handler.UserService.FindUserByID(*req)
	if err != nil {
		result = helpers.ResponseJSON(false, constants.INTERNAL_SERVER_ERROR, err.Error(), nil)
		return ctx.JSON(http.StatusInternalServerError, result)
	}

	// Kirimkan hasil user
	result = helpers.ResponseJSON(true, constants.SUCCESS_CODE, constants.EMPTY_VALUE, user)
	return ctx.JSON(http.StatusOK, result)
}

func (h UserHandler) Register(ctx echo.Context) error {
	var result models.Response

	req := new(models.UserRegisterRequest)
	if err := helpers.ValidateStruct(ctx, req); err != nil {
		log.Printf("Error Failed to validate request: %v", err)
		result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}
	userID, err := h.handler.UserService.Register(*req)
	if err != nil {
		result = helpers.ResponseJSON(false, constants.INTERNAL_SERVER_ERROR, err.Error(), nil)
		return ctx.JSON(http.StatusInternalServerError, result)
	}
	result = helpers.ResponseJSON(true, constants.SUCCESS_CODE, constants.EMPTY_VALUE, userID)
	return ctx.JSON(http.StatusCreated, result)
}

func (h UserHandler) DeleteUser(ctx echo.Context) error {
	var result models.Response

	req := new(models.RequestID)
	if err := helpers.ValidateStruct(ctx, req); err != nil {
		log.Printf("Error Failed to validate request: %v", err)
		result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	userID, err := h.handler.UserService.DeleteUser(*req)
	if err != nil {
		result = helpers.ResponseJSON(false, constants.INTERNAL_SERVER_ERROR, err.Error(), nil)
		return ctx.JSON(http.StatusInternalServerError, result)
	}

	result = helpers.ResponseJSON(true, constants.SUCCESS_CODE, constants.EMPTY_VALUE, userID)
	return ctx.JSON(http.StatusOK, result)
}

func (h UserHandler) Login(ctx echo.Context) error {
	var result models.Response

	req := new(models.UserLoginRequest)
	if err := helpers.ValidateStruct(ctx, req); err != nil {
		log.Printf("Error Failed to validate request: %v", err)
		result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}
	response, err := h.handler.UserService.Login(*req)
	if err != nil {
		log.Printf("Error Login: %v", err)
		result = helpers.ResponseJSON(false, constants.INTERNAL_SERVER_ERROR, err.Error(), nil)
		return ctx.JSON(http.StatusInternalServerError, result)
	}
	result = helpers.ResponseJSON(true, constants.SUCCESS_CODE, constants.EMPTY_VALUE, response)
	return ctx.JSON(http.StatusOK, result)
}

func (h UserHandler) FindProfile(ctx echo.Context) error {
	var result models.Response

	currentUser, err := helpers.GetCurrentUser(ctx)
	if err != nil {
		result := helpers.ResponseJSON(false, constants.FORBIDDEN_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusForbidden, result)
	}
	profile, err := h.handler.UserService.FindProfile(currentUser.ID)
	if err != nil {
		result = helpers.ResponseJSON(false, constants.INTERNAL_SERVER_ERROR, err.Error(), nil)
		return ctx.JSON(http.StatusInternalServerError, result)
	}
	result = helpers.ResponseJSON(true, constants.SUCCESS_CODE, constants.EMPTY_VALUE, profile)
	return ctx.JSON(http.StatusOK, result)
}

func (h UserHandler) FindUserByIDWithSignature(ctx echo.Context) error {
	var result models.Response

	// Ambil data request ID
	req := new(models.RequestID)
	if err := ctx.Bind(req); err != nil {
		log.Printf("Error Failed to bind request: %v", err)
		result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	// Lanjutkan untuk mencari user berdasarkan ID
	user, err := h.handler.UserService.FindUserByID(*req)
	if err != nil {
		result = helpers.ResponseJSON(false, constants.INTERNAL_SERVER_ERROR, err.Error(), nil)
		return ctx.JSON(http.StatusInternalServerError, result)
	}

	// Kirimkan hasil user
	result = helpers.ResponseJSON(true, constants.SUCCESS_CODE, constants.EMPTY_VALUE, user)
	return ctx.JSON(http.StatusOK, result)
}
