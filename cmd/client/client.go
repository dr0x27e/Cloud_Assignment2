package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"io"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: client <port> <endpoint>")
		return
	}

	port := os.Args[1]
	fmt.Println("Client Listening at Port: ", port)
	endpoint := os.Args[2]
	fmt.Println("At location: http://localhost:" + port + "/client/" + endpoint)

	// Define the handler for POST requests
	http.HandleFunc("/client/"+endpoint, func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a POST request
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method. Only POST is allowed.", http.StatusMethodNotAllowed)
			return
		}

		// Read the payload
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading body:", err)
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		// Log the received string
		log.Printf("Received request at /client/%s with body: %s", endpoint, string(body))

		// Respond with a status and message
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start the HTTP server
	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Client listening on %s at /client/%s\n", addr, endpoint)
	log.Fatal(http.ListenAndServe(addr, nil))
}
