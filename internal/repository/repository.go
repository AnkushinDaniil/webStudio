package repository

import (
	"github.com/jmoiron/sqlx"
	"main.go/internal/entity"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GetUser(username, password string) (entity.User, error)
}

type TimeslotList interface {
	Create(userID int, list entity.TimeslotsList) (int, error)
	GetAll(userID int) ([]entity.TimeslotsList, error)
	GetByID(userID, listID int) (entity.TimeslotsList, error)
	Delete(userID, listID int) error
	Update(userID, listID int, input entity.UpdateListInput) error
}

type TimeslotItem interface {
	Create(listID int, item entity.TimeslotItem) (int, error)
	GetAll(userID, listID int) ([]entity.TimeslotItem, error)
}

type Repository struct {
	Authorization
	TimeslotList
	TimeslotItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TimeslotList:  NewTimeslotListPostgres(db),
		TimeslotItem:  NewTimeslotItemPostgres(db),
	}
}
