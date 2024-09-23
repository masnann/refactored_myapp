package setup

import (
	"database/sql"
	"fmt"
	"myapp/config"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

// delete table
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

// Function to initialize database connection
func InitializeTestDB(t *testing.T) *sql.DB {
	if err := OpenConnectionDBTest(); err != nil {
		t.Fatalf("Failed to open test database connection: %v", err)
	}

	db := GetTestDB()

	return db
}

// OpenConnectionSQLMock creates a mock database connection for testing purposes.
func OpenConnectionSQLMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Error creating mock database:", err)
	}

	return db, mock
}
