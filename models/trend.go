package models

import (
	"time"
)

type TrendValues struct {
	IndicatorName      string        `json:"indicator_name"`
	IndicatorType      IndicatorType `json:"indicator_type"`
	IndicatorSMAPeriod int           `json:"indicator_sma"`
}

type Trend struct {
	Stock                 string       `datastore:"stock"`
	TrendScore            float64      `datastore:"trend_score"`
	Date                  time.Time    `datastore:"date"`
	IndicatorBuyScenario  BuyScenario  `json:"indicator_buy_scenario"`
	IndicatorSellScenario SellScenario `json:"indicator_sell_scenario"`
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
