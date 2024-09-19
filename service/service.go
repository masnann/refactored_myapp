package service

import (
	"myapp/helpers/utils"
	"myapp/repository"
)

type Service struct {
	UserRepo repository.UserRepositoryInterface
	Utils    utils.UtilsInterface
}

func NewService(
	userRepo repository.UserRepositoryInterface,
	utils utils.UtilsInterface,
) Service {
	return Service{
		UserRepo: userRepo,
		Utils:    utils,
	}
}
