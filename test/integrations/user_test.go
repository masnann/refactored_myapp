package integrations

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"myapp/models"
	"myapp/test"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)



func TestFindUserByID(t *testing.T) {
	db := test.InitializeTestDB(t)
	defer db.Close() 

	e := test.SetupEcho(db)

	// Prepare a user for testing
	userID, err := insertTestUser(db, "testuser", "testuser@example.com", "password")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Create a request to find the user by ID
	reqBody := models.RequestID{ID: userID}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/public/user/findbyid", bytes.NewReader(reqBodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Response
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.True(t, response.Success)
	assert.Equal(t, "200", response.StatusCode) // Adjust according to your implementation
	assert.Equal(t, "testuser", response.Result.(map[string]interface{})["username"])
	assert.Equal(t, "testuser@example.com", response.Result.(map[string]interface{})["email"])
}

func insertTestUser(db *sql.DB, username, email, password string) (int64, error) {
	var id int64
	err := db.QueryRow(
		"INSERT INTO users (username, email, password, status, created_at, updated_at) VALUES ($1, $2, $3, 'active', NOW(), NOW()) RETURNING id",
		username, email, password,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
