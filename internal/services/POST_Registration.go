package services

import (
	"Assignment2/internal/constants"
	"Assignment2/internal/database"
	"Assignment2/internal/structs"
	"Assignment2/internal/utils"
	"encoding/json"
	"net/http"
	"strings"
	"context"
	"time"
	"log"
	"io"
)

func POST_Registration(w http.ResponseWriter, r *http.Request) {
	// Reading body:
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Reading payload from body failed.")
		http.Error(w, "Reading payload failed.", http.StatusInternalServerError)
		return
	}

	log.Println("\nReceived request to add document for content ", string(content))

	// Check that a body was acutally sent.
	if len(string(content)) == 0 {
		log.Println("Content appears to be empty.")
		http.Error(w, "Payload is empty.", http.StatusBadRequest)
		return		
	}

	// Parsing JSON into struct:
	var config structs.Configuration
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Println("Error parsing JSON:", err)
		http.Error(w, "Invalid JSON format.", http.StatusBadRequest)
		return
	}

	// To Uppering ISOCode for consistency
	config.ISOCode = strings.ToUpper(config.ISOCode)

	// Checking that the isoCode is valid and setting country name:
	// Building the API URL:
	url := constants.Country_API + config.ISOCode + "?fields=name"

	// Fetching from Country_API API:
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Could not fetch from API")
		http.Error(w, "Wrong IsoCode?", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	// Checking response status:
	if resp.StatusCode != http.StatusOK {
		log.Println("API arror from: ", constants.Country_API)
		http.Error(w, "Error Fetching from API", http.StatusInternalServerError)
		return
	}

	// Reading body:
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not read response from: ", config.ISOCode)
		http.Error(w, "Could not read JSON body", http.StatusInternalServerError)
		return
	}

	// Decoding JSON into map[string]inteface{}:
	var countryData map[string]interface{}
	if err := json.Unmarshal(body, &countryData); err != nil {
		log.Println("Error Decoding JSON. API problem?")
		http.Error(w, "Error Decoding JSON", http.StatusInternalServerError)
		return
	}

	// Setting country name:
	Name := countryData["name"].(map[string]interface{})
	config.Country = Name["common"].(string)

	// Checking if Target Currencies is not set:
	if config.Features.TargetCurrencies == nil {
		config.Features.TargetCurrencies = []string{}
	}

	// Setting current time:
	config.LastChange = time.Now().Format("20060102 15:04")

	// Setting up firestore client:
	ctx := context.Background()
	client := database.GetFirebaseClient()

	// Adding document to firebase:
	collection := client.Collection(constants.Registration)
	docRef, _, err := collection.Add(ctx, config)
	if err != nil {
		log.Println("Error writing to firebase:", err)
		http.Error(w, "Failed to add document", http.StatusInternalServerError)
		return
	}

	// Invoke webhooks:
	utils.Invoke("REGISTER", config.ISOCode, "ID: " + docRef.ID +
	" REGISTERD with isoCode : " + config.ISOCode)

	// Creating temporary response struct:
	response:= map[string]string{
		"id":	      docRef.ID,
		"lastChange": config.LastChange,
	}

	// Converting to JSON:
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create response.", http.StatusInternalServerError)
		return
	}

	// Response:
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}
