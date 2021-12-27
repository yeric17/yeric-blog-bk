package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	DB_USER, DB_PASS, DB_NAME, DB_HOST, DB_PORT, DB_DRIVER, APP_PORT, ConnectionString string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	DB_USER = os.Getenv("DB_USER")
	DB_PASS = os.Getenv("DB_PASS")
	DB_NAME = os.Getenv("DB_NAME")
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = os.Getenv("DB_PORT")
	DB_DRIVER = os.Getenv("DB_DRIVER")
	APP_PORT = os.Getenv("APP_PORT")

	ConnectionString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)
}
