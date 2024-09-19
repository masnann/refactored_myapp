package app

import (
	"myapp/handler"
	"myapp/repository"
	userrepository "myapp/repository/userRepository"
	"myapp/service"
	userservice "myapp/service/userService"
)

func SetupApp(repo repository.Repository) handler.Handler {

	userRepo := userrepository.NewUserRepository(repo)

	service := service.NewService(userRepo)

	userService := userservice.NewUserService(service)

	handler := handler.NewHandler(userService)

	return handler
}
