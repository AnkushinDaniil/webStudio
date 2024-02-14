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

func (s *TimeslotListService) Create(userId int, list entity.TimeslotList) (int, error) {
	return s.repo.Create(userId, list)
}
