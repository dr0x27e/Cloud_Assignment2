package api

import (
    "Assignment2/internal/constants"
    "Assignment2/internal/structs"
    "encoding/json"
    "net/http"
    "fmt"
    "io"
)

// Currencies fetches exchange rates for a main currency against target currencies
func Currencies(main string, targets []string) (map[string]float64, error) {
    	// Building the API URL:
    	url := constants.Currency_API + main

    	// Fetching from Currency_API API:
    	resp, err := http.Get(url)
    	if err != nil {
        	return nil, fmt.Errorf("could not fetch rates for %s: %v", main, err)
    	}
    	defer resp.Body.Close()

    	// Checking response status:
    	if resp.StatusCode != http.StatusOK {
        	return nil, fmt.Errorf("API error for %s: %s", main, resp.Status)
    	}

    	// Reading body
    	body, err := io.ReadAll(resp.Body)
    	if err != nil {
        	return nil, fmt.Errorf("could not read response for %s: %v", main, err)
    	}

    	// Decode JSON into ExchangeRates struct
    	var respData structs.ExchangeRates
    	if err := json.Unmarshal(body, &respData); err != nil {
        	return nil, fmt.Errorf("error decoding JSON for %s: %v", main, err)
    	}

    	// Building the result map with only target currencies that exist:
    	result := make(map[string]float64)
    	for _, target := range targets {
        	if rate, exists := respData.Rates[target]; exists {
            	result[target] = rate
        	}
    	}

    	
    return result, nil
}
