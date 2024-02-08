package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"

	"main.go/internal/entity"
	"main.go/internal/repository"
)

type MasterService interface {
	Validate(master *entity.Master) error
	Create(master *entity.Master) (*entity.Master, error)
	FindAll() ([]entity.Master, error)
}

type service struct{}

var repo repository.MasterRepository

func NewMasterService(repository repository.MasterRepository) MasterService {
	repo = repository
	return &service{}
}

// Create implements MasterService.
func (*service) Create(master *entity.Master) (*entity.Master, error) {
	input, err := json.Marshal(master)
	if err != nil {
		log.Fatalf("Master marshalling was failed: %v", err)
		return master, err
	}
	hash := sha256.Sum256(input)
	master.Id = hex.EncodeToString(hash[:])
	return repo.Save(master)
}

// FindAll implements MasterService.
func (*service) FindAll() ([]entity.Master, error) {
	return repo.FindAll()
}

// Validate implements MasterService.
func (*service) Validate(master *entity.Master) error {
	if master == nil {
		err := errors.New("the master in nil")
		return err
	}
	if master.Name == "" {
		err := errors.New("the master's name in empty")
		return err
	}
	return nil
}
