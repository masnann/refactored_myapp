package handler

import (
	"myapp/helpers/utils"
	"myapp/service"
)

type Handler struct {
	Utils          utils.UtilsInterface
	UserService    service.UserServiceInterface
	PartnerService service.PartnerServiceInterface
}

func NewHandler(
	utils utils.UtilsInterface,
	userService service.UserServiceInterface,
	partnerService service.PartnerServiceInterface,
) Handler {
	return Handler{
		Utils:          utils,
		UserService:    userService,
		PartnerService: partnerService,
	}
}
