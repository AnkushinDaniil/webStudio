package entity

import "errors"

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
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Start       int    `json:"start"`
	End         int    `json:"end"`
	Done        bool   `json:"done"        db:"done"`
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
