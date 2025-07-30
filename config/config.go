package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	PostgresDSN     string
	MongoURI        string
	MongoDBName     string
	MongoCollection string
	RabbitMQURL     string
	QueueName       string
	HTTPPort        string
}

// Load reads environment variables into a config struct
func Load() (*Config, error) {

	cfg := &Config{
		PostgresDSN:     os.Getenv("POSTGRES_DSN"),
		MongoURI:        os.Getenv("MONGO_URI"),
		MongoDBName:     os.Getenv("MONGO_DB_NAME"),
		MongoCollection: os.Getenv("MONGO_COLLECTION"),
		RabbitMQURL:     os.Getenv("RABBITMQ_URL"),
		QueueName:       os.Getenv("QUEUE_NAME"),
		HTTPPort:        os.Getenv("HTTP_PORT"),
	}

	if cfg.PostgresDSN == "" || cfg.MongoURI == "" || cfg.MongoDBName == "" || cfg.RabbitMQURL == "" || cfg.HTTPPort == "" {
		return nil, fmt.Errorf("missing one or more required environment variables")
	}

	return cfg, nil
}

// SetupPostgres connects to PostgreSQL using the standard library
func SetupPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	log.Println("Connected to PostgreSQL")
	return db, nil
}

// SetupMongo connects to MongoDB  and returns the client
func SetupMongo(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("mongo connect error: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping failed: %w", err)
	}

	log.Println("Connected to MongoDB")
	return client, nil
}
