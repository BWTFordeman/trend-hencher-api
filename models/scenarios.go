package models

import (
	"encoding/json"
	"log"
	"os"
)

// ScenarioConfig represents a complete trading scenario configuration
type ScenarioConfig struct {
	Name                  string
	IndicatorBuyScenario  BuyScenario
	IndicatorSellScenario SellScenario
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
