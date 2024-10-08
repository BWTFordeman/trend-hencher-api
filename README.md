# trend-hencher-api


## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Indicators](#indicators)
- [License](#license)
- [Scoring](#scoring)

## Installation
To install this project, follow these steps:

1. Clone the repo
```bash
git clone https://github.com/BWTFordeman/tred-hencher-api.git
```
2. Setup environment variables
#### Local variable
    $env:ENVIRONMENT="local"
#### Credentials variable
    $env:GOOGLE_APPLICATION_CREDENTIALS="location\service-account-file.json"
#### Cloud project ID
    $env:GOOGLE_CLOUD_PROJECT="id"     

3. Run the project
 ```bash
go run ./main
```

## Usage
After installing, run the following request to check the market for new trends for chosen symbol:
PUT - localhost:8080/checkMarket?symbol=APPL

Run the following request to retrieve all current trends found:
GET - localhost:8080/trends

## Indicators
Indicators used by the trend hencher to find the best trends:

- SMA (Simple Moving Average)

Other indicators to utilize:

- EMA (Exponential Moving Average)
- DEMA (Double Exponential Moving Average)
- TEMA (Triple Exponential Moving Average)
- RSI (Relative Strength Index)
- Stochastic
- Stochastic RSI
- Super Trend
- Acceleration Bands
- Bollinger Bands
- Linear Regression
- Williams %R
- MFI (Money Flow Index)
- WMA (Weighted Moving Average)
- APO (Absolute Price Oscillator)
- Aroon
- Aroon oscillator
- ATR (Average True Range)
- CCI (Commodity Channel Index)
- CMO (Chande Momentum Oscillator)
- MACD (Moving Average Convergence Divergence)
- MFI (Money Flow Index)
- MinusDI (Minus Directional Indicator)
- Mom (Momentum)
- OBV (On Balance Volume)
- PPO (Percentage Price Oscillator)
- ROC (Rate Of Change)
- Trix (1-day ROC of a triple SMA)

## Scoring
Trend scoring is based on a few factors:
- Profitability
- Consistency
- Variance (risk)

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
