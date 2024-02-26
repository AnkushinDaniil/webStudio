package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	UsersTable          = "users"
	TimeslotListsTable  = "timeslots_lists"
	UsersListsTable     = "users_lists"
	TimeslotsItemsTable = "timeslots_items"
	ListsItemsTable     = "lists_items"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	dataBase, err := sqlx.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode),
	)
	if err != nil {
		return nil, err
	}

	if err = dataBase.Ping(); err != nil {
		return nil, err
	}

	return dataBase, nil
}
