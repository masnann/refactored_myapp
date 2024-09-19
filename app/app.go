package app

import (
	"myapp/handler"
	"myapp/helpers/utils"
	"myapp/repository"
	userrepository "myapp/repository/userRepository"
	"myapp/service"
	userservice "myapp/service/userService"
)

func SetupApp(repo repository.Repository) handler.Handler {

	userRepo := userrepository.NewUserRepository(repo)
	utils := utils.NewUtilsService(repo)

	service := service.NewService(userRepo, utils)

	userService := userservice.NewUserService(service)

	handler := handler.NewHandler(userService)

	return handler
}
