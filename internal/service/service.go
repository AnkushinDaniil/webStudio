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
	Create(userId int, list entity.TimeslotList) (int, error)
}

type TimeslotItem interface{}

type Service struct {
	Authorization
	TimeslotList
	TimeslotItem
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		TimeslotList:  NewTimeslotListService(repo.TimeslotList),
	}
}
