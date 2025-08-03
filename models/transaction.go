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
	IndicatorName       string        `bigquery:"indicator_name"`
	IndicatorType       IndicatorType `bigquery:"indicator_type"`
	IndicatorPeriod     int           `bigquery:"indicator_period"`
	IndicatorCheckValue Indicator     `bigquery:"indicator_check_value"`
}

type BuyScenario struct {
	Conditions []BuyCondition `bigquery:"conditions"`
}

type ConditionType int64

const (
	SellPercentage ConditionType = 1
	SellIndicator  ConditionType = 2
)

type SellCondition struct {
	ConditionType   ConditionType
	ProfitThreshold float64 `bigquery:"profit_threshold"`
	LossThreshold   float64 `bigquery:"loss_threshold"`
}

type SellScenario struct {
	Conditions []SellCondition
}

// General object that can have any indicator values:
type Indicator struct {
	IndicatorName        string
	IndicatorSMAPeriod   int
	IndicatorRSIPeriod   int
	IndicatorRSIStrength int
}
