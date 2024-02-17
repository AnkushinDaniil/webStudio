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

// func (s *TimeslotItemService) GetByID(userID, listID int) (entity.TimeslotsList, error) {
// 	return s.itemRepo.GetByID(userID, listID)
// }

// func (s *TimeslotItemService) Update(userID, listID int, input entity.UpdateListInput) error {
// 	if err := input.Validate(); err != nil {
// 		return err
// 	}

// 	return s.itemRepo.Update(userID, listID, input)
// }

// func (s *TimeslotItemService) Delete(userID, listID int) error {
// 	return s.itemRepo.Delete(userID, listID)
// }
