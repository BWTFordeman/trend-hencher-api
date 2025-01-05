package models

// ScenarioConfig represents a complete trading scenario configuration
type ScenarioConfig struct {
	Name                  string
	IndicatorBuyScenario  BuyScenario
	IndicatorSellScenario SellScenario
}

// GetPredefinedScenarios returns a list of all predefined trading scenarios
func GetPredefinedSingleTrendScenarios() []ScenarioConfig {
	return []ScenarioConfig{
		{
			Name: "SMA_14_Under",
			IndicatorBuyScenario: BuyScenario{
				Conditions: []BuyCondition{
					{
						IndicatorName:       "SMA",
						IndicatorType:       IndicatorUnder,
						IndicatorPeriod:     14,
						IndicatorCheckValue: "data",
					},
				},
			},
			IndicatorSellScenario: SellScenario{
				ProfitThreshold: 1.07,
				LossThreshold:   0.96,
			},
		},
		{
			Name: "SMA_20_Under",
			IndicatorBuyScenario: BuyScenario{
				Conditions: []BuyCondition{
					{
						IndicatorName:       "SMA",
						IndicatorType:       IndicatorUnder,
						IndicatorPeriod:     20,
						IndicatorCheckValue: "data",
					},
				},
			},
			IndicatorSellScenario: SellScenario{
				ProfitThreshold: 1.05,
				LossThreshold:   0.97,
			},
		},
		// Add more predefined scenarios here
	}
}
