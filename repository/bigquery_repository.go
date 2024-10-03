package repository

import (
	"context"
	"log"
	"trend-hencher-api/models"

	"cloud.google.com/go/bigquery"
)

type BigQueryRepository struct {
	client *bigquery.Client
	ctx    context.Context
}

func NewBigQueryRepository(ctx context.Context, client *bigquery.Client) *BigQueryRepository {
	return &BigQueryRepository{
		client: client,
		ctx:    ctx,
	}
}

func (r *BigQueryRepository) SaveTrend(trend *models.Trend) error {
	table := r.client.Dataset("trend_dataset").Table("Trend")
	inserter := table.Inserter()

	// Insert the Trend into BigQuery
	err := inserter.Put(r.ctx, trend)
	if err != nil {
		log.Printf("Failed to save trend in BigQuery: %v", err)
		return err
	}

	return nil
}

func (r *BigQueryRepository) SaveTransactions(transactions []models.Transaction) error {
	table := r.client.Dataset("trend_dataset").Table("Transaction")
	inserter := table.Inserter()

	// Insert the transactions into BigQuery
	err := inserter.Put(r.ctx, transactions)
	if err != nil {
		log.Printf("Failed to save transactions in BigQuery: %v", err)
		return err
	}

	return nil
}
