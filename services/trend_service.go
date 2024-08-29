package services

import (
	"trend-hencher-api/models"
	"trend-hencher-api/repository"

	"cloud.google.com/go/datastore"
)

type TrendService struct {
	repo *repository.DatastoreRepository
}

func NewTrendService(repo *repository.DatastoreRepository) *TrendService {
	return &TrendService{repo: repo}
}

func (s *TrendService) GetTrendByID(id int64) (*models.Trend, error) {
	return s.repo.GetTrendByID(id)
}

func (s *TrendService) GetAllTrends() ([]*models.Trend, error) {
	return s.repo.GetAllTrends()
}

func (s *TrendService) SaveTrend(trend *models.Trend) (*datastore.Key, error) {
	return s.repo.SaveTrend(trend)
}
