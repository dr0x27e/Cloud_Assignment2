package utils

import (
	"Assignment2/internal/constants"
	"Assignment2/internal/database"
	"context"
	"fmt"
)

func GetSubCollection(event string, country string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	
	// Firestore client:
	ctx := context.Background()
	client := database.GetFirebaseClient()

	// Fetching the right subcollection:
	collection := client.Collection(constants.Webhooks).Doc(event).Collection(constants.SubCollection)
	docs, err := collection.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("Error Fetching Documents from subcollection")
	}

	for _, doc := range docs {
		data := doc.Data()
		countryVal, _ := data["Country"].(string)
		if countryVal == country || countryVal == "" {
			result = append(result, data)
		}
	}

	return result, nil
}	
