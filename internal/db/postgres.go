package db

import (
	"context"
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func NewMongo() *mongo.Database {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Printf("Unable to create client %v", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Printf("Unable to connect - %v", err)
		return nil
	}

	return client.Database("ledger")
}
