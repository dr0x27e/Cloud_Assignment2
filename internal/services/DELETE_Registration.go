package services

import (
	"Assignment2/internal/constants"
	"Assignment2/internal/database"
	"Assignment2/internal/utils"
	"net/http"
	"context"
	"fmt"
	"log"
)

func DELETE_Registration(w http.ResponseWriter, r *http.Request) {
	log.Println("Received " + r.Method + " " + "request.")

	// Extract registration ID from the URL
	registerID := r.PathValue("id")
	if registerID == "" {
		http.Error(w, "Missing registration ID", http.StatusBadRequest)
		return
	}

	// Firestore deletion
	ctx := context.Background()
	client := database.GetFirebaseClient()

	// Deleting document And fetching IsoCode:
	collection := client.Collection(constants.Registration)
	doc, err := collection.Doc(registerID).Get(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch ISOCODE aborting delete", http.StatusInternalServerError)
		return
	}
	isoCode := doc.Data()["isoCode"].(string)
	_, err = collection.Doc(registerID).Delete(ctx)
	if err != nil {
		http.Error(w, "Failed to delete registration", http.StatusInternalServerError)
		return
	}
	fmt.Println("Successfully deleted registration: " + registerID)
	
	// Invoke webhooks:
	utils.Invoke(constants.EventDelete, isoCode,"ID: "+
	registerID+" Deleted.")
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
