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

    // Open a connection for testing (testDB)
    testDB, err = sql.Open("postgres", config.TestDBURL)
    if err != nil {
        return fmt.Errorf("failed to connect to test database: %w", err)
    }

    // Open a separate connection for cleanup (cleanupDB)
    cleanupDB, err := sql.Open("postgres", config.TestDBURL)
    if err != nil {
        return fmt.Errorf("failed to connect to cleanup database: %w", err)
    }
    defer cleanupDB.Close() // Close cleanupDB after use

    // Ping the testing database
    if err := testDB.Ping(); err != nil {
        return fmt.Errorf("failed to ping test database: %w", err)
    }

    if err := cleanupDatabase(cleanupDB); err != nil {
        return fmt.Errorf("failed to clean up test database: %w", err)
    }

    // Run migrations on testDB
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
	// Query to get all table
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_type = 'BASE TABLE';`

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to fetch table names: %w", err)
	}
	defer rows.Close()

	var tableName string
	var tables []string

	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	// Drop all table
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", table))
		if err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
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
