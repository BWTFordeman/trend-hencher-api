package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"
)

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func MeasureTime(start time.Time, name string) {
	elapsed := time.Since(start)
	seconds := elapsed.Seconds()                  // Get the duration in seconds (as float64)
	milliseconds := elapsed.Milliseconds() % 1000 // Get the milliseconds part only

	log.Printf("%s took %.0f s %d ms", name, seconds, milliseconds)
}

// Convert Unix timestamp and GMT offset to time.Time
func ConvertToEasternTime(unixTime int64, gmtOffset int) time.Time {
	utcTime := time.Unix(unixTime, 0)

	// Apply the GMT offset (adjusting for market-provided timezone)
	offsetDuration := time.Duration(gmtOffset) * time.Second
	adjustedTime := utcTime.Add(offsetDuration)

	// Convert to Eastern Time (ET)
	easternLocation, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("Failed to load Eastern timezone: %v", err)
	}

	return adjustedTime.In(easternLocation)
}

func CalculateAverage(profits []float64) float64 {
	sum := 0.0
	for _, profit := range profits {
		sum += profit
	}
	return sum / float64(len(profits))
}

func CalculateMedian(profits []float64) float64 {
	sort.Float64s(profits)

	n := len(profits)
	if n%2 == 0 {
		return (profits[n/2-1] + profits[n/2]) / 2.0
	}
	return profits[n/2]
}
