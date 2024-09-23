package setup

import (
	"testing"

	"myapp/repository"
	userrepository "myapp/repository/userRepository"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestSuiteRepository holds references to the mock database, repository, and user repository interface.
type TestSuiteRepository struct {
	Mock     sqlmock.Sqlmock
	Repo     repository.Repository
	UserRepo repository.UserRepositoryInterface
	
}

// SetupTestCaseRepository creates a mock database connection, instantiates repository objects, and returns a TestSuiteRepository.
func SetupTestCaseRepository(t *testing.T) *TestSuiteRepository {
	// Create a mock database connection for testing purposes.
	db, mock := OpenConnectionSQLMock(t)

	// Instantiate repository objects using the mocked database connection.
	repo := repository.NewRepository(db)
	userRepo := userrepository.NewUserRepository(repo)

	// Return a TestSuiteRepository containing the mock, repository, and repository interface.
	return &TestSuiteRepository{
		Mock:     mock,
		Repo:     repo,
		UserRepo: userRepo,
	}
}
