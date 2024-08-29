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
func (r *DatastoreRepository) GetAllTrends() ([]*models.Trend, error) {
	var trends []*models.Trend
	query := datastore.NewQuery("Trend")
	_, err := r.client.GetAll(r.ctx, query, &trends)
	if err != nil {
		log.Printf("Failed to retrieve trends: %v", err)
		return nil, err
	}
	return trends, nil
}
