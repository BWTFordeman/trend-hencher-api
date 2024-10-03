package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"trend-hencher-api/handlers"
	"trend-hencher-api/repository"
	"trend-hencher-api/services"

	"cloud.google.com/go/datastore"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables set in the system")
	}

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create Datastore client: %v", err)
	}

	log.Println("Starting server on :8080..")
	datastorerepo := repository.NewDatastoreRepository(ctx, client)
	trendService := services.NewTrendService(datastorerepo)
	trendHandler := handlers.NewTrendHandler(trendService)

	http.HandleFunc("/checkMarket", trendHandler.CheckMarket)
	http.HandleFunc("/trend", trendHandler.GetTrend)
	http.HandleFunc("/trends", trendHandler.GetAllTrends)
	http.HandleFunc("/saveTrend", trendHandler.SaveTrend)
	http.HandleFunc("/transactions", trendHandler.GetTransactions)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
