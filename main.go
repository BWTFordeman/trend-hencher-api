package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
)

type Trend struct {
	Message string `datastore:"message"`
}

type TrendResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/", helpHandler)
	http.HandleFunc("/trends", getTrendsHandler)
	http.HandleFunc("/save", saveHandler)
	http.ListenAndServe(":8080", nil)
}

func helpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Oops, nothing here. Want some help? Too bad! You won't get any!")
}

func getTrendsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create Datastore client: %v", err)
	}
	defer client.Close()

	// Create a query to retrieve all Trend entities
	query := datastore.NewQuery("Trend")

	var trends []Trend
	keys, err := client.GetAll(ctx, query, &trends)
	if err != nil {
		log.Printf("Failed to retrieve trends: %v", err)
		http.Error(w, "Failed to retrieve trends", http.StatusInternalServerError)
		return
	}

	// Prepare a slice of TrendResponse to return, including the numeric ID
	var response []TrendResponse
	for i, key := range keys {
		response = append(response, TrendResponse{
			ID:      key.ID,
			Message: trends[i].Message,
		})
	}

	// Convert the trends slice to JSON
	jsonResponse, err := json.Marshal(trends)
	if err != nil {
		log.Printf("Failed to marshal trends: %v", err)
		http.Error(w, "Failed to encode trends as JSON", http.StatusInternalServerError)
		return
	}

	// Send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var trend Trend
	err := json.NewDecoder(r.Body).Decode(&trend)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Set up Datastore client
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create Datastore client: %v", err)
	}
	defer client.Close()

	if err := saveTrend(ctx, client, &trend); err != nil {
		log.Printf("Failed to save trend: %v", err)
		http.Error(w, "Failed to save trend", http.StatusInternalServerError)
		return
	}

	// Response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Trend saved successfully")
}

func saveTrend(ctx context.Context, client *datastore.Client, trend *Trend) error {
	// Generate a new key for the entity
	key := datastore.IncompleteKey("Trend", nil)

	// Save the entity
	_, err := client.Put(ctx, key, trend)
	if err != nil {
		return err
	}

	return nil
}
