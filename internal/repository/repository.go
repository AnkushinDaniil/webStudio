package repository

type Authorization interface{}

type TimeslotList interface{}

type TimeslotItem interface{}

type Repository struct {
	Authorization
	TimeslotList
	TimeslotItem
}

func NewRepository() *Repository {
	return &Repository{}
}
