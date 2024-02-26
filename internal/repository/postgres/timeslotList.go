package postgres

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"main.go/internal/entity"
)

type TimeslotList interface {
	Create(userID int, list entity.TimeslotsList) (int, error)
	GetAll(userID int) ([]entity.TimeslotsList, error)
	GetByID(userID, listID int) (entity.TimeslotsList, error)
	Delete(userID, listID int) error
	Update(userID, listID int, input entity.UpdateListInput) error
}

type TimeslotListPostgres struct {
	db *sqlx.DB
}

func NewTimeslotListPostgres(db *sqlx.DB) *TimeslotListPostgres {
	return &TimeslotListPostgres{db: db}
}

func (r *TimeslotListPostgres) Create(userID int, list entity.TimeslotsList) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var listID int

	createListQuery := fmt.Sprintf(
		`
			INSERT INTO %s (title, description)
			    VALUES ($1, $2)
			RETURNING
			    id`,
		TimeslotListsTable,
	)
	row := transaction.QueryRow(createListQuery, list.Title, list.Description)

	if err = row.Scan(&listID); err != nil {
		if err1 := transaction.Rollback(); err1 != nil {
			return 0, err
		}

		return 0, err
	}

	createUsersListQuery := fmt.Sprintf(
		`
			INSERT INTO %s (user_id, list_id)
			    VALUES ($1, $2)`,
		UsersListsTable,
	)

	if _, err = transaction.Exec(createUsersListQuery, userID, listID); err != nil {
		if err1 := transaction.Rollback(); err1 != nil {
			return 0, err1
		}

		return 0, err
	}

	return listID, transaction.Commit()
}

func (r *TimeslotListPostgres) GetAll(userID int) ([]entity.TimeslotsList, error) {
	var lists []entity.TimeslotsList

	query := fmt.Sprintf(
		`
			SELECT
			    tl.id,
			    tl.title,
			    tl.description
			FROM
			    %s tl
			    INNER JOIN %s ul ON tl.id = ul.list_id
			WHERE
			    ul.user_id = $1`,
		TimeslotListsTable,
		UsersListsTable,
	)
	err := r.db.Select(&lists, query, userID)

	return lists, err
}

func (r *TimeslotListPostgres) GetByID(userID, listID int) (entity.TimeslotsList, error) {
	var list entity.TimeslotsList

	query := fmt.Sprintf(
		`
			SELECT
			    tl.id,
			    tl.title,
			    tl.description
			FROM
			    %s tl
			    INNER JOIN %s ul ON tl.id = ul.list_id
			WHERE
			    ul.user_id = $1
			    AND ul.list_id = $2`,
		TimeslotListsTable,
		UsersListsTable,
	)
	err := r.db.Get(&list, query, userID, listID)

	return list, err
}

func (r *TimeslotListPostgres) Update(userID, listID int, input entity.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argID))
		args = append(args, *input.Title)
		argID++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argID))
		args = append(args, *input.Description)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(
		`
			UPDATE
			    %s tl
			SET
			    %s
			FROM
			    %s ul
			WHERE
			    tl.id = ul.list_id
			    AND ul.list_id = $ %d
			    AND ul.user_id = $ %d`,
		TimeslotListsTable,
		setQuery,
		UsersListsTable,
		argID,
		argID+1,
	)

	args = append(args, listID, userID)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args...)

	_, err := r.db.Exec(query, args...)

	return err
}

func (r *TimeslotListPostgres) Delete(userID, listID int) error {
	query := fmt.Sprintf(
		`
			DELETE FROM %s tl USING %s ul
			WHERE tl.id = ul.list_id
			    AND ul.user_id = $1
			    AND ul.list_id = $2`,
		TimeslotListsTable,
		UsersListsTable,
	)
	_, err := r.db.Exec(query, userID, listID)

	return err
}
