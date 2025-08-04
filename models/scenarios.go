package models

import (
	"encoding/json"
	"log"
	"os"

	"github.com/markcheno/go-talib"
)

// ScenarioConfig represents a complete trading scenario configuration
type ScenarioConfig struct {
	Name                  string
	IndicatorBuyScenario  BuyScenario
	IndicatorSellScenario SellScenario
}

type IntradayData struct {
	Timestamp int64   `json:"timestamp"`
	GmtOffset int     `json:"gmtoffset"`
	Datetime  string  `json:"datetime"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int     `json:"volume"`
}

type IndicatorKey struct {
	Name   string
	Period int
}

// GetPredefinedScenarios returns a list of all predefined trading scenarios
func GetPredefinedScenarios() []ScenarioConfig {

	scenarios, err := LoadScenarioConfigs("scenarios.json")
	if err != nil {
		log.Fatal("Failed to load scenarios:", err)
	}
	return scenarios
}

func LoadScenarioConfigs(filepath string) ([]ScenarioConfig, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var scenarios []ScenarioConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&scenarios); err != nil {
		return nil, err
	}

	return scenarios, nil
}

func GetPredefinedIndicators(buyScenario BuyScenario, data []IntradayData) map[IndicatorKey][]float64 {
	closePrices := make([]float64, len(data))
	for i, entry := range data {
		closePrices[i] = entry.Close
	}

	cache := make(map[IndicatorKey][]float64)

	for _, cond := range buyScenario.Conditions {
		// Source
		key := IndicatorKey{cond.IndicatorName, cond.IndicatorPeriod}
		if _, exists := cache[key]; !exists {
			switch cond.IndicatorName {
			// TODO setup more indicators
			case "SMA":
				cache[key] = talib.Sma(closePrices, cond.IndicatorPeriod)
			case "RSI":
				cache[key] = talib.Rsi(closePrices, cond.IndicatorPeriod)
			case "Data":
				cache[key] = closePrices
			}
		}

		// Target (CheckValue)
		cv := cond.IndicatorCheckValue
		checkKey := IndicatorKey{cv.IndicatorName, cv.IndicatorPeriod}
		if _, exists := cache[checkKey]; !exists {
			switch cv.IndicatorName {
			case "SMA":
				cache[checkKey] = talib.Sma(closePrices, cv.IndicatorPeriod)
			case "RSI":
				cache[checkKey] = talib.Rsi(closePrices, cv.IndicatorPeriod)
			case "Data":
				cache[checkKey] = closePrices
			}
		}
	}

	return cache
}
