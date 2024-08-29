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
)

func main() {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create Datastore client: %v", err)
	}

	datastorerepo := repository.NewDatastoreRepository(ctx, client)
	trendService := services.NewTrendService(datastorerepo)
	trendHandler := handlers.NewTrendHandler(trendService)

	http.HandleFunc("/trend", trendHandler.GetTrend)
	http.HandleFunc("/trends", trendHandler.GetAllTrends)
	http.HandleFunc("/saveTrend", trendHandler.SaveTrend)
	http.HandleFunc("/transactions", trendHandler.GetTransactions)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
