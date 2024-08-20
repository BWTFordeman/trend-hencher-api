package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
)

type Trend struct {
	Stock       string      `datastore:"stock"`
	TrendScore  float64     `datastore:"trend_score"`
	Date        time.Time   `datastore:"date"`
	TrendValues TrendValues `datastore:"trend_values"`
}

type TrendResponse struct {
	ID int64 `json:"id"`
	Trend
}

type Transaction struct {
	TrendID     *datastore.Key `datastore:"trend_id"`
	DateBought  time.Time      `datastore:"date_bought"`
	DateSold    time.Time      `datastore:"date_sold"`
	PriceBought float64        `datastore:"price_bought"`
	PriceSold   float64        `datastore:"price_sold"`
	Volume      int64          `datastore:"volume"`
}

type TrendValues struct {
	IndicatorName  string        `json:"indicator_name"`
	IndicatorType  IndicatorType `json:"indicator_type"`
	IndicatorValue float64       `json:"indiator_value"`
}

type IndicatorType int64

const (
	IndicatorOver      IndicatorType = 1
	IndicatorUnder     IndicatorType = 2
	IndicatorCrossUp   IndicatorType = 3
	IndicatorCrossDown IndicatorType = 4
)

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
			ID:    key.ID,
			Trend: trends[i],
		})
	}

	// Convert the trends slice to JSON
	jsonResponse, err := json.Marshal(response)
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
	var requestData struct {
		Stock string `json:"stock"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Fetch data here...
	//trendScore, trendValues, transactions, err := fetchTrendData()
	trendValues := TrendValues{}
	trendValues.IndicatorName = "RSI"
	trendValues.IndicatorType = IndicatorOver
	trendValues.IndicatorValue = 70

	trend := Trend{
		Stock:       requestData.Stock,
		Date:        time.Now(),
		TrendScore:  56.4,
		TrendValues: trendValues,
	}

	// Set up Datastore client
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create Datastore client: %v", err)
	}
	defer client.Close()

	// Save trend results:
	trendKey := datastore.IncompleteKey("Trend", nil)
	completeKey, err := client.Put(ctx, trendKey, &trend)
	if err != nil {
		log.Printf("Failed to save trend: %v", err)
		http.Error(w, "Failed to save trend", http.StatusInternalServerError)
		return
	}

	// Save transactions for Trend:
	transactions := []Transaction{
		{
			DateBought:  time.Now(),
			DateSold:    time.Now().Add(1 * time.Hour),
			PriceBought: 150,
			PriceSold:   156,
			Volume:      100,
		},
		{
			DateBought:  time.Now().Add(1*time.Hour + 30*time.Minute),
			DateSold:    time.Now().Add(2 * time.Hour),
			PriceBought: 154,
			PriceSold:   159,
			Volume:      100,
		},
		{
			DateBought:  time.Now().Add(2*time.Hour + 30*time.Minute),
			DateSold:    time.Now().Add(3 * time.Hour),
			PriceBought: 160,
			PriceSold:   163,
			Volume:      100,
		},
	}
	for _, transaction := range transactions {
		transaction.TrendID = completeKey
		transactionKey := datastore.IncompleteKey("Transaction", nil)
		_, err := client.Put(ctx, transactionKey, &transaction)
		if err != nil {
			log.Printf("Failed to save transaction: %v", err)
			http.Error(w, "Failed to save transaction", http.StatusInternalServerError)
			return
		}
	}

	// Response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Trends saved successfully")
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

func saveTransaction(ctx context.Context, client *datastore.Client, transaction *Transaction) error {
	key := datastore.IncompleteKey("Transaction", nil)
	_, err := client.Put(ctx, key, transaction)
	return err
}
