package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"myapp/app"
	"myapp/config"
	"myapp/repository"
	"myapp/routes"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

var testDB *sql.DB

// OpenConnectionDBTest initializes the test database connection, runs migrations, and ensures that the schema is set up.
func OpenConnectionDBTest() error {
	var err error

	// Open a connection to the test database
	testDB, err = sql.Open("postgres", config.TestDBURL)
	if err != nil {
		return fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Ping the database to ensure the connection is established
	if err := testDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping test database: %w", err)
	}

	if err := cleanupDatabase(testDB); err != nil {
		return fmt.Errorf("failed to clean up test database: %w", err)
	}

	// Run database migrations to ensure schema is up to date
	if err := runMigrations(testDB); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// runMigrations applies all pending migrations to the database.
func runMigrations(db *sql.DB) error {
	migrationDir := filepath.Join("..", "..", "db", "migrations")

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	if err := goose.Up(db, migrationDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func cleanupDatabase(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users;")
	if err != nil {
		return fmt.Errorf("failed to clean up table users: %w", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS goose_db_version;")
	if err != nil {
		return fmt.Errorf("failed to clean up table goose_db_version: %w", err)
	}

	return nil
}

// GetTestDB returns the test database connection.
func GetTestDB() *sql.DB {
	return testDB
}

// SetupEcho sets up the Echo instance for testing.
func SetupEcho(db *sql.DB) *echo.Echo {
	e := echo.New()
	repo := repository.NewRepository(db)
	handler := app.SetupApp(repo)
	routes.ApiRoutes(e, handler)
	return e
}



func InitializeTestDB(t *testing.T) *sql.DB {
	if err := OpenConnectionDBTest(); err != nil {
		t.Fatalf("Failed to open test database connection: %v", err)
	}

	db := GetTestDB()

	return db
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
