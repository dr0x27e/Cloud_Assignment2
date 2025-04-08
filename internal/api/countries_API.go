package api

import (
	"Assignment2/internal/constants"
	"encoding/json"
	"net/http"
	"fmt"
	"io"
)

func Countries(Iso2 string) (map[string]interface{}, string, error) {
	// Building the API URL:
	url := constants.Country_API + Iso2 + "?fields=currencies,capital,latlng,population,area"

	// Fetching from Currency_API API:
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("could not fetch %s: %v", Iso2, err)
	}
	defer resp.Body.Close()

	// Checking response status:
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("API error for %s: %s", Iso2, resp.Status)
	}

	// Reading body:
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("could not read response for %s: %v", Iso2, err)
	}

	// Decode JSON into map[string]interface{}:
	var respData map[string]interface{}
	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, "", fmt.Errorf("error decoding JSON for %s: %v", Iso2, err)
	}

	// Extract first currency from the "currencies" map (if present):
	var firstCurrency string
	if currencies, ok := respData["currencies"].(map[string]interface{}); ok {
		for currencyCode := range currencies {
			firstCurrency = currencyCode
			break
		}
	}

	// Remove "currencies" field from response map:
	delete(respData, "currencies")

	return respData, firstCurrency, nil
}
