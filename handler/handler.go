package handler

import "myapp/service"

type Handler struct {
	UserService service.UserServiceInterface
}

func NewHandler(
	userService service.UserServiceInterface,
) Handler {
	return Handler{
		UserService: userService,
	}
}
