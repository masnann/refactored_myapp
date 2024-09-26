package integrations

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	"myapp/models"
	"myapp/test/setup"

	"github.com/stretchr/testify/assert"
)

func insertTestUser(db *sql.DB, username, email, password string) (int64, error) {
	var id int64
	err := db.QueryRow(
		"INSERT INTO users (username, email, password, status, created_at, updated_at) VALUES ($1, $2, $3, 'active', '', '') RETURNING id",
		username, email, password,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func TestFindUserByID(t *testing.T) {
	db := setup.InitializeTestDB(t)

	e := setup.SetupEcho(db)

	// Prepare a user for testing
	userID, err := insertTestUser(db, "testuser", "testuser@example.com", "password")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Define models response
	var response models.Response
	path := "/api/v1/private/findbyid"
	method := http.MethodPost
	token := setup.GenerateSuperAdminToken()

	// Test case for validation error
	t.Run("Failure Case - Error Validation", func(t *testing.T) {
		rec := setup.SendAPIRequest(t, e, method, path, nil, token)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		assert.False(t, response.Success)
		assert.Equal(t, "400", response.StatusCode)
		assert.Equal(t, "Validation error: Field 'ID' is required", response.Message)
	})

	// Test case for user not found
	t.Run("Failure Case - User Not Found", func(t *testing.T) {
		nonExistingUserID := userID + 9999
		reqBodyNegative := models.RequestID{ID: nonExistingUserID}
		rec := setup.SendAPIRequest(t, e, method, path, reqBodyNegative, token)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		assert.False(t, response.Success)
		assert.Equal(t, "500", response.StatusCode)
		assert.Equal(t, "user not found", response.Message)
	})

	// Test case for user found
	t.Run("Success Case - User Found", func(t *testing.T) {
		reqBody := models.RequestID{ID: userID}
		rec := setup.SendAPIRequest(t, e, method, path, reqBody, token)

		assert.Equal(t, http.StatusOK, rec.Code)

		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		assert.True(t, response.Success)
		assert.Equal(t, "200", response.StatusCode)
		assert.Equal(t, "testuser", response.Result.(map[string]interface{})["username"])
		assert.Equal(t, "testuser@example.com", response.Result.(map[string]interface{})["email"])
	})
}

func TestUserRegister(t *testing.T) {
	// Initialize the test database
	db := setup.InitializeTestDB(t)

	// Setup Echo context and routes
	e := setup.SetupEcho(db)
	// Prepare the request payload
	req := models.UserRegisterRequest{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	var response models.Response
	path := "/api/v1/public/user/register"
	method := http.MethodPost

	t.Run("Success Case - Register", func(t *testing.T) {

		// Create HTTP POST request to the user register endpoint
		rec := setup.SendAPIRequest(t, e, method, path, req, "")

		// Check the response status code
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Parse the response body

		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Assert the success of the registration process
		assert.True(t, response.Success)
		assert.Equal(t, "200", response.StatusCode)
		assert.NotNil(t, response.Result)
	})

	t.Run("Failure Case - Error Validation", func(t *testing.T) {
		req := models.UserRegisterRequest{
			Password: "password123",
		}
		// Create HTTP POST request to the user register endpoint
		rec := setup.SendAPIRequest(t, e, method, path, req, "")

		// Check the response status code for validation error
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response body
		var response models.Response
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Assert that validation failed
		assert.False(t, response.Success)
		assert.Equal(t, "400", response.StatusCode)
		assert.Contains(t, response.Message, "Validation error: Field 'Username' is required")
	})

	t.Run("Failure Case - Error Database", func(t *testing.T) {
		db.Close()

		// Create HTTP POST request to the user register endpoint
		rec := setup.SendAPIRequest(t, e, method, path, req, "")
		// Check the response status code (expecting 503 Service Unavailable)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response body
		var response models.Response
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Assert that the response indicates a database error
		assert.False(t, response.Success)
		assert.Equal(t, "500", response.StatusCode)
		assert.Contains(t, response.Message, "failed to register user")
	})

}

func TestDeleteUser(t *testing.T) {
	db := setup.InitializeTestDB(t)

	e := setup.SetupEcho(db)

	var response models.Response
	path := "/api/v1/public/user/delete"
	method := http.MethodPost

	userID, err := insertTestUser(db, "testuser", "testuser@example.com", "password")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	t.Run("Success Case - User Deleted", func(t *testing.T) {

		req := models.RequestID{ID: userID}
		rec := setup.SendAPIRequest(t, e, method, path, req, "")

		assert.Equal(t, http.StatusOK, rec.Code)

		response = setup.ParseResponseBody(t, rec)

		assert.True(t, response.Success)
		assert.Equal(t, "200", response.StatusCode)
		assert.Equal(t, "", response.Message)
	})

	t.Run("Failure Case - Error Validation", func(t *testing.T) {
		req := models.RequestID{}
		rec := setup.SendAPIRequest(t, e, method, path, req, "")

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		response = setup.ParseResponseBody(t, rec)

		assert.False(t, response.Success)
		assert.Equal(t, "400", response.StatusCode)
		assert.Contains(t, response.Message, "Validation error: Field 'ID' is required")
	})

	t.Run("Failure Case - User Not Found", func(t *testing.T) {
		nonExistingUserID := userID + 9999
		req := models.RequestID{ID: nonExistingUserID}
		rec := setup.SendAPIRequest(t, e, method, path, req, "")

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		response = setup.ParseResponseBody(t, rec)

		assert.False(t, response.Success)
		assert.Equal(t, "500", response.StatusCode)
		assert.Equal(t, "user not found", response.Message)
	})
}
