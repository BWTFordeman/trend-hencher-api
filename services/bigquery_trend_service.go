package services

import (
	"trend-hencher-api/models"
	"trend-hencher-api/repository"
)

type BigQueryTrendService struct {
	repo *repository.BigQueryRepository
}

func NewBigQueryService(repo *repository.BigQueryRepository) *BigQueryTrendService {
	return &BigQueryTrendService{repo: repo}
}

func (s *BigQueryTrendService) SaveTrend(trend *models.Trend) error {
	return s.repo.SaveTrend(trend)
}

func (s *BigQueryTrendService) SaveTransactions(transactions []models.Transaction) error {
	return s.repo.SaveTransactions(transactions)
}
