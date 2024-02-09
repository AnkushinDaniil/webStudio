package service

import "main.go/internal/repository"

type Authorization interface{}

type TimeslotList interface{}

type TimeslotItem interface{}

type Service struct {
	Authorization
	TimeslotList
	TimeslotItem
}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}
