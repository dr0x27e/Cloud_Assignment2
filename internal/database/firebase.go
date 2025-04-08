package database

import (
	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"Assignment2/internal/structs"
	"firebase.google.com/go"
	"context"
	"log"
)


var firebaseClient *structs.FirebaseClient

// Initializing the Firebase client:
func InitFirebaseClient() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("./internal/database/privateKey.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v", err)
	}

	// Assign client to global variable:
	firebaseClient = &structs.FirebaseClient{Client: client}
}

// Get the Firestore client:
func GetFirebaseClient() *firestore.Client {
	return firebaseClient.Client
}

// Close the Firestore client:
func CloseFirebaseClient() {
	if firebaseClient != nil && firebaseClient.Client != nil {
		if err := firebaseClient.Client.Close(); err != nil {
			log.Fatalf("Error closing Firestore client: %v", err)
		}
	}
}
