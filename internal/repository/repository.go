package repository

import (
	"github.com/jmoiron/sqlx"
	"main.go/internal/entity"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GetUser(username, password string) (entity.User, error)
}

type TimeslotList interface{}

type TimeslotItem interface{}

type Repository struct {
	Authorization
	TimeslotList
	TimeslotItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
