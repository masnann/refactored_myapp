package service

import "myapp/models"

type UserServiceInterface interface {
	FindUserByID(req models.RequestID) (models.UserModels, error)
	Register(req models.UserRegisterRequest) (int64, error)
	DeleteUser(req models.RequestID) (int64, error)
	FindUserByEmail(req models.UserFindUserByEmailRequest) (models.UserModels, error)
	Login(req models.UserLoginRequest) (models.UserLoginResponse, error)
}
