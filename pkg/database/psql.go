package database

import (
	"database/sql"
	"fmt"
)

const (
	host     = "db"
	user     = "postgres"
	password = "root"
	dbname   = "app_db"
)

func OpenDatabaseConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	return db, err
}
