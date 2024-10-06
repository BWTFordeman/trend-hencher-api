package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"trend-hencher-api/utils"
)

// Fetch intraday data for a symbol chosen, either from API or testdata retrieved locally.
func fetchIntradayData(stockSymbol string) ([]IntradayData, error) {
	defer utils.MeasureTime(time.Now(), "fetchIntradayData")

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
	filePath := "../testdata/test-data.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read local file: %v", err)
	}

	var intradayData []IntradayData
	err = json.Unmarshal(data, &intradayData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse local JSON file: %v", err)
	}

	// Set correct timezone for intraday data:
	intradayData = filterIntradayData(intradayData)

	return intradayData, nil
}

func filterIntradayData(intradayData []IntradayData) []IntradayData {
	// Load the Oslo timezone
	osloLocation, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		log.Fatalf("Failed to load Oslo timezone: %v", err)
	}

	filteredData := []IntradayData{}

	for _, data := range intradayData {
		// Convert the timestamp to Eastern Time
		easternTime := convertToEasternTime(data.Timestamp, data.GmtOffset)

		// Check if the time falls within the open market hours (9:30 AM - 4:00 PM ET)
		if (easternTime.Hour() > openHourET || (easternTime.Hour() == openHourET && easternTime.Minute() >= openMinuteET)) &&
			easternTime.Hour() < closeHourET {

			// Convert the timestamp to Oslo time
			osloTime := easternTime.In(osloLocation)

			// Add the data to the filtered list, but adjust the timestamp and datetime for Oslo time
			filteredData = append(filteredData, IntradayData{
				Timestamp: osloTime.Unix(), // Convert back to Unix timestamp if needed
				GmtOffset: 3600,            // Set Oslo GMT offset manually (+1 hour standard, +2 hours DST)
				Datetime:  osloTime.Format("2006-01-02 15:04:05"),
				Open:      data.Open,
				High:      data.High,
				Low:       data.Low,
				Close:     data.Close,
				Volume:    data.Volume,
			})
		}
	}

	return filteredData
}
