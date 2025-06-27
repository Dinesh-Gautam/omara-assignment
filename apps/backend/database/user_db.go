package database

import (
	"strategic-insight-analyst/backend/models"
)

func FindUserByID(uid string) (*models.User, error) {
	var user models.User
	query := "SELECT id, email, auth_method, created_at FROM users WHERE id = $1"
	err := DB.QueryRow(query, uid).Scan(&user.ID, &user.Email, &user.AuthMethod, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *models.User) error {
	insertQuery := "INSERT INTO users (id, email, auth_method, created_at) VALUES ($1, $2, $3, $4)"
	_, err := DB.Exec(insertQuery, user.ID, user.Email, user.AuthMethod, user.CreatedAt)
	return err
}
