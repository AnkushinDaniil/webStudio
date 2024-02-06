package repository

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"main.go/entity"
)

type MasterRepository interface {
	Save(master *entity.Master) (*entity.Master, error)
	FindAll() ([]entity.Master, error)
}

type repo struct{}

func NewMasterRepository() MasterRepository {
	return &repo{}
}

const (
	projectId      string = "studio-jul"
	collectionName string = "masters"
)

func (*repo) Save(master *entity.Master) (*entity.Master, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("Failed to create a Firestore client: %v", err)
		return nil, err
	}

	defer client.Close()
	_, _, err = client.Collection(collectionName).Add(ctx, map[string]interface{}{
		"Id":   master.Id,
		"Name": master.Name,
	})
	if err != nil {
		log.Fatalf("Failed adding a new master: %v", err)
		return nil, err
	}
	return master, nil
}

func (*repo) FindAll() ([]entity.Master, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("Failed to create a Firestore client: %v", err)
		return nil, err
	}

	defer client.Close()
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
