package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DB_USER, DB_PASS, DB_NAME, DB_HOST, DB_PORT, DB_DRIVER, APP_PORT, ConnectionString, MODE string
	Mail                                                                                     = MailStruct{}
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		MODE = "dev"
	} else {
		fmt.Println("loaded env")
		MODE = os.Getenv("MODE")
	}

	if MODE == "dev" {
		fmt.Println("Running in development mode")
		DB_USER = "root"
		DB_PASS = "1234567"
		DB_NAME = "yeric-blog"
		DB_HOST = "localhost"
		DB_PORT = "5433"
		DB_DRIVER = "postgres"
		APP_PORT = "7070"
		Mail.Host = "smtp.gmail.com"
		Mail.Port = 587
		Mail.Username = "mail.manager.cf2345181@gmail.com"
		Mail.Password = "cF6XN$ozi0b3"
		Mail.From = "mail.manager.cf2345181@gmail.com"
	} else {
		fmt.Println("Running in production mode")
		DB_USER = os.Getenv("DB_USER")
		DB_PASS = os.Getenv("DB_PASS")
		DB_NAME = os.Getenv("DB_NAME")
		DB_HOST = os.Getenv("DB_HOST")
		DB_PORT = os.Getenv("DB_PORT")
		DB_DRIVER = os.Getenv("DB_DRIVER")
		APP_PORT = os.Getenv("APP_PORT")
		Mail.Host = os.Getenv("MAIL_HOST")
		Mail.Port, err = strconv.Atoi(os.Getenv("MAIL_PORT"))

		if err != nil {
			fmt.Println("MAIL_PORT is not a number")
			os.Exit(1)
		}

		Mail.Username = os.Getenv("MAIL_USER")
		Mail.Password = os.Getenv("MAIL_PASS")
		Mail.From = os.Getenv("MAIL_USER")
	}

	ConnectionString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)
}

type MailStruct struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}
