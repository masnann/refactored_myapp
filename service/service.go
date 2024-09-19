package service

import "myapp/repository"

type Service struct {
	UserRepo repository.UserRepositoryInterface
}

func NewService(
	userRepo repository.UserRepositoryInterface,
) Service {
	return Service{
		UserRepo: userRepo,
	}
}
