package app

import (
	"myapp/handler"
	"myapp/helpers/utils"
	"myapp/repository"
	partnerrepository "myapp/repository/partnerRepository"
	rolepermissionrepository "myapp/repository/rolePermissionRepository"
	userrepository "myapp/repository/userRepository"
	"myapp/service"
	partnerservice "myapp/service/partnerService"
	userservice "myapp/service/userService"
)

func SetupApp(repo repository.Repository) handler.Handler {

	userRepo := userrepository.NewUserRepository(repo)
	utils := utils.NewUtilsService(repo)
	rolePermissionRepo := rolepermissionrepository.NewPermissionRepository(repo)
	partnerRepo := partnerrepository.NewPartnerRepository(repo)

	service := service.NewService(utils, userRepo, rolePermissionRepo, partnerRepo)

	userService := userservice.NewUserService(service)
	partnerService := partnerservice.NewPartnerService(service)

	handler := handler.NewHandler(utils, userService, partnerService)

	return handler
}
