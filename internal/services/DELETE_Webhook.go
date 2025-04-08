package services

import (
	"Assignment2/internal/constants"
	"Assignment2/internal/database"
	"Assignment2/internal/utils"
	"net/http"
	"context"
	"log"
)

func DeleteWeb(w http.ResponseWriter, r *http.Request) {
	log.Println("Received " + r.Method + " " + "request.")

	// Extract registration ID from the URL
	webhookID := r.PathValue("id")
	if webhookID == "" {
		http.Error(w, "Missing registration ID", http.StatusBadRequest)
		return
	}

	log.Println("Deleting webhook with specific ID:", webhookID)

	// Firestore deletion
	ctx := context.Background()
	client := database.GetFirebaseClient()

	// Fetching the right sub collection:
	subCollection, err := utils.GetWebhookByID(webhookID)
	if err != nil {
		log.Println("Unnown event type")
		http.Error(w, "Failed to fetch subCollection", http.StatusInternalServerError)
		return
	}
	
	// Deleting document:
	collection := client.Collection(constants.Webhooks).Doc(subCollection).Collection(constants.SubCollection)
	_, err = collection.Doc(webhookID).Delete(ctx)
	if err != nil {
		log.Println("Error deleting document:", err)
		http.Error(w, "Failed to delete webhook", http.StatusInternalServerError)
		return
	}

	log.Println("Successfully deleted webhook:", webhookID)
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
