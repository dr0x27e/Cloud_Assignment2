package services

import (
	"Assignment2/internal/constants"
	"cloud.google.com/go/firestore"
	"Assignment2/internal/database"
	"Assignment2/internal/utils"
	"encoding/json"
	"net/http"
	"context"
	"time"
	"io"
)

func PATCH_Registration(w http.ResponseWriter, r *http.Request) {
	// Reading body:
	content, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Reading Payload Failed", http.StatusInternalServerError)
		return
	}

	// Check if the body is empty:
	if len(content) == 0 {
		http.Error(w, "Payload is empty", http.StatusBadRequest)
		return
	}

	// Parsing the JSON body into a map for partial update:
	var updates map[string]interface{}
	if err := json.Unmarshal(content, &updates); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Extracting the document ID:
	docID := r.PathValue("id")
	if docID == "" {
		http.Error(w, "Missing document ID", http.StatusBadRequest)
		return
	}

	// Firestore client:
	ctx := context.Background()
	client := database.GetFirebaseClient()

	// Referencing firebase document:
	docRef := client.Collection(constants.Registration).Doc(docID)
	
	// Handeling Isocode update if provided:
	if isoRaw, ok := updates["isoCode"]; ok {
		if isoCode, ok := isoRaw.(string); ok && isoCode != "" {
			// Fetching country name from API:
			url := constants.Country_API + isoCode + "?fields=name"
			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, "Failed to fetch the country data", http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				http.Error(w, "Error fetching from API (Wrong iso?)", http.StatusInternalServerError)
				return
			}

			// Reading body:
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, "Could not read API reponse body", http.StatusInternalServerError)
				return
			}

			// Parsing JSON:
			var countryData map[string]interface{}
			if err := json.Unmarshal(body, &countryData); err != nil {
				http.Error(w, "Error Decoding country data", http.StatusInternalServerError)
				return
			}

			Name := countryData["name"].(map[string]interface{})
			Common := Name["common"].(string)

			if Common == "" {
				http.Error(w, "Invalid ISOcode", http.StatusBadRequest)
				return
			}

			// Updating both isoCode and Country name:
			updates["isoCode"] = isoCode
			updates["country"] = Common
		}
	}

	// Updating lastChange time:
	updates["lastChange"] = time.Now().Format("20060102 15:04")

	// Prepare update fields for Registration:
	var updateFields []firestore.Update
	for key, value := range updates {
		if key == "features" {
			// Handle features map
			if features, ok := value.(map[string]interface{}); ok {
				for fKey, fValue := range features {
					updateFields = append(updateFields, firestore.Update{
						Path:  "features." + fKey,
						Value: fValue,
					})
				}
			}
		} else {
			updateFields = append(updateFields, firestore.Update{
				Path:  key,
				Value: value,
			})
		}
	}

	// Applying update:
	_, err = docRef.Update(ctx, updateFields)
	if err != nil {
		http.Error(w, "Failed to update document", http.StatusInternalServerError)
		return
	}

	// Fetching the isoCode:
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch updated document", http.StatusInternalServerError)
		return
	}
	data := docSnap.Data()

	// Invoke webhooks:
	utils.Invoke("CHANGE", data["isoCode"].(string), "ID: " + docRef.ID + 
	" Got changed with method: PATCH")

	// Response:
	w.WriteHeader(http.StatusCreated)
}
