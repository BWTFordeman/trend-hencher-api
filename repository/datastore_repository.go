package repository

import (
	"context"
	"log"
	"trend-hencher-api/models"

	"cloud.google.com/go/datastore"
)

type DatastoreRepository struct {
	client *datastore.Client
	ctx    context.Context
}

func NewDatastoreRepository(ctx context.Context, client *datastore.Client) *DatastoreRepository {
	return &DatastoreRepository{
		client: client,
		ctx:    ctx,
	}
}

// SaveTrend stores a Trend entity in Datastore
func (r *DatastoreRepository) SaveTrend(trend *models.Trend) (*datastore.Key, error) {
	key := datastore.IncompleteKey("Trend", nil)
	savedKey, err := r.client.Put(r.ctx, key, trend)
	if err != nil {
		log.Printf("Failed to save trend: %v", err)
		return nil, err
	}
	return savedKey, nil
}

// GetTrend retrieves a Trend entity by its ID
func (r *DatastoreRepository) GetTrendByID(id int64) (*models.Trend, error) {
	key := datastore.IDKey("Trend", id, nil)
	var trend models.Trend
	err := r.client.Get(r.ctx, key, &trend)
	if err != nil {
		log.Printf("Failed to get trend: %v", err)
		return nil, err
	}
	return &trend, nil
}

// GetAllTrends retrieves all Trend entities
func (r *DatastoreRepository) GetAllTrends() ([]models.TrendResponse, error) {
	var trends []models.Trend
	query := datastore.NewQuery("Trend")
	keys, err := r.client.GetAll(r.ctx, query, &trends)
	if err != nil {
		log.Printf("Failed to retrieve trends: %v", err)
		return nil, err
	}

	var response []models.TrendResponse
	for i, key := range keys {
		response = append(response, models.TrendResponse{
			ID:          key.ID,
			Stock:       trends[i].Stock,
			TrendScore:  trends[i].TrendScore,
			Date:        trends[i].Date,
			TrendValues: trends[i].TrendValues,
		})
	}
	return response, nil
}

// GetTransactions retrieves transactions related to a Trend
func (r *DatastoreRepository) GetTransactions(trendID int64) ([]models.TransactionResponse, error) {
	key := datastore.IDKey("Trend", trendID, nil)
	var transactions []models.Transaction
	query := datastore.NewQuery("Transaction").FilterField("trend_id", "=", key)
	keys, err := r.client.GetAll(r.ctx, query, &transactions)
	if err != nil {
		log.Printf("Failed to retrieve transactions: %v", err)
		return nil, err
	}

	var response []models.TransactionResponse
	for i, key := range keys {
		response = append(response, models.TransactionResponse{
			ID:          key.ID,
			TrendID:     trendID,
			DateBought:  transactions[i].DateBought,
			DateSold:    transactions[i].DateSold,
			PriceBought: transactions[i].PriceBought,
			PriceSold:   transactions[i].PriceSold,
			Volume:      transactions[i].Volume,
		})
	}
	return response, nil
}
