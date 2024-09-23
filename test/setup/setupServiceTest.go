package setup

import (
	"myapp/mocks"
	"myapp/service"
	userservice "myapp/service/userService"
	"testing"
)

type TestSuiteService struct {
	// Define the mock repositories for testing
	UserRepo *mocks.UserRepositoryInterface
	Utils    *mocks.UtilsInterface

	// Define the service instances for testing
	Service     service.Service
	UserService userservice.UserService
}

func SetupTestCaseService(t *testing.T) *TestSuiteService {
	// Initialize the mock repositories instances for testing
	userRepo := mocks.NewUserRepositoryInterface(t)
	utils := mocks.NewUtilsInterface(t)

	// Initialize the service instances with the defined repository mocks
	svc := service.NewService(userRepo, utils)
	userService := userservice.NewUserService(svc)

	return &TestSuiteService{
		// Initialize the mock repositories instances for testing
		UserRepo: userRepo,
		Utils:    utils,

		// Initialize the service instances for testing
		Service:     svc,
		UserService: userService,
	}
}
