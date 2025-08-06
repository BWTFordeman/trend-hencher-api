package models

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
	ConditionType       ConditionType
	ProfitThreshold     float64       `bigquery:"profit_threshold"`
	LossThreshold       float64       `bigquery:"loss_threshold"`
	IndicatorName       string        `bigquery:"indicator_name"`
	IndicatorType       IndicatorType `bigquery:"indicator_type"`
	IndicatorPeriod     int           `bigquery:"indicator_period"`
	IndicatorCheckValue Indicator     `bigquery:"indicator_check_value"`
}

type SellScenario struct {
	Conditions []SellCondition
}
