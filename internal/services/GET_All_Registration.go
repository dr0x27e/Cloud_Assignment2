package services

import (
	"Assignment2/internal/constants"
	"Assignment2/internal/database"
	"encoding/json"
	"net/http"
	"context"
	"log"
)


func GET_All_Registration(w http.ResponseWriter, r* http.Request) {
	log.Println("\nReceived request to display all registers\n")

	// Firestore client:
	ctx := context.Background()
	client := database.GetFirebaseClient()


	// Fetching all focuments inside Registrations collection:
	iter := client.Collection(constants.Registration).Documents(ctx)

	// Array to store documents:
	var documents []map[string]interface{}

	// Iterating through collection:
	for {
		doc, err := iter.Next()
		if err != nil {
			break // No more documents
		}
	
		// Converting thr data to a map and adding ID
		docData := doc.Data()
		docData["id"] = doc.Ref.ID

		// Appening document to array:
		documents = append(documents, docData)
	}

	// Converting array to JSON:
	response, err := json.Marshal(documents)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Response:
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
