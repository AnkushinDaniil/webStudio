package entity

type TimeslotList struct {
	Id          int    `json:"id"          db:"id"`
	Title       string `json:"title"       db:"title"       binding:"required"`
	Description string `json:"description" db:"description"`
}

type UsersList struct {
	Id     int
	UserId int
	ListId int
}

type TimeslotItem struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Start       int    `json:"start"`
	End         int    `json:"end"`
	Done        bool   `json:"done"        db:"done"`
}

type ListsItem struct {
	Id     int
	ListId int
	ItemId int
}
