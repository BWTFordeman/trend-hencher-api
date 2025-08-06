package models

import "github.com/markcheno/go-talib"

type Indicator struct {
	IndicatorName     string
	IndicatorPeriod   int
	IndicatorStrength float64
}

type IndicatorType int64

const (
	IndicatorOver      IndicatorType = 1
	IndicatorUnder     IndicatorType = 2
	IndicatorCrossUp   IndicatorType = 3
	IndicatorCrossDown IndicatorType = 4
)

type IndicatorCondition interface {
	GetIndicatorName() string
	GetIndicatorType() IndicatorType
	GetIndicatorPeriod() int
	GetCheckValue() Indicator
}

type IndicatorKey struct {
	Name   string
	Period int
}

func GetPredefinedIndicators(buyScenario BuyScenario, sellScenario SellScenario, data []IntradayData) map[IndicatorKey][]float64 {
	highPrices := make([]float64, len(data))
	lowPrices := make([]float64, len(data))
	closePrices := make([]float64, len(data))
	for i, entry := range data {
		highPrices[i] = entry.High
		lowPrices[i] = entry.Low
		closePrices[i] = entry.Close
	}

	cache := make(map[IndicatorKey][]float64)

	processCondition := func(name string, period int) {
		key := IndicatorKey{Name: name, Period: period}
		if _, exists := cache[key]; exists {
			return
		}

		switch name {
		case "SMA":
			cache[key] = talib.Sma(closePrices, period)
		case "RSI":
			cache[key] = talib.Rsi(closePrices, period)
		case "WILLR":
			cache[key] = talib.WillR(highPrices, lowPrices, closePrices, period)
		case "Data":
			cache[key] = closePrices
		}
	}

	// Loop through all buy conditions
	for _, cond := range buyScenario.Conditions {
		processCondition(cond.IndicatorName, cond.IndicatorPeriod)

		cv := cond.IndicatorCheckValue
		if cv.IndicatorName != "" && cv.IndicatorPeriod > 0 {
			processCondition(cv.IndicatorName, cv.IndicatorPeriod)
		}
	}

	// Loop through all sell conditions
	for _, cond := range sellScenario.Conditions {
		if cond.ConditionType != SellIndicator {
			continue
		}

		processCondition(cond.IndicatorName, cond.IndicatorPeriod)

		cv := cond.IndicatorCheckValue
		if cv.IndicatorName != "" && cv.IndicatorPeriod > 0 {
			processCondition(cv.IndicatorName, cv.IndicatorPeriod)
		}
	}

	return cache
}
