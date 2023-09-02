package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func InitDatabase() *sql.DB {

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	connStr := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + name + " sslmode=" + sslmode

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Default().Panicf("error connecting to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Default().Panicf("error pinging to database: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Default().Panicf("error setting dialect: %v", err)
	}

	if err := goose.Up(db, "db/migrations"); err != nil {
		log.Default().Panicf("error on migration: %v", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)

	return db
}
