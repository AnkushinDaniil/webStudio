package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"main.go/internal/entity"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user entity.User) (int, error) {
	var userID int

	query := fmt.Sprintf(
		`
			INSERT INTO %s (name, username, password_hash)
			    VALUES ($1, $2, $3)
			RETURNING
			    id`,
		UsersTable,
	)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)

	if err := row.Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *AuthPostgres) GetUser(username, password string) (entity.User, error) {
	var user entity.User

	query := fmt.Sprintf(`
		SELECT
		    id
		FROM
		    %s
		WHERE
		    username = $1
		    AND password_hash = $2`, UsersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
