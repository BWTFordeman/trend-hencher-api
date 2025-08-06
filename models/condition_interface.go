package models

func (b BuyCondition) GetIndicatorName() string {
	return b.IndicatorName
}

func (b BuyCondition) GetIndicatorType() IndicatorType {
	return b.IndicatorType
}

func (b BuyCondition) GetIndicatorPeriod() int {
	return b.IndicatorPeriod
}

func (b BuyCondition) GetCheckValue() Indicator {
	return b.IndicatorCheckValue
}

func (s SellCondition) GetIndicatorName() string {
	return s.IndicatorName
}

func (s SellCondition) GetIndicatorType() IndicatorType {
	return s.IndicatorType
}

func (s SellCondition) GetIndicatorPeriod() int {
	return s.IndicatorPeriod
}

func (s SellCondition) GetCheckValue() Indicator {
	return s.IndicatorCheckValue
}
