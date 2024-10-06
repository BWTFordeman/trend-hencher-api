package utils

import (
	"encoding/json"
	"log"
	"net/http"
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
