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
)

func WebhookReg(w http.ResponseWriter, r *http.Request) {
	// Firestore client:
	ctx := context.Background()
	client := database.GetFirebaseClient()

	var webhook structs.WebhookRegistrationModel

	// Decode request body:
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields:
	
	// Validating EVENT:
	if !constants.ValidEvents[webhook.Event] {
		http.Error(w, "Not a valid Event", http.StatusBadRequest)
		return
	}

	// Validating URL:
	if webhook.Country != ""  {
		// Fetching data to see if its a real country:
		if (len(webhook.Country) > 2 || len(webhook.Country) < 2) {
			http.Error(w, "Invalid IsoCode", http.StatusBadRequest)
			return
		}
		
		url := constants.Country_API + webhook.Country + "?fields=name"

		// Fetching from Currency_API API:
        resp, err := http.Get(url)
        if err != nil {
			http.Error(w, "Error fetching from Country_API", http.StatusInternalServerError)
			return
        }
        defer resp.Body.Close()

        // Checking response status:
        if resp.StatusCode != http.StatusOK {
			http.Error(w, "Invalid IsoCode", http.StatusBadRequest)
			return
        }

		// Taking the IsoCode to upper for consistency:
		webhook.Country = strings.ToUpper(webhook.Country)
	}
	
	// Creating a custom ID for the webhook, with its event sub collection initials
	// as the 3 leading letters (for faster and more efficient fetching later):
	randomID, err := utils.GenerateRandomID(10) // 20 hex characters long.
	if err != nil {
		http.Error(w, "Failed to generate random ID", http.StatusInternalServerError)
		return
	}
	customID := webhook.Event[:3] + randomID

	// Add document to Firestore subcollection based on EVENT type:
	collection := client.Collection(constants.Webhooks).Doc(webhook.Event).Collection(constants.SubCollection)
	_, err = collection.Doc(customID).Set(ctx, webhook)
	if err != nil {
		http.Error(w, "Failed to add webhook", http.StatusInternalServerError)
		return
	}

	// Respond with generated ID
	response := map[string]string{"id": customID}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}
