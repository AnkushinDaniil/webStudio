package service

import (
	"main.go/internal/entity"
	"main.go/internal/repository"
)

type TimeslotListService struct {
	repo repository.TimeslotList
}

func NewTimeslotListService(repo repository.TimeslotList) *TimeslotListService {
	return &TimeslotListService{repo: repo}
}

func (s *TimeslotListService) Create(userId int, list entity.TimeslotsList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TimeslotListService) GetAll(userId int) ([]entity.TimeslotsList, error) {
	return s.repo.GetAll(userId)
}

func (s *TimeslotListService) GetById(userId, listId int) (entity.TimeslotsList, error) {
	return s.repo.GetById(userId, listId)
}
