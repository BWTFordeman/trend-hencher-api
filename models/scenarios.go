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
						IndicatorName:   "SMA",
						IndicatorType:   IndicatorUnder,
						IndicatorPeriod: 14,
						IndicatorCheckValue: Indicator{
							IndicatorName: "data", // Buy when SMA is under data value
						},
					},
				},
			},
			IndicatorSellScenario: SellScenario{
				Conditions: []SellCondition{
					{
						ConditionType:   SellPercentage,
						ProfitThreshold: 1.07,
						LossThreshold:   0.96,
					},
				},
			},
		},
		{
			Name: "SMA_20_Under",
			IndicatorBuyScenario: BuyScenario{
				Conditions: []BuyCondition{
					{
						IndicatorName:   "SMA",
						IndicatorType:   IndicatorUnder,
						IndicatorPeriod: 20,
						IndicatorCheckValue: Indicator{
							IndicatorName: "data",
						},
					},
				},
			},
			IndicatorSellScenario: SellScenario{
				Conditions: []SellCondition{
					{
						ConditionType:   SellPercentage,
						ProfitThreshold: 1.05,
						LossThreshold:   0.97,
					},
				},
			},
		},
		{
			Name: "RSI14_Under30",
			IndicatorBuyScenario: BuyScenario{
				Conditions: []BuyCondition{
					{
						IndicatorName:   "RSI",
						IndicatorType:   IndicatorUnder,
						IndicatorPeriod: 14,
						IndicatorCheckValue: Indicator{
							IndicatorName:        "RSI",
							IndicatorRSIStrength: 30, // Buy when RSI < 30 (oversold)
						},
					},
				},
			},
			IndicatorSellScenario: SellScenario{
				Conditions: []SellCondition{
					{
						ConditionType:   SellPercentage,
						ProfitThreshold: 1.06,
						LossThreshold:   0.97,
					},
				},
			},
		},
		{
			Name: "RSI14_CrossOver50",
			IndicatorBuyScenario: BuyScenario{
				Conditions: []BuyCondition{
					{
						IndicatorName:   "RSI",
						IndicatorType:   IndicatorCrossUp,
						IndicatorPeriod: 14,
						IndicatorCheckValue: Indicator{
							IndicatorName:        "RSI",
							IndicatorRSIStrength: 50, // Buy when RSI crosses above 50
						},
					},
				},
			},
			IndicatorSellScenario: SellScenario{
				Conditions: []SellCondition{
					{
						ConditionType:   SellPercentage,
						ProfitThreshold: 1.08,
						LossThreshold:   0.95,
					},
				},
			},
		},
		// Add more predefined scenarios here
	}
}
