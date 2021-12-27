package models

import (
	"database/sql"
	"fmt"

	"yeric-blog/config"

	_ "github.com/lib/pq"
)

var (
	Connection *sql.DB
)

func init() {
	Connection = GetConnection()
}

func GetConnection() *sql.DB {

	db, err := sql.Open(config.DB_DRIVER, config.ConnectionString)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected postgresSQL")
	return db
}
