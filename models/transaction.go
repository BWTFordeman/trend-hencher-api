package models

import (
	"time"

	"cloud.google.com/go/datastore"
)

type Transaction struct {
	TrendID     *datastore.Key `datastore:"trend_id"`
	DateBought  time.Time      `datastore:"date_bought"`
	DateSold    time.Time      `datastore:"date_sold"`
	PriceBought float64        `datastore:"price_bought"`
	PriceSold   float64        `datastore:"price_sold"`
	Volume      int64          `datastore:"volume"`
}

type TransactionResponse struct {
	ID          int64     `json:"id"`
	TrendID     int64     `json:"trend_id"`
	DateBought  time.Time `json:"date_bought"`
	DateSold    time.Time `json:"date_sold"`
	PriceBought float64   `json:"price_bought"`
	PriceSold   float64   `json:"price_sold"`
	Volume      int64     `json:"volume"`
}

type IndicatorType int64

const (
	IndicatorOver      IndicatorType = 1
	IndicatorUnder     IndicatorType = 2
	IndicatorCrossUp   IndicatorType = 3
	IndicatorCrossDown IndicatorType = 4
)
