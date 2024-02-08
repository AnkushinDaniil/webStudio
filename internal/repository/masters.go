package repository

import "main.go/internal/entity"

type MasterRepository interface {
	Save(master *entity.Master) (*entity.Master, error)
	FindAll() ([]entity.Master, error)
}
