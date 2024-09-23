package userservice

import (
	"errors"
	"myapp/helpers"
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

func (s UserService) Register(req models.UserRegisterRequest) (int64, error) {

	hashedPassword, err := s.service.Utils.GenerateHash(req.Password)
	if err != nil {
		msg := "failed to generate hash"
		return 0, errors.New(msg)
	}

	newData := models.UserModels{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Status:    "",
		CreatedAt: helpers.TimeStampNow(),
		UpdatedAt: "",
	}

	result, err := s.service.UserRepo.Register(newData)
	if err != nil {
		msg := "failed to register user"
		return 0, errors.New(msg)
	}

	return result, nil
}

func (s UserService) DeleteUser(req models.RequestID) (int64, error) {
	user, err := s.service.UserRepo.FindUserByID(req.ID)
	if err != nil {
		msg := "user not found"
		return 0, errors.New(msg)
	}

	result, err := s.service.UserRepo.DeleteUser(user.ID)
	if err != nil {
		msg := "failed to delete user"
		return 0, errors.New(msg)
	}

	return result, nil
}
