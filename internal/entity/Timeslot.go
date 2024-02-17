package entity

import (
	"errors"
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
	if i.Title == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
