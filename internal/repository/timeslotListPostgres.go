package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"main.go/internal/entity"
)

type TimeslotListPostgres struct {
	db *sqlx.DB
}

func NewTimeslotListPostgres(db *sqlx.DB) *TimeslotListPostgres {
	return &TimeslotListPostgres{db: db}
}

func (r *TimeslotListPostgres) Create(userId int, list entity.TimeslotsList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf(
		"INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id",
		timeslotListsTable,
	)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err = row.Scan(&id); err != nil {
		if err = tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf(
		"INSERT INTO %s (user_id, list_id) VALUES ($1, $2)",
		usersListsTable,
	)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TimeslotListPostgres) GetAll(userId int) ([]entity.TimeslotsList, error) {
	var lists []entity.TimeslotsList
	query := fmt.Sprintf(
		`SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul 
                                       on tl.id = ul.list_id WHERE ul.user_id = $1`,
		timeslotListsTable,
		usersListsTable,
	)
	err := r.db.Select(&lists, query, userId)
	return lists, err
}

func (r *TimeslotListPostgres) GetById(userId, listId int) (entity.TimeslotsList, error) {
	var list entity.TimeslotsList
	query := fmt.Sprintf(
		`SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul 
                                       on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
		timeslotListsTable,
		usersListsTable,
	)
	err := r.db.Get(&list, query, userId, listId)
	return list, err
}
