package main

import (
	"context"
	"errors"
	"ledger/internal/repository/mongo"
	"ledger/internal/repository/postgres"
	"ledger/internal/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"ledger/config"
	"ledger/internal/handler"
	"ledger/internal/queue"
)

func main() {
	// Load config/env vars
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// PostgreSQL setup
	pgDB, err := config.SetupPostgres(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}

	// MongoDB setup
	mongoClient, err := config.SetupMongo(cfg.MongoURI)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	// Initialize publisher
	transactionPublisher, err := queue.NewTransactionPublisher(cfg.RabbitMQURL, cfg.QueueName)
	if err != nil {
		log.Fatalf("failed to create transaction publisher: %v", err)
	}
	// Initialize account handler
	accountRepo := postgres.NewAccountRepository(pgDB)
	accountService := service.NewAccountService(accountRepo)
	accountHandler := handler.NewAccountHandler(accountService)
	// Initialize ledger Repository
	ledgerRepo := mongo.NewLedgerRepository(mongoClient, cfg.MongoDBName, cfg.MongoCollection)

	// Initialize transaction service
	transactionService := service.NewTransactionService(accountRepo, ledgerRepo, *transactionPublisher)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// Start consumer in background
	go func() {
		err := queue.StartTransactionConsumer(cfg.RabbitMQURL, cfg.QueueName, transactionService)
		if err != nil {
			log.Fatalf("failed to start transaction consumer: %v", err)
		}
	}()

	// Setup HTTP router
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	api.HandleFunc("/accounts/{id}", accountHandler.GetAccount).Methods("GET")
	api.HandleFunc("/accounts", accountHandler.GetAllAccounts).Methods("GET")
	api.HandleFunc("/accounts/{id}/balance", accountHandler.UpdateBalance).Methods("PUT")
	api.HandleFunc("/accounts/{id}", accountHandler.DeleteAccount).Methods("DELETE")
	api.HandleFunc("/transactions", transactionHandler.ProcessTransaction).Methods("POST")
	api.HandleFunc("/transactions/{id}", transactionHandler.GetTransactionHistory).Methods("GET")

	// Start HTTP server
	server := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: router,
	}
	log.Println("Starting server on", cfg.HTTPPort)
	if err := server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
		log.Fatalf("failed to start server: %v", err)
	}
	log.Println("Server stopped gracefully")
	log.Println("Exiting application")
	err = mongoClient.Disconnect(context.Background())
	if err != nil {
		log.Printf("failed to disconnect from MongoDB: %v", err)
		return

	}
	err = pgDB.Close()
	if err != nil {
		log.Printf("failed to close PostgreSQL connection: %v", err)
		return
	}

}
