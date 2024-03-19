package entity

import (
	"errors"
	"reflect"
	"time"
)

type TimeslotsList struct {
	ID          int    `db:"id"          json:"id"`
	Title       string `db:"title"       json:"title"       binding:"required"`
	Description string `db:"description" json:"description"`
}

type UsersList struct {
	ID     int
	UserID int
	ListID int
}

type TimeslotItem struct {
	ID          int       `json:"id"          db:"id"`
	Title       string    `json:"title"       db:"title"       binding:"required"`
	Description string    `json:"description" db:"description"`
	Start       time.Time `json:"start"       db:"beginning"   binding:"required"`
	End         time.Time `json:"end"         db:"finish"      binding:"required"`
	Done        bool      `json:"done"        db:"done"`
	Username    string    `json:"username"    db:"username"`
	Color       string    `json:"color"       db:"color"`
}

type ItemsByRange struct {
	Start time.Time `form:"start"`
	End   time.Time `form:"end"`
}

type ListsItem struct {
	ID     int
	ListID int
	ItemID int
}

type UpdateListInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (i *UpdateListInput) Validate() error {
	val := reflect.ValueOf(i).Elem()
	for j := 0; j < val.NumField(); j++ {
		if !val.Field(j).IsNil() {
			return nil
		}
	}

	return errors.New("update structure has no values")
}

type UpdateItemInput struct {
	Title       *string    `json:"title"       db:"title"`
	Description *string    `json:"description" db:"description"`
	Start       *time.Time `json:"start"       db:"beginning"`
	End         *time.Time `json:"end"         db:"finish"`
	Done        *bool      `json:"done"        db:"done"`
}

func (i *UpdateItemInput) Validate() error {
	val := reflect.ValueOf(i).Elem()
	for j := 0; j < val.NumField(); j++ {
		if !val.Field(j).IsNil() {
			return nil
		}
	}

	return errors.New("update structure has no values")
}
