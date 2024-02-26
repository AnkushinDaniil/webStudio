package service

import (
	"main.go/internal/entity"
	"main.go/internal/repository/postgres"
)

type TimeslotListRepository interface {
	Create(userID int, list entity.TimeslotsList) (int, error)
	GetAll(userID int) ([]entity.TimeslotsList, error)
	GetByID(userID, listID int) (entity.TimeslotsList, error)
	Delete(userID, listID int) error
	Update(userID, listID int, input entity.UpdateListInput) error
}

type TimeslotListService struct {
	repo TimeslotListRepository
}

func NewTimeslotListService(repo postgres.TimeslotList) *TimeslotListService {
	return &TimeslotListService{repo: repo}
}

func (s *TimeslotListService) Create(userID int, list entity.TimeslotsList) (int, error) {
	return s.repo.Create(userID, list)
}

func (s *TimeslotListService) GetAll(userID int) ([]entity.TimeslotsList, error) {
	return s.repo.GetAll(userID)
}

func (s *TimeslotListService) GetByID(userID, listID int) (entity.TimeslotsList, error) {
	return s.repo.GetByID(userID, listID)
}

func (s *TimeslotListService) Update(userID, listID int, input entity.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userID, listID, input)
}

func (s *TimeslotListService) Delete(userID, listID int) error {
	return s.repo.Delete(userID, listID)
}
