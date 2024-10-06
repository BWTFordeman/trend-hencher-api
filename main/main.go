package main

import (
	"log"
	"net/http"
)

func main() {
	initEnv()
	trendHandler := initTrendHandler()

	log.Println("Starting server on :8082..")
	http.HandleFunc("/checkMarket", trendHandler.CheckMarket)
	http.HandleFunc("/trend", trendHandler.GetTrend)
	http.HandleFunc("/trends", trendHandler.GetAllTrends)
	http.HandleFunc("/saveTrend", trendHandler.SaveTrend)
	http.HandleFunc("/transactions", trendHandler.GetTransactions)

	log.Fatal(http.ListenAndServe(":8082", nil))
}
