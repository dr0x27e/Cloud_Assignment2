package services

import (
	"Assignment2/internal/constants"
	"Assignment2/internal/database"
	"Assignment2/internal/utils"
	"encoding/json"
	"net/http"
	"context"
)

func GET_Id_Registration(w http.ResponseWriter, r *http.Request) {
	// Extract registration ID from the URL
	registerID := r.PathValue("id")
	if registerID == "" {
		http.Error(w, "Missing registration ID", http.StatusBadRequest)
		return
	}
	
	// Firestore client:
	ctx := context.Background()
	client := database.GetFirebaseClient()
	
	// Fetching registration:
	res := client.Collection(constants.Registration).Doc(registerID)
	doc, err := res.Get(ctx)
	if err != nil {
		http.Error(w, "Error getting document", http.StatusInternalServerError)
		return
	}
	
	// Parsing data:
	data := doc.Data()
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error marshalling payload", http.StatusInternalServerError)
		return
	}
	
	// Response:
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	// Invoke webhooks:
	utils.Invoke("INVOKE", data["isoCode"].(string), "ID: " +
	registerID + " Got invoked: GET_ID")
}
