package repository

import (
	"user-service/internal/database"
	"user-service/internal/models"
)

func CreateUser(
	email string,
	passwordHash string,
) error {

	_, err := database.DB.Exec(
		`
		INSERT INTO users(email, password_hash)
		VALUES($1, $2)
		`,
		email,
		passwordHash,
	)

	return err
}

func GetUserByEmail(
	email string,
) (*models.User, error) {

	user := &models.User{}

	err := database.DB.QueryRow(
		`
		SELECT id, email, password_hash, role
		FROM users
		WHERE email = $1
		`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
