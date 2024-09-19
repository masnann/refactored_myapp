package userservice

import (
	"errors"
	"log"
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
		log.Println("Error generating hash: ", err)
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
		log.Println("Error registering user: ", err)
		msg := "failed to register user"
		return 0, errors.New(msg)
	}

	return result, nil
}
