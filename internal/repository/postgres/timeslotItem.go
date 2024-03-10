package postgres

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"main.go/internal/entity"
)

type TimeslotItem interface {
	Create(listID int, item entity.TimeslotItem) (int, error)
	GetAll(userID, listID int) ([]entity.TimeslotItem, error)
	GetByID(userID, itemID int) (entity.TimeslotItem, error)
	Delete(userID, itemID int) error
	Update(userID, itemID int, input entity.UpdateItemInput) error
}

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
		`
			INSERT INTO %s (title, description, beginning, finish)
			    VALUES ($1, $2, $3, $4)
			RETURNING
			    id`,
		TimeslotsItemsTable,
	)
	row := transaction.QueryRow(createItemQuery, item.Title, item.Description, item.Start, item.End)

	if err = row.Scan(&itemID); err != nil {
		if err1 := transaction.Rollback(); err1 != nil {
			return 0, err1
		}

		return 0, err
	}

	createListsItemsQuery := fmt.Sprintf(
		`
			INSERT INTO %s (list_id, item_id)
			    VALUES ($1, $2)`,
		ListsItemsTable,
	)

	if _, err = transaction.Exec(createListsItemsQuery, listID, itemID); err != nil {
		if err1 := transaction.Rollback(); err1 != nil {
			return 0, err
		}

		return 0, err
	}

	return itemID, transaction.Commit()
}

func (r *TimeslotItemPostgres) GetAll(userID, listID int) ([]entity.TimeslotItem, error) {
	var items []entity.TimeslotItem

	query := fmt.Sprintf(
		`
			SELECT
			    ti.id,
			    ti.title,
			    ti.description,
			    ti.beginning,
			    ti.finish,
			    ti.done
			FROM
			    %s ti
			    INNER JOIN %s li ON li.item_id = ti.id
			    INNER JOIN %s ul ON ul.list_id = li.list_id
			WHERE
			    li.list_id = $1
			    AND ul.user_id = $2`,
		TimeslotsItemsTable,
		ListsItemsTable,
		UsersListsTable,
	)

	if err := r.db.Select(&items, query, listID, userID); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TimeslotItemPostgres) GetByID(userID, itemID int) (entity.TimeslotItem, error) {
	var item entity.TimeslotItem

	query := fmt.Sprintf(
		`
			SELECT
			    ti.id,
			    ti.title,
			    ti.description,
			    ti.beginning,
			    ti.finish,
			    ti.done
			FROM
			    %s ti
			    INNER JOIN %s li ON li.item_id = ti.id
			    INNER JOIN %s ul ON ul.list_id = li.list_id
			WHERE
			    ti.id = $1
			    AND ul.user_id = $2`,
		TimeslotsItemsTable,
		ListsItemsTable,
		UsersListsTable,
	)
	if err := r.db.Get(&item, query, itemID, userID); err != nil {
		return item, err
	}

	return item, nil
}

func (r *TimeslotItemPostgres) Update(userID, itemID int, input entity.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	refVal := reflect.ValueOf(&input).Elem()
	refType := reflect.TypeOf(input)

	for i := 0; i < refVal.NumField(); i++ {
		field := refVal.Field(i)
		if !field.IsNil() {
			setValues = append(
				setValues,
				fmt.Sprintf("%s=$%d", refType.Field(i).Tag.Get("db"), argID),
			)
			args = append(args, field.Elem().Interface())
			argID++
		}
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(
		`
			UPDATE
			    %s ti
			SET
			    %s
			FROM
			    %s li,
			    %s ul
			WHERE
			    ti.id = li.item_id
			    AND li.list_id = ul.list_id
			    AND ul.user_id = $ %d
			    AND ti.id = $ %d`,
		TimeslotsItemsTable,
		setQuery,
		ListsItemsTable,
		UsersListsTable,
		argID,
		argID+1,
	)

	args = append(args, userID, itemID)

	logrus.Debugf("updateQuery: %s \n", query)
	logrus.Debugf("args: %s \n", args...)

	_, err := r.db.Exec(query, args...)

	return err
}

func (r *TimeslotItemPostgres) Delete(userID, itemID int) error {
	query := fmt.Sprintf(
		`
			DELETE FROM %s ti USING %s li, %s ul
			WHERE ti.id = li.item_id
			    AND ul.user_id = $1
			    AND li.list_id = ul.list_id
			    AND ul.user_id = $1
			    AND ti.id = $2`,
		TimeslotsItemsTable,
		ListsItemsTable,
		UsersListsTable,
	)
	_, err := r.db.Exec(query, userID, itemID)

	return err
}

func (r *TimeslotItemPostgres) GetByRange(input entity.ItemsByRange) ([]entity.TimeslotItemWithUsername, error) {
	var items []entity.TimeslotItemWithUsername

	query := fmt.Sprintf(
		`
			SELECT
			    ti.id,
			    ti.title,
			    ti.description,
			    ti.beginning,
			    ti.finish,
			    ti.done,
			    u.username
			FROM
			    %s ti
			    INNER JOIN %s li ON li.item_id = ti.id
			    INNER JOIN %s ul ON ul.list_id = li.list_id
				INNER JOIN %s u ON u.id = ul.user_id
			WHERE
			    ti.beginning >= $1
			    AND ti.finish <= $2`,
		TimeslotsItemsTable,
		ListsItemsTable,
		UsersListsTable,
		UsersTable,
	)

	if err := r.db.Select(&items,
		query,
		input.Start,
		input.End,
	); err != nil {
		return nil, err
	}

	return items, nil
}
