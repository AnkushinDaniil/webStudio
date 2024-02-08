package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"

	"google.golang.org/api/iterator"
	"main.go/entity"
	"main.go/repository"
)

type MasterService interface {
	Validate(master *entity.Master) error
	Create(master *entity.Master) (*entity.Master, error)
	FindAll() ([]entity.Master, error)
}

type service struct{}

var repo repository.MasterRepository = repository.NewFirestoreRepository()

func NewMasterService() MasterService {
	return &service
}

func Validate(master *entity.Master) error {
	if master == nil {
		err := errors.New("The master in nil")
		return err
	}
	if master.Name == "" {
		err := errors.New("The master's name in empty")
		return err
	}
}

func Create(master *entity.Master) (*entity.Master, error) {
	input, err := json.Marshal(master)
	if err != nil {
		log.Fatalf("Master marshalling was failed: %v", err)
		return master, err
	}
	hash := sha256.Sum256(input)
	master.Id = hex.EncodeToString(hash[:])
	return repo.Save(master)
}

func FindAll() ([]entity.Master, error) {
	var masters []entity.Master
	itr := client.Collection(collectionName).Documents(ctx)
	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate the lists of masters: %v", err)
			return nil, err
		}
		masters = append(masters, entity.Master{
			Id:   doc.Data()["Id"].(string),
			Name: doc.Data()["Name"].(string),
		})
	}

	return masters, nil
}
