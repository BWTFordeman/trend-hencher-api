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
	Key     *datastore.Key `datastore:"__key__"` // Datastore field fo key
	Message string         `datastore:"message"`
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/save", saveHandler)
	http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, trend setters!")
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

	log.Printf("saved values %v, %v", ctx, key)

	// Save the entity
	completeKey, err := client.Put(ctx, key, trend)
	if err != nil {
		return err
	}
	log.Printf("complete key %v", completeKey)

	// Update the Trend struct with the generated complete key
	trend.Key = completeKey
	return nil
}
