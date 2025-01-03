package service

import "myapp/models"

type UserServiceInterface interface {
	FindUserByID(req models.RequestID) (models.UserModels, error)
	Register(req models.UserRegisterRequest) (int64, error)
	DeleteUser(req models.RequestID) (int64, error)
	FindUserByEmail(req models.UserFindUserByEmailRequest) (models.UserModels, error)
	Login(req models.UserLoginRequest) (models.UserLoginResponse, error)
	FindProfile(userID int64) (models.UserModels, error)
}

type PartnerServiceInterface interface {
	PartnerCreate(req models.PartnerCreateRequest) (int64, error)
	PartnerAssignKey(req models.PartnerAssignKeyRequest) (int64, error)
	PartnerFindByPartnerIDandSecretKey(partnerID int64, secretKey string) (models.PartnerModels, error)
}
