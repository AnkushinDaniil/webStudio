package repository

import "github.com/jmoiron/sqlx"

type Authorization interface{}

type TimeslotList interface{}

type TimeslotItem interface{}

type Repository struct {
	Authorization
	TimeslotList
	TimeslotItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
