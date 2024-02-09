package entity

type User struct {
	Id           int    `json:"-"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	PassWordHash string `json:"passwordhash"`
}
