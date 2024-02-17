package entity

type User struct {
	ID       int    `json:"-"        db:"id"`
	Name     string `json:"name"     db:"name"     binding:"required"`
	Username string `json:"username" db:"username" binding:"required"`
	Password string `json:"password" db:"password" binding:"required"`
}
