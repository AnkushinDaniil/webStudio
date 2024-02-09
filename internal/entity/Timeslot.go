package entity

type Timeslot struct {
	Id         int    `json:"-"`
	Owner      string `json:"name"`
	Timestamps [2]int `json:"timestamps"`
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
}

type ListsItem struct {
	Id     int
	ListId int
	ItemId int
}
