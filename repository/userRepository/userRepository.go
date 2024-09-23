package userrepository

import (
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

func (r UserRepository) FindUserByID(id int64) (models.UserModels, error) {
	var user models.UserModels
	query := `
		SELECT 
			id, username, email, password, status, created_at, updated_at
		FROM 
			users WHERE id = ? AND status = 'active'`

	query = helpers.ReplaceSQL(query, "?")

	row := r.repo.DB.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("Error query FindUserByID: ", err)
		return user, err
	}
	return user, nil
}

func (r UserRepository) Register(req models.UserModels) (int64, error) {
	var ID int64
	query := `
		INSERT INTO users (username, email, password, status, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id`

	query = helpers.ReplaceSQL(query, "?")
	err := r.repo.DB.QueryRow(query, req.Username, req.Email, req.Password, req.Status, req.CreatedAt, req.UpdatedAt).Scan(&ID)
	if err != nil {
		log.Println("Error querying register: ", err)
		return ID, err
	}

	return ID, nil
}

func (r UserRepository) DeleteUser(userID int64) (int64, error) {
	var updatedID int64
	query := `UPDATE users SET status = 'inactive' WHERE id =? RETURNING id`
	query = helpers.ReplaceSQL(query, "?")
	err := r.repo.DB.QueryRow(query, userID).Scan(&updatedID)
	if err != nil {
		log.Println("Error querying register: ", err)
		return userID, err
	}

	return userID, nil
}
