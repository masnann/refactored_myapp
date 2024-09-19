package repository

import "myapp/models"

type UserRepositoryInterface interface {
	FindUserByID(id int64) (models.UserModels, error)
}
