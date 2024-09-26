package repository

import (
	"database/sql"
	"errors"
	"myapp/helpers"
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
            FROM users WHERE id = ? AND status = 'active'
        `)

		query = helpers.ReplaceSQL(query, "?")
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
            FROM users WHERE id = ? AND status = 'active'
        `)

		query = helpers.ReplaceSQL(query, "?")
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
            FROM users WHERE id = ? AND status = 'active'
        `)

		query = helpers.ReplaceSQL(query, "?")
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

func TestRepository_DeleteUser(t *testing.T) {
	ts := setup.SetupTestCaseRepository(t)

	t.Run("Failure: Database Error", func(t *testing.T) {
		// Construct the expected SQL query with the placeholder that will be replaced by helpers.
		query := regexp.QuoteMeta(`
			UPDATE users SET status = 'inactive' WHERE id = ? RETURNING id
		`)
		query = helpers.ReplaceSQL(query, "?")
		// Mock the expected query with a custom error.
		expectedErr := errors.New("database error")
		ts.Mock.ExpectQuery(query).WithArgs(3).WillReturnError(expectedErr)

		// Call the DeleteUser method with the user ID of 3.
		userID, err := ts.UserRepo.DeleteUser(3)

		// Assertions to check for the correct error and user ID returned.
		assert.Error(t, err)
		assert.Equal(t, int64(3), userID)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, ts.Mock.ExpectationsWereMet())
	})

	t.Run("Success: User Deactivated", func(t *testing.T) {
		// Construct the expected SQL query.
		query := regexp.QuoteMeta(`
			UPDATE users SET status = 'inactive' WHERE id = ? RETURNING id
		`)
		query = helpers.ReplaceSQL(query, "?")
		// Mock the expected query with the returning ID.
		expectedID := int64(3)
		ts.Mock.ExpectQuery(query).WithArgs(expectedID).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(expectedID),
		)

		// Call the DeleteUser method with the user ID of 3.
		userID, err := ts.UserRepo.DeleteUser(3)

		// Assertions to check for success.
		assert.NoError(t, err)
		assert.Equal(t, expectedID, userID)
		assert.NoError(t, ts.Mock.ExpectationsWereMet())
	})
}

func TestRepository_Register(t *testing.T) {
	ts := setup.SetupTestCaseRepository(t)
	req := models.UserModels{
		Username:  "testuser",
		Email:     "testuser@example.com",
		Password:  "hashedpassword",
		Status:    "active",
		CreatedAt: helpers.TimeStampNow(),
		UpdatedAt: "",
	}

	t.Run("Success: User Registered", func(t *testing.T) {
		// Expected query for the insertion
		query := `
			INSERT INTO users (username, email, password, status, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`

		// Mock the expected query with the returning ID.
		expectedID := int64(1)

		// Mocking the query execution and result
		ts.Mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(
			req.Username, req.Email, req.Password, req.Status, req.CreatedAt, req.UpdatedAt,
		).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(expectedID),
		)

		// Call the Register method
		ID, err := ts.UserRepo.Register(req)

		// Assertions to verify that the method returned the correct ID and no errors
		assert.NoError(t, err)
		assert.Equal(t, expectedID, ID)
		assert.NoError(t, ts.Mock.ExpectationsWereMet())
	})

	t.Run("Failure: Database Error", func(t *testing.T) {
		// Expected query for the insertion
		query := `
			INSERT INTO users (username, email, password, status, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`

		// Mock the expected query with a custom error.
		expectedErr := errors.New("database error")

		// Simulate a database error when executing the query
		ts.Mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(
			req.Username, req.Email, req.Password, req.Status, req.CreatedAt, req.UpdatedAt,
		).WillReturnError(expectedErr)

		// Call the Register method
		ID, err := ts.UserRepo.Register(req)

		// Assertions to verify the error handling
		assert.Error(t, err)
		assert.Equal(t, int64(0), ID)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, ts.Mock.ExpectationsWereMet())
	})
}

func TestRepository_FindUserByEmail(t *testing.T) {
	ts := setup.SetupTestCaseRepository(t)

	req := "test@example.com"

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

	query := regexp.QuoteMeta(`
            SELECT id, username, email, password, status, created_at, updated_at
            FROM users WHERE email = ?
        `)
	query = helpers.ReplaceSQL(query, "?")
	t.Run("Success Case - User Found", func(t *testing.T) {
		// Mock the expected query with arguments and corresponding rows.
		rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "status", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email, expectedUser.Password, expectedUser.Status, expectedUser.CreatedAt, expectedUser.UpdatedAt)
		ts.Mock.ExpectQuery(query).WithArgs(expectedUser.Email).WillReturnRows(rows)

		// Call the FindUserByID function and assert the results.
		user, err := ts.UserRepo.FindUserByEmail(expectedUser.Email)
		assert.NoError(t, err, "Error finding user")
		assert.Equal(t, expectedUser, user, "Retrieved user does not match expected user")

		// Verify that all mock expectations were met.
		assert.NoError(t, ts.Mock.ExpectationsWereMet(), "Mock expectations not met")
	})

	t.Run("Failure Case _ User Not Found", func(t *testing.T) {

		ts.Mock.ExpectQuery(query).WithArgs(req).WillReturnError(sql.ErrNoRows)

		// Call the FindUserByID function and assert the results.
		user, err := ts.UserRepo.FindUserByEmail(req)
		assert.Error(t, err, "Expected error finding user")
		assert.Equal(t, models.UserModels{}, user, "Retrieved user should be empty")
		assert.Equal(t, sql.ErrNoRows, err, "Expected specific error (sql.ErrNoRows)")

		// Verify that all mock expectations were met.
		assert.NoError(t, ts.Mock.ExpectationsWereMet(), "Mock expectations not met")
	})

}
