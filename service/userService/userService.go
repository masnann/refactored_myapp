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
	_, err := s.service.UserRepo.FindUserByEmail(req.Email)
	if err == nil {
		msg := "email already registered"
		return 0, errors.New(msg)
	}

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

	newRole := models.AssignRoleToUserRequest{
		UserID: result,
		RoleID: 2,
	}

	err = s.service.RolePermissionRepo.AssignRoleToUserRequest(newRole)
	if err != nil {
		msg := "failed to assign role"
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

// FindUserByEmail implements service.UserServiceInterface.
func (s UserService) FindUserByEmail(req models.UserFindUserByEmailRequest) (models.UserModels, error) {
	result, err := s.service.UserRepo.FindUserByEmail(req.Email)
	if err != nil {
		msg := "user not found"
		return result, errors.New(msg)
	}
	return result, nil
}

func (s UserService) Login(req models.UserLoginRequest) (models.UserLoginResponse, error) {
	var result models.UserLoginResponse

	user, err := s.service.UserRepo.FindUserByEmail(req.Email)
	if err != nil {
		log.Println("Error finding user by email: ", err)
		return result, errors.New("user not found")
	}

	isValidPassword, err := s.service.Utils.CompareHash(user.Password, req.Password)
	if !isValidPassword || err != nil {
		log.Println("Error comparing password: ", err)
		return result, errors.New("invalid password")
	}

	role, err := s.service.RolePermissionRepo.FindUserRole(user.ID)
	if err != nil {
		log.Println("Error finding user role: ", err)
		return result, errors.New("failed to find user role")
	}

	accessToken, err := s.service.Utils.GenerateJWT(user.ID, user.Email, role.RoleName)
	if err != nil {
		log.Println("Error generating JWT: ", err)
		return result, errors.New("failed to generate access token")
	}

	refreshToken, err := s.service.Utils.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Println("Error generating refresh token: ", err)
		return result, errors.New("failed to generate refresh token")
	}

	// permissions, err := s.service.RolePermissionRepo.FindPermissionsForUser(user.ID)
	// if err != nil {
	// 	log.Println("Error finding user permissions: ", err)
	// 	return result, errors.New("failed to find user permissions")
	// }

	response := models.UserLoginResponse{
		UserID:       user.ID,
		RoleName:     role.RoleName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		//Permission:   permissions,
	}

	return response, nil
}
