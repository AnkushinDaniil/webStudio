package service

import (
	"net/http"

	"main.go/internal/entity"
	"main.go/internal/repository"
)

type TimeslotItemService struct {
	itemRepo repository.TimeslotItem
	listRepo repository.TimeslotList
}

func NewTimeslotItemService(
	itemRepo repository.TimeslotItem,
	listRepo repository.TimeslotList,
) *TimeslotItemService {
	return &TimeslotItemService{listRepo: listRepo, itemRepo: itemRepo}
}

func (s *TimeslotItemService) Create(userID, listID int, item entity.TimeslotItem) (int, error) {
	if _, err := s.listRepo.GetByID(userID, listID); err != nil {
		return http.StatusBadRequest, err
	}

	return s.itemRepo.Create(listID, item)
}

func (s *TimeslotItemService) GetAll(userID, listID int) ([]entity.TimeslotItem, error) {
	return s.itemRepo.GetAll(userID, listID)
}

func (s *TimeslotItemService) GetByID(userID, itemID int) (entity.TimeslotItem, error) {
	return s.itemRepo.GetByID(userID, itemID)
}

func (s *TimeslotItemService) Delete(userID, itemID int) error {
	return s.itemRepo.Delete(userID, itemID)
}

func (s *TimeslotItemService) Update(userID, itemID int, input entity.UpdateItemInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.itemRepo.Update(userID, itemID, input)
}
