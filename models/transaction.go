package models

import (
	"cloud.google.com/go/datastore"
)

type Transaction struct {
	TrendID     *datastore.Key `datastore:"trend_id"`
	DateBought  string         `datastore:"date_bought"`
	DateSold    string         `datastore:"date_sold"`
	PriceBought float64        `datastore:"price_bought"`
	PriceSold   float64        `datastore:"price_sold"`
	Volume      int64          `datastore:"volume"`
}

type TransactionResponse struct {
	ID          int64   `json:"id"`
	TrendID     int64   `json:"trend_id"`
	DateBought  string  `json:"date_bought"`
	DateSold    string  `json:"date_sold"`
	PriceBought float64 `json:"price_bought"`
	PriceSold   float64 `json:"price_sold"`
	Volume      int64   `json:"volume"`
}

type IndicatorType int64

const (
	IndicatorOver      IndicatorType = 1
	IndicatorUnder     IndicatorType = 2
	IndicatorCrossUp   IndicatorType = 3
	IndicatorCrossDown IndicatorType = 4
)

type BuyCondition struct {
	IndicatorName       string                 `json:"indicator_name"`
	IndicatorType       IndicatorType          `json:"indicator_type"`
	IndicatorConfig     map[string]interface{} `json:"indicator_config"`
	IndicatorPeriod     int                    `json:"indicator_period"`
	IndicatorCheckValue string                 `json:"indicator_check_value"`
}

type BuyScenario struct {
	Conditions []BuyCondition `json:"conditions"`
}

type SellScenario struct {
	ProfitThreshold float64 `json:"profit_threshold"`
	LossThreshold   float64 `json:"loss_threshold"`
}
