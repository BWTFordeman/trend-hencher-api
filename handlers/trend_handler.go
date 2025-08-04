package handlers

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
	"trend-hencher-api/models"
	"trend-hencher-api/services"
	"trend-hencher-api/utils"

	"github.com/google/uuid"
)

type TrendHandler struct {
	trendService         *services.TrendService
	bigQueryTrendService *services.BigQueryTrendService
}

const (
	openHourET   = 9  // Market opens at 9 AM ET
	openMinuteET = 30 // Market opens at 9:30 AM ET
	closeHourET  = 16 // Market closes at 4 PM ET
)

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
		log.Printf("Error fetching data; %v", err)
		http.Error(w, "Failed to retrieve or parse data", http.StatusUnauthorized)
		return
	}

	// Run trends:
	trendsCreated, err := createTrends(h, intradayData, stockSymbol)
	if err != nil {
		http.Error(w, "Failed creating trends", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, trendsCreated)
}

func createTrends(h *TrendHandler, data []models.IntradayData, symbol string) (string, error) {
	defer utils.MeasureTime(time.Now(), "createTrends")
	err := createSingleTrends(h, data, symbol)
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

func createSingleTrends(h *TrendHandler, data []models.IntradayData, symbol string) error {
	// Get all predefined scenarios
	scenarios := models.GetPredefinedScenarios()

	// Run through each scenario
	for _, scenario := range scenarios {
		trendID := uuid.New().String()

		transactions, err := createTransactions(data, scenario.IndicatorBuyScenario, scenario.IndicatorSellScenario, trendID)

		trendScore := calculateTrendScore(transactions)
		if err != nil {
			log.Printf("scoring error for scenario %s: %v", scenario.Name, err)
			continue // Skip this scenario if there's an error
		}

		log.Println("score for scenario: ", trendScore)
		/*
			trend := models.Trend{
				TrendID:               trendID,
				Stock:                 symbol,
				TrendScore:            trendScore,
				YearlyProfit:			yearlyProfit,
				Date:                  time.Now(),
				IndicatorBuyScenario:  scenario.IndicatorBuyScenario,
				IndicatorSellScenario: scenario.IndicatorSellScenario,
			}

			// Save Trend
			if err := h.bigQueryTrendService.SaveTrend(&trend); err != nil {
				log.Printf("error saving trend for scenario %s: %v", scenario.Name, err)
				continue
			}

			// Generate transactionIDs
			for i := range transactions {
				transactions[i].TransactionID = uuid.New().String()
			}

			// Save transactions
			if err := h.bigQueryTrendService.SaveTransactions(transactions); err != nil {
				log.Printf("error saving transactions for scenario %s: %v", scenario.Name, err)
				continue
			} */

		log.Printf("Successfully processed scenario: %s", scenario.Name)
	}

	return nil
}

func createTransactions(data []models.IntradayData, buyScenario models.BuyScenario, sellScenario models.SellScenario, trendID string) ([]models.Transaction, error) {
	transactions := []models.Transaction{}
	inPosition := false
	var lastBuy models.Transaction

	indicatorCache := models.GetPredefinedIndicators(buyScenario, data)
	// I only have data value in sellscenario now so won't need to worry about indicators there atm.
	// TODO add sellScenario to GetPredefinedIndicators and use inidcatorCache in shouldSell (only required after adding indicators in shouldSell formula)

	// Get transactions:
	for i := 1; i < len(data); i++ {
		price := data[i].Close

		// Check for BuyScenario
		if !inPosition {
			if shouldBuy(data, buyScenario, i, indicatorCache) {
				lastBuy = models.Transaction{
					DateBought:  data[i].Datetime,
					PriceBought: price,
					Volume:      int64(1000000 / price), // Assuming total invested per trade is 1.000.000~
				}
			}
			transactions = append(transactions, lastBuy)
			inPosition = true
		}

		// Check for SellScenario
		if inPosition {
			if shouldSell(sellScenario, lastBuy.PriceBought, price) {
				lastBuy.DateSold = data[i].Datetime
				lastBuy.PriceSold = price
				lastBuy.TrendID = trendID
				transactions[len(transactions)-1] = lastBuy
				inPosition = false
			}
		}
	}

	// Remove last transaction that was not sold:
	if inPosition {
		transactions = transactions[:len(transactions)-1]
	}

	return transactions, nil
}

// This checks current data(by index) against buyScenario conditions and return whether to buy or wait for correct conditions to buy
func shouldBuy(data []models.IntradayData, buyScenario models.BuyScenario, index int, indicatorCache map[models.IndicatorKey][]float64) bool {
	for _, cond := range buyScenario.Conditions {
		indicatorSourceData := indicatorCache[models.IndicatorKey{Name: cond.IndicatorName, Period: cond.IndicatorPeriod}]
		indicatorTargetData := indicatorCache[models.IndicatorKey{Name: cond.IndicatorCheckValue.IndicatorName, Period: cond.IndicatorCheckValue.IndicatorPeriod}]
		if !checkBuyCondition(indicatorSourceData, indicatorTargetData, cond.IndicatorType, index) {
			return false
		}
	}
	return true
}

func checkBuyCondition(sourceData []float64, targetData []float64, indicatorType models.IndicatorType, index int) bool {

	currSource := sourceData[index]
	currTarget := targetData[index]

	switch indicatorType {
	case models.IndicatorOver:
		return currSource > currTarget
	case models.IndicatorUnder:
		return currSource < currTarget
	case models.IndicatorCrossUp:
		prevSource := sourceData[index-1]
		prevTarget := targetData[index-1]
		return prevSource < prevTarget && currSource >= currTarget
	case models.IndicatorCrossDown:
		prevSource := sourceData[index-1]
		prevTarget := targetData[index-1]
		return prevSource > prevTarget && currSource <= currTarget
	default:
		return false
	}
}

func shouldSell(sellScenario models.SellScenario, buyPrice, currentPrice float64) bool {
	for _, sellCondition := range sellScenario.Conditions {
		switch sellCondition.ConditionType {
		case models.SellPercentage:
			if currentPrice > buyPrice*(sellCondition.ProfitThreshold/100) || currentPrice < buyPrice*(sellCondition.LossThreshold/100) {
				return true
			}
		case models.SellIndicator:
			// TODO implement
			return true
			//TODO add more types
		}
	}

	return true
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
	log.Printf("totalProfit: %.2f", totalProfit)
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

	averageProfit := utils.CalculateAverage(percentageProfits)
	medianProfit := utils.CalculateMedian(percentageProfits)

	// Calculate the variance (difference between average and median percentage profit)
	variance := math.Abs(averageProfit - medianProfit)

	maxVariance := 1.0
	varianceScore := 1 - (variance / maxVariance)
	if varianceScore < 0 {
		varianceScore = 0 // Ensure the score doesn't go below 0
	}
	return varianceScore
}
