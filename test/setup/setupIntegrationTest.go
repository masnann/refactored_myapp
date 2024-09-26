package setup

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"myapp/app"
	"myapp/config"
	"myapp/models"
	"myapp/repository"
	"myapp/routes"

	"github.com/golang-jwt/jwt"
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
func SendAPIRequest(t *testing.T, e *echo.Echo, method, path string, payload interface{}, token string) *httptest.ResponseRecorder {
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

	// If a token is provided, add it to the Authorization header
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Record the response
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec
}

func GenerateSuperAdminToken() string {
	claims := jwt.MapClaims{
		"userID": 1,
		"email":  "superadmin@example.com",
		"role":   "Customer",
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return ""
	}

	return accessToken
}

// General helper to parse the response body
func ParseResponseBody(t *testing.T, rec *httptest.ResponseRecorder) models.Response {
	var response models.Response
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	return response
}
