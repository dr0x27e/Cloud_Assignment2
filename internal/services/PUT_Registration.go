package services

import (
	"Assignment2/internal/constants"
	"Assignment2/internal/database"
	"Assignment2/internal/structs"
	"Assignment2/internal/utils"
	"encoding/json"
	"net/http"
	"context"
	"strings"
	"time"
	"log"
	"io"
)

// Update registrations function:
func PUT_Registration(w http.ResponseWriter, r* http.Request) {
	// Reading body:
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Reading payload from body failed.")
		http.Error(w, "Reading payload failed.", http.StatusInternalServerError)
		return
	}

	log.Println("\nReceived request to add document for content\n", string(content))

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

	// Turning isoCode into to upper for consistency:
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

	// Decoding JSON into struct:
	var countryData map[string]interface{}
	if err := json.Unmarshal(body, &countryData); err != nil {
		log.Println("Error Decoding JSON. API problem?")
		http.Error(w, "Error Decoding JSON", http.StatusInternalServerError)
		return
	}

	Name := countryData["name"].(map[string]interface{})
	Common := Name["common"].(string)

	// Checking if an acutal name was sent back:
	if Common == "" {
		log.Println("Not An Acutal Country")
		http.Error(w, "Wrong Country name", http.StatusBadRequest)
		return
	}

	// Extracting the document id:
	docID := r.PathValue("id")
	if docID == "" {
		log.Println("No ID given")
		http.Error(w, "Missing document ID", http.StatusBadRequest)
		return
	}

	// Firestore client:
	ctx := context.Background()
	client := database.GetFirebaseClient()

	// Referance Firestore document:
	docRef := client.Collection(constants.Registration).Doc(docID)

	// Checking if it is a valid id:
	_, err = docRef.Get(ctx)
	if err != nil {
		log.Println("Invalid document ID")
		http.Error(w, "Document does not exist (invalid document ID)", http.StatusNotFound)
		return
	}

	// Setting country name (Even if it is set right or not)
	config.Country = Common

	// Updating LastChange variable:
	config.LastChange = time.Now().Format("20060102 15:04")

	// Updateing document:
	_, err = docRef.Set(ctx, config)
	if err != nil {
		log.Println("Error updating document:", err)
		http.Error(w, "Failed to update document", http.StatusInternalServerError)
		return
	}
	
	// Invoke webhooks:
	utils.Invoke("CHANGE", config.ISOCode, "ID: " + docRef.ID + 
	" Got edited method: PUT")

	// Successful response:
	w.WriteHeader(http.StatusCreated)
}
