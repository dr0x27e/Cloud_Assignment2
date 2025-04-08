package api

import (
	"Assignment2/internal/constants"
	"Assignment2/internal/database"
	"Assignment2/internal/utils"
	"encoding/json"
	"net/http"
	"context"
	"time"
)

func Dashboards(w http.ResponseWriter, r *http.Request) {
	// Checking if document id is given:
	docID := r.PathValue("id")

	if docID == "" {
		http.Error(w, "Missing document ID", http.StatusBadRequest)
		return
	}

	// Firestore client:
	ctx := context.Background()
	client := database.GetFirebaseClient()

	// Fetch the document:
	doc, err := client.Collection(constants.Registration).Doc(docID).Get(ctx)
	if err != nil {
		http.Error(w, "Invalid document", http.StatusNotFound)
		return
	}
	
	// Converting to map interface. fetching isocode and features:
	data := doc.Data()
	isoCode  := data["isoCode"].(string)
	features := data["features"].(map[string]interface{})

	// Removing redundant lastchange field and changing with last Retrival:
	delete(data, "lastChange")
	data["lastRetrival"] = time.Now().Format("20060102 15:04")
	
	// Fetching country data:
	countryInterface, currency, err := Countries(isoCode)
	if err != nil {
		http.Error(w, "Error Fetching Country data", http.StatusInternalServerError)
		return
	}
	
	// Convert []interface{} -> []string
	targetCurrencies := utils.ConvSliceInter[string](features["targetCurrencies"].([]interface{}))

	// Fetching Currency exchange rates:
	Currencies, err := Currencies(currency, targetCurrencies)
	if err != nil {
		http.Error(w, "Error fetching Currency exchange rates", http.StatusInternalServerError)
		return
	}

	// Fetching meteo.go:
	metego, err := Weather(utils.ConvSliceInter[float64](countryInterface["latlng"].([]interface{})))
	if err != nil {
		http.Error(w, "Error fetching weather conditions", http.StatusInternalServerError)
		return
	}
	
	// Combining output from weather and country:
	for key, value := range metego {
		countryInterface[key] = value
	}

	// Iterate through features and only keep the features that are true:
	// Change features that are true into the right values.
	for key, value := range features {
		if v, ok := value.(bool); ok && !v {
			delete(features, key)
		} else if ok && v {
			features[key] = countryInterface[key]
		}
	}

	// Setting target currencies:
	features["targetCurrencies"] = Currencies

	// Returning data as JSON:
	result, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error marshaling data to JSON: ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
