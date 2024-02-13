package service

import (
	"main.go/internal/entity"
	"main.go/internal/repository"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
}

type TimeslotList interface{}

type TimeslotItem interface{}

type Service struct {
	Authorization
	TimeslotList
	TimeslotItem
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}
