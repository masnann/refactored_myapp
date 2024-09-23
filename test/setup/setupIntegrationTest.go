package setup

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"myapp/app"
	"myapp/models"
	"myapp/repository"
	"myapp/routes"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

// SetupEcho sets up the Echo instance for testing.
func SetupEcho(db *sql.DB) *echo.Echo {
	e := echo.New()
	repo := repository.NewRepository(db)
	handler := app.SetupApp(repo)
	routes.ApiRoutes(e, handler)
	return e
}

// General helper function to send HTTP requests and return the response
func SendAPIRequest(t *testing.T, e *echo.Echo, method, path string, payload interface{}) *httptest.ResponseRecorder {
	var reqBodyBytes []byte
	var err error

	// If a payload is provided, marshal it into JSON
	if payload != nil {
		reqBodyBytes, err = json.Marshal(payload)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	// Create the HTTP request with the provided method, path, and payload
	req := httptest.NewRequest(method, path, bytes.NewReader(reqBodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Record the response
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec
}

// General helper to parse the response body
func ParseResponseBody(t *testing.T, rec *httptest.ResponseRecorder) models.Response {
	var response models.Response
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	return response
}
