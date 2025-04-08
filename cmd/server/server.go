package main

import (
	"Assignment2/internal/database"
	"Assignment2/internal/webhooks"
	"Assignment2/internal/api"
	"net/http"
	"log"
	"os"
)

func main() {

	// INITIALIZING FIREBASE CLIENT:
	database.InitFirebaseClient()
	// Closing the firebase client at the end of the service:
	defer database.CloseFirebaseClient()

	// Getting the port
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Instantiate the router
	router := http.NewServeMux()

	// Handlers
	router.HandleFunc("/dashboard/v1/dashboards/{id}", api.Dashboards)
	
	router.HandleFunc("/dashboard/v1/registrations/", api.Registrations)
	router.HandleFunc("/dashboard/v1/registrations/{id}", api.Registrations)

	router.HandleFunc("/dashboard/v1/notifications/", webhooks.Notifications)
	router.HandleFunc("/dashboard/v1/notifications/{id}", webhooks.Notifications)
	
	router.HandleFunc("/dashboard/v1/status/", api.Status)
	router.HandleFunc("/dashboard/v1/", api.Empty)

	// Test endpoint to call Python API
	// router.HandleFunc("/api/test", api.Python_stub)

	// Python API Prediction endpoint:
	router.HandleFunc("/api/predict/", api.Prediction_Python)

	// Start HTTP server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":" + port, router))
}
