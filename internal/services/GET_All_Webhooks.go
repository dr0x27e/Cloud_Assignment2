package services

import (
	"Assignment2/internal/constants"
	"Assignment2/internal/database"
	"encoding/json"
	"net/http"
	"context"
	"log"
)

func GetAllWeb(w http.ResponseWriter, r *http.Request) {
	log.Println("\nReceived request to display all webhooks\n")

	// Firestore client:
	ctx := context.Background()
	client := database.GetFirebaseClient()
	
	// Fetching all webhooks:
	var allWebhooks []map[string]interface{}
	for _, event := range constants.AllEvents {
		collection := client.Collection(constants.Webhooks).Doc(event).Collection(constants.SubCollection)

		docs, err := collection.Documents(ctx).GetAll()
		if err != nil {
			log.Printf("Failed to get documents for event %s: %v", event, err)
			continue // skip this event but continue with the others
		}

		for _, doc := range docs {
			data := doc.Data()
			data["id"] = doc.Ref.ID     // Add the document ID to the data
			allWebhooks = append(allWebhooks, data)
		}
	}

	// Marshal result to JSON
	jsonResponse, err := json.Marshal(allWebhooks)
	if err != nil {
		http.Error(w, "Failed to marshal webhooks", http.StatusInternalServerError)
		return
	}

	// Response:
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
