package db

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var localUrl = ""

func NewDatabaseConnection() *sqlx.DB {

	dsn := os.Getenv("GIN_MODE")

	if dsn != "release" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("failed to load env", err)
		}
	}

	dbUrl := os.Getenv("DSN")

	if dbUrl == "" {
		dbUrl = localUrl
	}

	Db, err := sqlx.Connect("mysql", dbUrl)

	if err != nil {
		panic(err)
	}

	return Db

}
