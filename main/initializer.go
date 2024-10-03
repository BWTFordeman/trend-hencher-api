package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/datastore"
	"github.com/joho/godotenv"

	"trend-hencher-api/handlers"
	"trend-hencher-api/repository"
	"trend-hencher-api/services"
)

func initEnv() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables set in the system")
	}
}

// initHandler initializes the necessary services and returns a TrendHandler
func initTrendHandler() *handlers.TrendHandler {
	// Create the context
	ctx := context.Background()

	// Initialize Datastore client
	datastoreClient, err := datastore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create Datastore client: %v", err)
	}

	// Initialize BigQuery client
	bigQueryClient, err := bigquery.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create BigQuery client: %v", err)
	}

	// Initialize repositories
	bigqueryRepo := repository.NewBigQueryRepository(ctx, bigQueryClient)
	datastorerepo := repository.NewDatastoreRepository(ctx, datastoreClient)

	// Initialize services
	bigQueryService := services.NewBigQueryService(bigqueryRepo)
	trendService := services.NewTrendService(datastorerepo)

	// Return the initialized TrendHandler
	return handlers.NewTrendHandler(trendService, bigQueryService)
}
