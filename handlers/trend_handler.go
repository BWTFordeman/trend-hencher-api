package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"trend-hencher-api/models"
	"trend-hencher-api/services"
	"trend-hencher-api/utils"
)

type TrendHandler struct {
	trendService *services.TrendService
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
