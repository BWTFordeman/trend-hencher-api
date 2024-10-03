package handlers

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"
	"trend-hencher-api/models"
	"trend-hencher-api/services"
	"trend-hencher-api/utils"

	"github.com/google/uuid"
	"github.com/markcheno/go-talib"
)

type TrendHandler struct {
	trendService         *services.TrendService
	bigQueryTrendService *services.BigQueryTrendService
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

func NewTrendHandler(trendService *services.TrendService, bigQueryTrendService *services.BigQueryTrendService) *TrendHandler {
	return &TrendHandler{
		trendService:         trendService,
		bigQueryTrendService: bigQueryTrendService,
	}
}

func (h *TrendHandler) GetTrend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queryID := r.URL.Query().Get("id")
	if queryID == "" {
		http.Error(w, "Missing trend ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(queryID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid trend ID", http.StatusBadRequest)
		return
	}

	trend, err := h.trendService.GetTrendByID(id)
	if err != nil {
		http.Error(w, "Trend not found", http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, http.StatusOK, trend)
}

func (h *TrendHandler) GetAllTrends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	trends, err := h.trendService.GetAllTrends()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve trends"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, trends)
}

func (h *TrendHandler) SaveTrend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var trend models.Trend
	err := json.NewDecoder(r.Body).Decode(&trend)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	trend.Date = time.Now() // Set the current date and time

	_, err = h.trendService.SaveTrend(&trend)
	if err != nil {
		http.Error(w, "Failed to save trend", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Trend saved successfully"})
}

func (h *TrendHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queryID := r.URL.Query().Get("id")
	if queryID == "" {
		http.Error(w, "Missing trend ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(queryID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid trend ID", http.StatusBadRequest)
		return
	}

	transactions, err := h.trendService.GetTransactions(id)
	if err != nil {
		http.Error(w, "Transactions not found", http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, http.StatusOK, transactions)
}

func (h *TrendHandler) CheckMarket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stockSymbol := r.URL.Query().Get("symbol")
	if stockSymbol == "" {
		http.Error(w, "Missing stock symbol", http.StatusBadRequest)
		return
	}

	// Fetch data from API:
	log.Println("Checking market...")
	intradayData, err := fetchIntradayData(stockSymbol)
	if err != nil {
		http.Error(w, "Failed to retrieve or parse data", http.StatusUnauthorized)
		return
	}

	// Run trends:
	trendsCreated, err := createTrends(h, w, intradayData, stockSymbol)
	if err != nil {
		http.Error(w, "Failed creating trends", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, trendsCreated)
}

func createTrends(h *TrendHandler, w http.ResponseWriter, data []IntradayData, symbol string) (string, error) {
	err := createSingleTrends(h, w, data, symbol)
	// TODO add single trends for other indicators than SMA...

	// createDoubleTrends()
	// createCrossOverDoubleTrends()
	// createTripleTrends()
	// createComplexTrends()

	// TODO get more indicators

	if err != nil {
		return "Failed to create trends", err
	}

	return "Created trends", nil
}

func createSingleTrends(h *TrendHandler, w http.ResponseWriter, data []IntradayData, symbol string) error {

	trendID := uuid.New().String()
	indicatorBuyScenario := models.BuyScenario{
		Conditions: []models.BuyCondition{
			{
				IndicatorName:       "SMA",
				IndicatorType:       models.IndicatorUnder,
				IndicatorPeriod:     14,
				IndicatorCheckValue: "data",
			},
		},
	}
	log.Println("some error")
	// Another possible sellScenario could be lossThreshold set to recent low and profit to 2x loss. (If the recent low is very high compared to current then maybe no trade?)
	indicatorSellScenario := models.SellScenario{
		ProfitThreshold: 1.07,
		LossThreshold:   0.96,
	}
	trendScore, transactions, err := scoreTrend(data, indicatorBuyScenario, indicatorSellScenario, trendID)
	if err != nil {
		log.Println("scoring error", trendScore)
		return err
	}
	log.Println("score gotten", trendScore)
	trend := models.Trend{
		TrendID:               trendID,
		Stock:                 symbol,
		TrendScore:            trendScore,
		Date:                  time.Now(),
		IndicatorBuyScenario:  indicatorBuyScenario,
		IndicatorSellScenario: indicatorSellScenario,
	}
	// Save Trend
	err = h.bigQueryTrendService.SaveTrend(&trend)
	if err != nil {
		return err
	}
	log.Println("Trend saved: ", trend)

	// Generate transactionIDs
	for i := range transactions {
		transactions[i].TransactionID = uuid.New().String()
	}

	// Save transactions
	err = h.bigQueryTrendService.SaveTransactions(transactions)
	if err != nil {
		return err
	}

	return nil
}

func scoreTrend(data []IntradayData, buyScenario models.BuyScenario, sellScenario models.SellScenario, trendID string) (float64, []models.Transaction, error) {
	transactions := []models.Transaction{}
	inPosition := false
	var lastBuy models.Transaction

	// Get indicator data:
	var closePrices []float64
	for _, entry := range data {
		closePrices = append(closePrices, entry.Close)
	}
	SMAData := talib.Sma(closePrices, buyScenario.Conditions[0].IndicatorPeriod)

	// Get transactions:
	for i := 1; i < len(data); i++ {
		price := data[i].Close

		// Check for BuyScenario
		if !inPosition && checkBuyScenario(buyScenario, closePrices, SMAData, i) {
			lastBuy = models.Transaction{
				DateBought:  data[i].Datetime,
				PriceBought: price,
				Volume:      int64(price / 1000000), // Assuming total invested per trade is 1.000.000~
			}
			transactions = append(transactions, lastBuy)
			inPosition = true
		}

		if inPosition && checkSellScenario(sellScenario, lastBuy.PriceBought, price) {
			lastBuy.DateSold = data[i].Datetime
			lastBuy.PriceSold = price
			lastBuy.TrendID = trendID
			transactions[len(transactions)-1] = lastBuy
			inPosition = false
		}
	}

	// Remove last transaction that was not sold:
	if inPosition {
		transactions = transactions[:len(transactions)-1]
	}

	trendScore := calculateTrendScore(transactions)

	return trendScore, transactions, nil
}

func checkBuyScenario(buyScenario models.BuyScenario, checkData []float64, indicatorData []float64, index int) bool {
	for _, condition := range buyScenario.Conditions {
		switch condition.IndicatorType {
		case models.IndicatorCrossDown:
			if checkData[index] >= indicatorData[index] {
				return false
			}
		case models.IndicatorCrossUp:
			if checkData[index] <= indicatorData[index] {
				return false
			}
		case models.IndicatorOver:
			if checkData[index] > indicatorData[index] {
				return false
			}
		case models.IndicatorUnder:
			if checkData[index] < indicatorData[index] {
				return false
			}
		}
	}
	return true
}

func checkSellScenario(sellScenario models.SellScenario, buyPrice, currentPrice float64) bool {
	if currentPrice >= buyPrice*(1+sellScenario.ProfitThreshold/100) ||
		currentPrice <= buyPrice*(1-sellScenario.LossThreshold/100) {
		return true
	}
	return false
}

// Fake normalizations are being done - meaning any trend can have a score above 1
// but most won't. When they go above 1 they are most likely very good trends!
func calculateTrendScore(transactions []models.Transaction) float64 {

	// Occurrence (assuming a max of 100)
	normalizedOccurrence := len(transactions) / 100
	occurrenceWeight := 0.15

	// Profitability (Assuming a max of 100%, and total invested: 1.000.000~ per trade)
	totalProfit := 0.0
	for _, transaction := range transactions {
		profit := (transaction.PriceSold - transaction.PriceBought) * float64(transaction.Volume)
		totalProfit += profit
	}
	normalizedProfitability := totalProfit / 1000000
	profitabilityWeight := 0.45

	// Consistency
	winningTransactions := 0
	totalTransactions := len(transactions)
	normalizedConsistency := 0.0

	for _, transaction := range transactions {
		if transaction.PriceSold > transaction.PriceBought {
			winningTransactions++
		}
	}

	if totalTransactions != 0 {
		normalizedConsistency = float64(winningTransactions) / float64(totalTransactions)
	}
	consistencyWeight := 0.25

	// Variance
	normalizedVariance := calculateVarianceScore(transactions)
	varianceWeight := 0.15

	trendScore := occurrenceWeight*float64(normalizedOccurrence) + profitabilityWeight*normalizedProfitability + consistencyWeight*normalizedConsistency + varianceWeight*normalizedVariance
	return math.Round(trendScore*1000) / 1000
}

func calculateVarianceScore(transactions []models.Transaction) float64 {
	var percentageProfits []float64

	// Calculate percentage profit for each transaction
	for _, transaction := range transactions {
		percentageProfit := ((transaction.PriceSold - transaction.PriceBought) / transaction.PriceBought) * 100
		percentageProfits = append(percentageProfits, percentageProfit)
	}

	averageProfit := calculateAverage(percentageProfits)
	medianProfit := calculateMedian(percentageProfits)

	// Calculate the variance (difference between average and median percentage profit)
	variance := math.Abs(averageProfit - medianProfit)

	maxVariance := 1.0
	varianceScore := 1 - (variance / maxVariance)
	if varianceScore < 0 {
		varianceScore = 0 // Ensure the score doesn't go below 0
	}
	return varianceScore
}

func calculateMedian(profits []float64) float64 {
	sort.Float64s(profits)

	n := len(profits)
	if n%2 == 0 {
		return (profits[n/2-1] + profits[n/2]) / 2.0
	}
	return profits[n/2]
}

func calculateAverage(profits []float64) float64 {
	sum := 0.0
	for _, profit := range profits {
		sum += profit
	}
	return sum / float64(len(profits))
}
