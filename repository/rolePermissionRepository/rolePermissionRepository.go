package rolepermissionrepository

import (
	"errors"
	"log"
	"myapp/helpers"
	"myapp/models"
	"myapp/repository"
)

type RolePermissionRepository struct {
	repo repository.Repository
}

func NewPermissionRepository(repo repository.Repository) RolePermissionRepository {
	return RolePermissionRepository{
		repo: repo,
	}
}

func (r RolePermissionRepository) AssignRoleToUserRequest(req models.AssignRoleToUserRequest) error {

	query := `
        INSERT INTO user_role (user_id, role_id) 
        VALUES (?,?)
	`

	query = helpers.ReplaceSQL(query, "?")
	_, err := r.repo.DB.Exec(query, req.UserID, req.RoleID)
	if err != nil {
		log.Println("Error querying create user role: ", err)
		return errors.New("error query")
	}

	return nil
}

func (r RolePermissionRepository) FindUserRole(userID int64) (models.FindUserRoleResponse, error) {
	var userRole models.FindUserRoleResponse
	query := ` 
		SELECT
			u.username,
			u.email,
			r.id as role_id, 
			r.name as role_name
		FROM 
			users u
		JOIN 
			user_role ur ON u.id = ur.user_id
		JOIN 
			roles r ON ur.role_id = r.id 
		WHERE
			u.id = ?
		`
	query = helpers.ReplaceSQL(query, "?")
	rows, err := r.repo.DB.Query(query, userID)
	if err != nil {
		log.Println("Error querying find user role: ", err)
		return userRole, errors.New("error query")
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&userRole.Username, &userRole.Email, &userRole.RoleID, &userRole.RoleName)
		if err != nil {
			log.Println("Error scanning row: ", err)
			return userRole, errors.New("error scanning row")
		}
	}
	return userRole, nil
}
