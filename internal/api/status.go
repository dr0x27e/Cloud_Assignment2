package api

import (
	"Assignment2/internal/constants"
	"Assignment2/internal/database"
	"Assignment2/internal/structs"
	"encoding/json"
	"net/http"
	"context"
	"time"
	"log"
)

// Global variable to store service start time
var serviceStartTime = time.Now()

// checkAPIStatus sends a GET request to an external API and returns its status code.
func checkAPIStatus(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error reaching API:", url, "-", err)
		return http.StatusServiceUnavailable // 503 if the API is unreachable
	}
	return resp.StatusCode
}

func Status(w http.ResponseWriter, r *http.Request) {
	log.Println("Received status request.")

	ctx := context.Background()
	client := database.GetFirebaseClient()

	countriesStatus := checkAPIStatus(constants.TEST_COUNTRY)
	meteoStatus     := checkAPIStatus(constants.TEST_OPENMETEO)
	currencyStatus  := checkAPIStatus(constants.TEST_CURRENCY)
	pythonStatus    := checkAPIStatus(constants.TEST_PYTHON_API)

	// webhooks DataBase Status
	webhooks := http.StatusOK
	_, err := client.Collection(constants.Webhooks).Documents(ctx).GetAll()
	if err != nil {
		webhooks = http.StatusInternalServerError
	}

	// Registration data
	Registration := http.StatusOK
	_, err2 := client.Collection(constants.Webhooks).Documents(ctx).GetAll()
	if err2 != nil {
		Registration = http.StatusInternalServerError
	}

	// Webhook count directly inside the function
	webhookCount := 0

	for _, event := range constants.AllEvents {
		collection := client.Collection(
			constants.Webhooks).Doc(
				event).Collection(
					constants.SubCollection)

		iter := collection.Documents(ctx)
		// Iterating through collection to count documents:
		for {
			_, err := iter.Next()
			if err != nil {
				break // No more documents
			}
			webhookCount++
		}
	}

	// Compute uptime
	uptime := int64(time.Since(serviceStartTime).Seconds())

	// Build response
	statusResponse := structs.StatusResponse{
		CountriesAPI:         countriesStatus,
		MeteoAPI:             meteoStatus,
		CurrencyAPI:          currencyStatus,
		PythonAPI:            pythonStatus,
		WebhooksDatabase:     webhooks,
		RegistrationDatabase: Registration,
		Webhooks:             webhookCount,
		Version:              "v1",
		Uptime:               uptime,
	}

	// Convert to JSON
	jsonResponse, err := json.Marshal(statusResponse)
	if err != nil {
		http.Error(w, "Error generating JSON response", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}


