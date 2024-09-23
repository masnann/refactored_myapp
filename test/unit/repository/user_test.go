package repository

import (
	"database/sql"
	"errors"
	"myapp/models"
	"myapp/test/setup"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestRepository_FindByID tests the FindUserByID function of the UserRepo.
func TestRepository_FindByID(t *testing.T) {
	// Setup a mock database connection for testing.
	ts := setup.SetupTestCaseRepository(t)

	// Define the expected user for successful test case.
	expectedUser := models.UserModels{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Status:    "active",
		CreatedAt: "",
		UpdatedAt: "",
	}

	// **Success Case - User Found**

	t.Run("Success: User Found", func(t *testing.T) {
		// Construct the expected SQL query with proper quoting.
		query := regexp.QuoteMeta(`
            SELECT id, username, email, password, status, created_at, updated_at
            FROM users WHERE id = $1 AND status = 'active'
        `)

		// Mock the expected query with arguments and corresponding rows.
		rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "status", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email, expectedUser.Password, expectedUser.Status, expectedUser.CreatedAt, expectedUser.UpdatedAt)
		ts.Mock.ExpectQuery(query).WithArgs(expectedUser.ID).WillReturnRows(rows)

		// Call the FindUserByID function and assert the results.
		user, err := ts.UserRepo.FindUserByID(expectedUser.ID)
		assert.NoError(t, err, "Error finding user")
		assert.Equal(t, expectedUser, user, "Retrieved user does not match expected user")

		// Verify that all mock expectations were met.
		assert.NoError(t, ts.Mock.ExpectationsWereMet(), "Mock expectations not met")
	})

	// **Failure Case - User Not Found**

	t.Run("Failure: User Not Found", func(t *testing.T) {
		// Construct the expected SQL query.
		query := regexp.QuoteMeta(`
            SELECT id, username, email, password, status, created_at, updated_at
            FROM users WHERE id = $1 AND status = 'active'
        `)

		// Mock the expected query with error (sql.ErrNoRows).
		ts.Mock.ExpectQuery(query).WithArgs(2).WillReturnError(sql.ErrNoRows)

		// Call the FindUserByID function and assert the results.
		user, err := ts.UserRepo.FindUserByID(2)
		assert.Error(t, err, "Expected error finding user")
		assert.Equal(t, models.UserModels{}, user, "Retrieved user should be empty")
		assert.Equal(t, sql.ErrNoRows, err, "Expected specific error (sql.ErrNoRows)")

		// Verify that all mock expectations were met.
		assert.NoError(t, ts.Mock.ExpectationsWereMet(), "Mock expectations not met")
	})
	// **Failure Case - Database Error**

	t.Run("Failure: Database Error", func(t *testing.T) {
		// Construct the expected SQL query.
		query := regexp.QuoteMeta(`
            SELECT id, username, email, password, status, created_at, updated_at
            FROM users WHERE id = $1 AND status = 'active'
        `)

		// Mock the expected query with a custom error.
		expectedErr := errors.New("database error")
		ts.Mock.ExpectQuery(query).WithArgs(3).WillReturnError(expectedErr)

		// Call the FindUserByID
		user, err := ts.UserRepo.FindUserByID(3)

		assert.Error(t, err)
		assert.Equal(t, models.UserModels{}, user)
		assert.Equal(t, expectedErr, err)

		assert.NoError(t, ts.Mock.ExpectationsWereMet())
	})
}
