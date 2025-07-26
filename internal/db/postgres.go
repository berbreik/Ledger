package db

import (
	"database/sql"
	"log"
	"os"
	"time"
)

func NewPostgres() *sql.DB {
	dns := os.Getenv("POSTGRES_DNS")
	db, err := sql.Open("postgres", dns)
	if err != nil {
		log.Printf("Failed to open postgres connection: %v\n", err)
		return nil
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err := db.Ping(); err != nil {
		log.Printf("Failed to connect to postgres: %v", err)
	}
	return db
}
