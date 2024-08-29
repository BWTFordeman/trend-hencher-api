package models

import "time"

type IndicatorType int64

type TrendValues struct {
	IndicatorName  string        `json:"indicator_name"`
	IndicatorType  IndicatorType `json:"indicator_type"`
	IndicatorValue float64       `json:"indiator_value"`
}

type Trend struct {
	Stock       string      `datastore:"stock"`
	TrendScore  float64     `datastore:"trend_score"`
	Date        time.Time   `datastore:"date"`
	TrendValues TrendValues `datastore:"trend_values"`
}

type TrendResponse struct {
	ID          int64       `json:"id"`
	Stock       string      `json:"stock"`
	TrendScore  float64     `json:"trend_score"`
	Date        time.Time   `json:"date"`
	TrendValues TrendValues `json:"trend_values"`
}
