package app

import (
	"myapp/handler"
	"myapp/helpers/utils"
	"myapp/repository"
	rolepermissionrepository "myapp/repository/rolePermissionRepository"
	userrepository "myapp/repository/userRepository"
	"myapp/service"
	userservice "myapp/service/userService"
)

func SetupApp(repo repository.Repository) handler.Handler {

	userRepo := userrepository.NewUserRepository(repo)
	utils := utils.NewUtilsService(repo)
	rolePermissionRepo := rolepermissionrepository.NewPermissionRepository(repo)

	service := service.NewService(utils, userRepo, rolePermissionRepo)

	userService := userservice.NewUserService(service)

	handler := handler.NewHandler(userService)

	return handler
}
