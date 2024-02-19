package service

import (
	"main.go/internal/entity"
	"main.go/internal/repository"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TimeslotList interface {
	Create(userID int, list entity.TimeslotsList) (int, error)
	GetAll(userID int) ([]entity.TimeslotsList, error)
	GetByID(userID, listID int) (entity.TimeslotsList, error)
	Delete(userID, listID int) error
	Update(userID, listID int, input entity.UpdateListInput) error
}

type TimeslotItem interface {
	Create(userID, listID int, input entity.TimeslotItem) (int, error)
	GetAll(userID, listID int) ([]entity.TimeslotItem, error)
	GetByID(userID, itemID int) (entity.TimeslotItem, error)
	Delete(userID, itemID int) error
	Update(userID, itemID int, input entity.UpdateItemInput) error
}

type Service struct {
	Authorization
	TimeslotList
	TimeslotItem
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		TimeslotList:  NewTimeslotListService(repo.TimeslotList),
		TimeslotItem:  NewTimeslotItemService(repo.TimeslotItem, repo.TimeslotList),
	}
}
