package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"main.go/internal/entity"
)

type TimeslotItemPostgres struct {
	db *sqlx.DB
}

func NewTimeslotItemPostgres(db *sqlx.DB) *TimeslotItemPostgres {
	return &TimeslotItemPostgres{db: db}
}

func (r *TimeslotItemPostgres) Create(listID int, item entity.TimeslotItem) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemID int
	createItemQuery := fmt.Sprintf(
		`INSERT INTO %s (title, description, beginning, finish) values ($1, $2, $3, $4) RETURNING id`,
		timeslotsItemsTable,
	)
	row := transaction.QueryRow(createItemQuery, item.Title, item.Description, item.Start, item.End)

	if err = row.Scan(&itemID); err != nil {
		if err1 := transaction.Rollback(); err != nil {
			return 0, err1
		}

		return 0, err
	}

	createListsItemsQuery := fmt.Sprintf(
		`INSERT INTO %s (list_id, item_id) VALUES ($1, $2)`,
		listsItemsTable,
	)

	if _, err = transaction.Exec(createListsItemsQuery, listID, itemID); err != nil {
		if err = transaction.Rollback(); err != nil {
			return 0, err
		}

		return 0, err
	}

	return itemID, transaction.Commit()
}

func (r *TimeslotItemPostgres) GetAll(userID, listID int) ([]entity.TimeslotItem, error) {
	var items []entity.TimeslotItem

	query := fmt.Sprintf(
		`SELECT ti.id, ti.title, ti.description,
		ti.beginning, ti.finish, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id
		INNER JOIN %s ul on ul.list_id = li.list_id WHERE li.list_id = $1 AND ul.user_id = $2`,
		timeslotsItemsTable,
		listsItemsTable,
		usersListsTable,
	)

	if err := r.db.Select(&items, query, listID, userID); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TimeslotItemPostgres) GetByID(userID, listID int) (entity.TimeslotsList, error) {
	var list entity.TimeslotsList

	query := fmt.Sprintf(
		`SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul 
                                       on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
		timeslotListsTable,
		usersListsTable,
	)
	err := r.db.Get(&list, query, userID, listID)

	return list, err
}

// func (r *TimeslotItemPostgres) Update(userID, listID int, input entity.UpdateListInput) error {
// 	setValues := make([]string, 0)
// 	args := make([]interface{}, 0)
// 	argID := 1

// 	if input.Title != nil {
// 		setValues = append(setValues, fmt.Sprintf("title=$%d", argID))
// 		args = append(args, *input.Title)
// 		argID++
// 	}

// 	if input.Description != nil {
// 		setValues = append(setValues, fmt.Sprintf("description=$%d", argID))
// 		args = append(args, *input.Description)
// 		argID++
// 	}

// 	setQuery := strings.Join(setValues, ", ")
// 	query := fmt.Sprintf(
// 		`UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d`,
// 		timeslotListsTable,
// 		setQuery,
// 		usersListsTable,
// 		argID,
// 		argID+1,
// 	)

// 	args = append(args, listID, userID)

// 	logrus.Debugf("updateQuery: %s", query)
// 	logrus.Debugf("args: %s", args...)

// 	_, err := r.db.Exec(query, args...)

// 	return err
// }

// func (r *TimeslotItemPostgres) Delete(userID, listID int) error {
// 	query := fmt.Sprintf(
// 		`DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2`,
// 		timeslotListsTable,
// 		usersListsTable,
// 	)
// 	_, err := r.db.Exec(query, userID, listID)

// 	return err
// }
