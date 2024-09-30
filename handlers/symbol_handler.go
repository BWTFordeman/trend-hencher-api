package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Fetch intraday data for a symbol chosen, either from API or testdata retrieved locally.
func fetchIntradayData(stockSymbol string) ([]IntradayData, error) {
	environment := os.Getenv("ENVIRONMENT")

	if environment == "local" {
		fmt.Println("Fetching local test data")
		return fetchFromLocalFile()
	}

	return fetchFromAPI(stockSymbol)
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
