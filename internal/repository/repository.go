package repository

import (
	"github.com/jmoiron/sqlx"
	"main.go/internal/entity"
	"main.go/internal/repository/postgres"
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
	GetByID(userID, itemID int) (entity.TimeslotItem, error)
	Delete(userID, itemID int) error
	Update(userID, itemID int, input entity.UpdateItemInput) error
	GetByRange(input entity.ItemsByRange) ([]entity.TimeslotItem, error)
}

type Repository struct {
	Authorization
	TimeslotList
	TimeslotItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthorizationPostgres(db),
		TimeslotList:  postgres.NewTimeslotListPostgres(db),
		TimeslotItem:  postgres.NewTimeslotItemPostgres(db),
	}
}
