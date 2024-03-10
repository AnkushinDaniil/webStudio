package service

import (
	"net/http"

	"main.go/internal/entity"
)

type TimeslotItemRepository interface {
	Create(listID int, item entity.TimeslotItem) (int, error)
	GetAll(userID, listID int) ([]entity.TimeslotItem, error)
	GetByID(userID, itemID int) (entity.TimeslotItem, error)
	Delete(userID, itemID int) error
	Update(userID, itemID int, input entity.UpdateItemInput) error
	GetByRange(input entity.ItemsByRange) ([]entity.TimeslotItemWithUsername, error)
}

type TimeslotItemService struct {
	itemRepo TimeslotItemRepository
	listRepo TimeslotListRepository
}

func NewTimeslotItemService(
	itemRepo TimeslotItemRepository,
	listRepo TimeslotListRepository,
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

func (s *TimeslotItemService) GetByRange(input entity.ItemsByRange) ([]entity.TimeslotItemWithUsername, error) {
	return s.itemRepo.GetByRange(input)
}
