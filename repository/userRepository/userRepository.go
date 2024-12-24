package userrepository

import (
	"fmt"
	"log"
	"myapp/helpers"
	"myapp/models"
	"myapp/repository"
)

type UserRepository struct {
	repo repository.Repository
}

func NewUserRepository(repo repository.Repository) UserRepository {
	return UserRepository{
		repo: repo,
	}
}

const (
	UserTable         = "users"
	UserInsertColumns = "username, email, password, status, address, created_at, updated_at"
	UserSelectColumns = "id, username, email, password, status, created_at, updated_at, coalesce(address,'')"
)

// UserScanArgs returns a slice of pointers to the fields of the UserModels struct
func userScanArgs(user *models.UserModels) []interface{} {
	return []interface{}{
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Address,
	}
}

func (r UserRepository) FindUserByID(id int64) (models.UserModels, error) {
	var user models.UserModels
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE id = ? AND status = 'active'`, UserSelectColumns, UserTable)

	query = helpers.ReplaceSQL(query, "?")

	row := r.repo.DB.QueryRow(query, id)
	err := row.Scan(userScanArgs(&user)...)
	if err != nil {
		log.Println("Error query FindUserByID: ", err)
		return user, err
	}
	return user, nil
}

func (r UserRepository) Register(req models.UserModels) (int64, error) {
	var ID int64
	query := fmt.Sprintf(`
		INSERT INTO %s (%s) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id`, UserTable, UserInsertColumns)

	query = helpers.ReplaceSQL(query, "?")
	err := r.repo.DB.QueryRow(query, req.Username, req.Email, req.Password, req.Status, req.Address, req.CreatedAt, req.UpdatedAt).Scan(&ID)
	if err != nil {
		log.Println("Error querying register: ", err)
		return ID, err
	}

	return ID, nil
}

func (r UserRepository) DeleteUser(userID int64) (int64, error) {
	var updatedID int64
	query := fmt.Sprintf(`
		UPDATE %s 
		SET status = 'inactive' 
		WHERE id = ?
		RETURNING id`, UserTable)

	query = helpers.ReplaceSQL(query, "?")
	err := r.repo.DB.QueryRow(query, userID).Scan(&updatedID)
	if err != nil {
		log.Println("Error querying delete user: ", err)
		return userID, err
	}

	return userID, nil
}

func (r UserRepository) FindUserByEmail(email string) (models.UserModels, error) {
	var user models.UserModels
	query := fmt.Sprintf(`
		SELECT %s 
		FROM %s 
		WHERE email = ?`, UserSelectColumns, UserTable)

	query = helpers.ReplaceSQL(query, "?")

	row := r.repo.DB.QueryRow(query, email)
	err := row.Scan(userScanArgs(&user)...)
	if err != nil {
		log.Println("Error querying FindUserByEmail: ", err)
		return user, err
	}
	return user, nil
}
