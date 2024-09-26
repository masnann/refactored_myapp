package service

import (
	"errors"
	"myapp/models"
	"myapp/test/setup"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test_FindByID tests the FindUserByID function of the UserService.
func TestService_FindByID(t *testing.T) {
	// Setup a mock test environment using setup.SetupTestCaseService.
	ts := setup.SetupTestCaseService(t)

	// Define the request ID and expected user for successful test case.
	req := models.RequestID{ID: 1}
	expectedUser := models.UserModels{
		ID:       1,
		Username: "Test User",
		Email:    "testuser@example.com",
	}

	// **User Not Found Test Case**

	t.Run("User Not Found", func(t *testing.T) {
		// Define the expected error for user not found scenario.
		expectedErr := errors.New("user not found")

		// Mock the user repository's FindUserByID method to return an empty user and the expected error.
		ts.UserRepo.On("FindUserByID", req.ID).Return(models.UserModels{}, expectedErr).Once()

		// Call the UserService's FindUserByID method and capture the result and error.
		result, err := ts.UserService.FindUserByID(req)

		// Verify that all mock expectations were met.
		ts.UserRepo.AssertExpectations(t)

		// Assert the error and returned user as expected.
		assert.Error(t, err)
		assert.Equal(t, models.UserModels{}, result)
		assert.Equal(t, expectedErr, err)
	})

	// **User Found Test Case**

	t.Run("User Found", func(t *testing.T) {
		// Mock the user repository's FindUserByID method to return the expected user and no error.
		ts.UserRepo.On("FindUserByID", req.ID).Return(expectedUser, nil).Once()

		// Call the UserService's FindUserByID method and capture the result and error.
		result, err := ts.UserService.FindUserByID(req)

		// Verify that all mock expectations were met.
		ts.UserRepo.AssertExpectations(t)

		// Assert no error and the returned user matches the expected user.
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})
}

// TestUnit_DeleteUser tests the DeleteUser function of the UserService.
func TestService_DeleteUser(t *testing.T) {
	// Setup a mock test environment for this specific test case.
	ts := setup.SetupTestCaseService(t)
	req := models.RequestID{ID: 1}
	expectedUser := models.UserModels{
		ID:       1,
		Username: "Test User",
		Email:    "testuser@example.com",
	}

	// **Failure Case - User Not Found**

	t.Run("Failure Case - User Not Found", func(t *testing.T) {
		// Define the expected error for user not found scenario.
		expectedErr := errors.New("user not found")

		// Mock the user repository's FindUserByID method to return an empty user and the expected error.
		ts.UserRepo.On("FindUserByID", req.ID).Return(models.UserModels{}, expectedErr).Once()

		// Call the UserService's DeleteUser method and capture the result and error.
		result, err := ts.UserService.DeleteUser(req)

		// Verify that all mock expectations were met.
		ts.UserRepo.AssertExpectations(t)

		// Assert the error and returned user count as expected.
		assert.Error(t, err)
		assert.Equal(t, int64(0), result)
		assert.Equal(t, expectedErr, err)
	})

	// **Failure Case - Failed To Delete User**

	t.Run("Failure Case - Failed To Delete User", func(t *testing.T) {

		expectedErr := errors.New("failed to delete user")

		ts.UserRepo.On("FindUserByID", req.ID).Return(expectedUser, nil).Once()
		ts.UserRepo.On("DeleteUser", expectedUser.ID).Return(int64(0), expectedErr).Once()

		result, err := ts.UserService.DeleteUser(req)

		ts.UserRepo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, int64(0), result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Success Case - Success To Delete User", func(t *testing.T) {
		expectedResult := int64(1)

		ts.UserRepo.On("FindUserByID", req.ID).Return(expectedUser, nil).Once()
		ts.UserRepo.On("DeleteUser", expectedUser.ID).Return(expectedResult, nil).Once()

		result, err := ts.UserService.DeleteUser(req)

		ts.UserRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		assert.Nil(t, err)
	})
}

func TestService_Register(t *testing.T) {

	req := models.UserRegisterRequest{
		Username: "Test User",
		Email:    "testuser@example.com",
		Password: "testpassword",
	}

	expectedUser := models.UserModels{
		ID:       1,
		Username: "Test User",
		Email:    "testuser@example.com",
	}
	ts := setup.SetupTestCaseService(t)

	t.Run("Failure Case - User Already Registered", func(t *testing.T) {

		expectedErr := errors.New("email already registered")

		ts.UserRepo.On("FindUserByEmail", req.Email).Return(expectedUser, nil).Once()

		result, err := ts.UserService.Register(req)

		ts.UserRepo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, int64(0), result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Failure Case - Error Hash Password", func(t *testing.T) {
		expectedErr := errors.New("failed to generate hash")
		expectedErrs := errors.New("email already registered")
		ts.UserRepo.On("FindUserByEmail", req.Email).Return(models.UserModels{}, expectedErrs).Once()
		ts.Utils.On("GenerateHash", req.Password).Return("", expectedErr).Once()

		_, err := ts.UserService.Register(req)

		ts.Utils.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Failure Case - Failed To Register User", func(t *testing.T) {

		expectedErr := errors.New("failed to register user")
		expectedErrs := errors.New("email already registered")

		ts.UserRepo.On("FindUserByEmail", req.Email).Return(models.UserModels{}, expectedErrs).Once()
		ts.Utils.On("GenerateHash", req.Password).Return("hashed_password", nil).Once()
		ts.UserRepo.On("Register", mock.Anything).Return(int64(0), expectedErr).Once()

		result, err := ts.UserService.Register(req)

		ts.Utils.AssertExpectations(t)
		ts.UserRepo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, int64(0), result)
	})

	t.Run("Success Case - Successfully Register User", func(t *testing.T) {

		expectedErrs := errors.New("email already registered")
		ts.UserRepo.On("FindUserByEmail", req.Email).Return(models.UserModels{}, expectedErrs).Once()
		ts.Utils.On("GenerateHash", req.Password).Return("hashed_password", nil).Once()
		ts.UserRepo.On("Register", mock.Anything).Return(expectedUser.ID, nil).Once()

		result, err := ts.UserService.Register(req)

		ts.Utils.AssertExpectations(t)
		ts.UserRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, result)
		assert.Nil(t, err)
	})
}

func TestService_FindUserByEmail(t *testing.T) {
	ts := setup.SetupTestCaseService(t)
	req := models.UserFindUserByEmailRequest{
		Email: "test@example.com",
	}

	t.Run("Failure Case - User Not Found", func(t *testing.T) {
		expectedErr := errors.New("user not found")

		ts.UserRepo.On("FindUserByEmail", req.Email).Return(models.UserModels{}, expectedErr).Once()

		result, err := ts.UserService.FindUserByEmail(req)

		ts.UserRepo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, models.UserModels{}, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Success Case - User Found", func(t *testing.T) {
		expectedUser := models.UserModels{
			ID:       1,
			Username: "Test User",
			Email:    "test@example.com",
		}

		ts.UserRepo.On("FindUserByEmail", req.Email).Return(expectedUser, nil).Once()

		result, err := ts.UserService.FindUserByEmail(req)

		ts.UserRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		assert.Nil(t, err)
	})
}
