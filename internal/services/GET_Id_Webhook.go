package services

import (
	"Assignment2/internal/constants"
	"Assignment2/internal/database"
	"Assignment2/internal/utils"
	"encoding/json"
	"net/http"
	"context"
	"log"
)

func GET_Id_Webhook(w http.ResponseWriter, r *http.Request) {
	log.Println("Received " + r.Method + " " + "request.")

	// Extract webhook ID from the URL
	webhookID := r.PathValue("id")
	if webhookID == "" {
		http.Error(w, "Missing registration ID", http.StatusBadRequest)
		return
	}
	
	// Firestore client:
	ctx := context.Background()
	client := database.GetFirebaseClient()

	// Fetching the right sub collection:
	subCollection, err := utils.GetWebhookByID(webhookID)
	if err != nil {
		log.Println("Unnown event type")
		http.Error(w, "Failed to fetch subCollection", http.StatusInternalServerError)
		return
	}

	// Fetching the sub collection:
	collection := client.Collection(constants.Webhooks).Doc(subCollection).Collection(constants.SubCollection)
	doc, err := collection.Doc(webhookID).Get(ctx)
	if err != nil {
		log.Println("Error getting document:", err.Error())
		http.Error(w, "Error getting document", http.StatusInternalServerError)
		return
	}

	// Parsing the data:
	data := doc.Data()
	data["id"] = doc.Ref.ID
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshalling payload:", err.Error())
		http.Error(w, "Error marshalling payload", http.StatusInternalServerError)
		return
	}

	// Set response headers and send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
