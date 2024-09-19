package userrepository

import (
	"database/sql"
	"errors"
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
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		log.Println("Error scanning row: ", err)
		return user, errors.New("error scanning row")
	}
	return user, nil
}
