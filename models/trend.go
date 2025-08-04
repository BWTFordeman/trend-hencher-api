package models

import (
	"time"
)

type TrendValues struct {
	IndicatorName   string        `json:"indicator_name"`
	IndicatorType   IndicatorType `json:"indicator_type"`
	IndicatorPeriod int           `json:"indicator_period"`
}

type Trend struct {
	TrendID               string       `bigquery:"trend_id"`
	Stock                 string       `bigquery:"stock"`
	TrendScore            float64      `bigquery:"trend_score"`
	Date                  time.Time    `bigquery:"date"`
	IndicatorBuyScenario  BuyScenario  `bigquery:"indicator_buy_scenario"`
	IndicatorSellScenario SellScenario `bigquery:"indicator_sell_scenario"`
}

type TrendResponse struct {
	ID                    int64        `json:"id"`
	Stock                 string       `json:"stock"`
	TrendScore            float64      `json:"trend_score"`
	Date                  time.Time    `json:"date"`
	TrendValues           TrendValues  `json:"trend_values"`
	IndicatorBuyScenario  BuyScenario  `json:"indicator_buy_scenario"`
	IndicatorSellScenario SellScenario `json:"indicator_sell_scenario"`
}
