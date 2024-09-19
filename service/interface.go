package service

import "myapp/models"

type UserServiceInterface interface {
	FindUserByID(req models.RequestID) (models.UserModels, error)
}
