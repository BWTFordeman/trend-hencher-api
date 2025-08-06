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
