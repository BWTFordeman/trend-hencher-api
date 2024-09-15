package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"trend-hencher-api/models"
	"trend-hencher-api/services"
	"trend-hencher-api/utils"
)

type TrendHandler struct {
	trendService *services.TrendService
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

func NewTrendHandler(trendService *services.TrendService) *TrendHandler {
	return &TrendHandler{trendService: trendService}
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
	if r.Method != http.MethodGet {
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
		http.Error(w, "Failed to retrieve or parse data", http.StatusInternalServerError)
		return
	}
	log.Println("data: ", intradayData)

	// TODO Setup trends:

	// TODO Run trends to get scores:

	// TODO save trends with score:
	/* trendsCreated, err := h.trendService.CreateTrends()
	if err != nil {
		http.Error(w, "No trends created", http.StatusNotFound)
		return
	} */

	utils.WriteJSON(w, http.StatusOK, nil)
}

func fetchIntradayData(stockSymbol string) ([]IntradayData, error) {
	environment := os.Getenv("ENVIRONMENT")

	if environment == "local" {
		fmt.Println("Fetching local test data")
		return fetchFromLocalFile()
	}

	return fetchFromAPI(stockSymbol)
}

func fetchFromLocalFile() ([]IntradayData, error) {
	filePath := "test-data.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read local file: %v", err)
	}

	var intradayData []IntradayData
	err = json.Unmarshal(data, &intradayData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse local JSON file: %v", err)
	}

	return intradayData, nil
}

func fetchFromAPI(stockSymbol string) ([]IntradayData, error) {
	apiToken := os.Getenv("EODHD_API_TOKEN")
	if apiToken == "" {
		return nil, fmt.Errorf("API token is not set")
	}

	url := fmt.Sprintf("https://eodhd.com/api/intraday/%s.US?interval=1m&api_token=%s&fmt=json", stockSymbol, apiToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Go-http-client/1.1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch intraday data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data, status code: %d", resp.StatusCode)
	}

	var intradayData []IntradayData
	err = json.NewDecoder(resp.Body).Decode(&intradayData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return intradayData, nil
}
