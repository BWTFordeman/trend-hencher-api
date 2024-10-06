package handlers

import (
	"os"
	"testing"
)

func TestFetchIntradayData(t *testing.T) {
	os.Setenv("ENVIRONMENT", "local")

	data, err := fetchIntradayData("APPL")
	if err != nil {
		t.Errorf("fetchIntradayData should not give error; got: %s", err.Error())
	}

	// Validate returned data
	if data == nil {
		t.Error("Expected non-nil data from fetchIntradayData")
	}
}

func TestFetchIntradayDataWithdifferentSymbols(t *testing.T) {
	os.Setenv("ENVIRONMENT", "local")

	_, err := fetchIntradayData("MSFT")
	if err != nil {
		t.Errorf("fetchIntradayData should not give error; got: %s", err.Error())
	}

	// Can setup new request to check for valid symbol through:
	// https://eodhd.com/api/exchanges-list/?api_token={YOUR_API_TOKEN}&fmt=json
	/* _, err = fetchIntradayData("INVALID")
	if err == nil {
		t.Errorf("fetchIntradayData should give error but didn't get any")
	} */
}

func TestFetchIntradayDataWithProdData(t *testing.T) {
	os.Setenv("ENVIRONMENT", "production")

	_, err := fetchIntradayData("APPL")
	if err == nil {
		t.Errorf("fetchIntradayData should give error but didn't get any")
	}
}

func TestFetchIntradayDataWithoutEnvironmentENV(t *testing.T) {
	_, err := fetchIntradayData("APPL")
	if err == nil {
		t.Errorf("fetchIntradayData should give error but didn't get any")
	}
}
