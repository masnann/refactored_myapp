package service

import (
	"myapp/helpers/utils"
	"myapp/repository"
)

type Service struct {
	Utils              utils.UtilsInterface
	UserRepo           repository.UserRepositoryInterface
	RolePermissionRepo repository.RolePermissionRepositoryInterface
}

func NewService(
	utils utils.UtilsInterface,
	userRepo repository.UserRepositoryInterface,
	rolePermissionRepo repository.RolePermissionRepositoryInterface,
) Service {
	return Service{
		Utils:              utils,
		UserRepo:           userRepo,
		RolePermissionRepo: rolePermissionRepo,
	}
}
