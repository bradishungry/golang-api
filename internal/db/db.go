package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

//There is probably a better way to do this but it works
var db *sql.DB

func InitDB(dataSourceName string) error {
	var err error

	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}

	return db.Ping()
}

func GetDBHandle() *sql.DB {
	return db
}
