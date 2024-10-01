package models

type Transaction struct {
	TransactionID string  `bigquery:"transaction_id"`
	TrendID       string  `bigquery:"trend_id"`
	DateBought    string  `bigquery:"date_bought"`
	DateSold      string  `bigquery:"date_sold"`
	PriceBought   float64 `bigquery:"price_bought"`
	PriceSold     float64 `bigquery:"price_sold"`
	Volume        int64   `bigquery:"volume"`
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
