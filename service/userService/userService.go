package userservice

import (
	"errors"
	"myapp/models"
	"myapp/service"
)

type UserService struct {
	service service.Service
}

func NewUserService(service service.Service) UserService {
	return UserService{
		service: service,
	}
}

func (s UserService) FindUserByID(req models.RequestID) (models.UserModels, error) {
	result, err := s.service.UserRepo.FindUserByID(req.ID)
	if err != nil {
		msg := "user not found"
		return result, errors.New(msg)
	}
	return result, nil
}
