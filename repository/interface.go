package repository

import "myapp/models"

type UserRepositoryInterface interface {
	FindUserByID(id int64) (models.UserModels, error)
	Register(req models.UserModels) (int64, error)
	DeleteUser(userID int64) (int64, error)
	FindUserByEmail(email string) (models.UserModels, error)
}

type RolePermissionRepositoryInterface interface {
	AssignRoleToUserRequest(req models.AssignRoleToUserRequest) error
	FindUserRole(userID int64) (models.FindUserRoleResponse, error)
}
