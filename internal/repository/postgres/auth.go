package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"main.go/internal/entity"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GetUser(username, password string) (entity.User, error)
}

type AuthorizationPostgres struct {
	db *sqlx.DB
}

func NewAuthorizationPostgres(db *sqlx.DB) *AuthorizationPostgres {
	return &AuthorizationPostgres{db: db}
}

func (r *AuthorizationPostgres) CreateUser(user entity.User) (int, error) {
	var userID int

	query := fmt.Sprintf(
		`
			INSERT INTO %s (name, color, username, password_hash)
			    VALUES ($1, $2, $3, $4)
			RETURNING
			    id`,
		UsersTable,
	)
	row := r.db.QueryRow(query, user.Name, user.Color, user.Username, user.Password)

	if err := row.Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *AuthorizationPostgres) GetUser(username, password string) (entity.User, error) {
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
