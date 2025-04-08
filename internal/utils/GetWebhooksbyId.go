package utils

import (
	"fmt"
)

// Function to fetch a webhooks "sub collection by their ID"
func GetWebhookByID(id string) (string, error) {
	eventPrefix := id[:3]
	switch eventPrefix {
	case "INV":
		return "INVOKE", nil
	case "REG":
		return "REGISTER", nil
	case "DEL":
		return "DELETE", nil
	case "CHA":
		return "CHANGE", nil
	default:
		return "", fmt.Errorf("Unknown EVENT type")
	}
}
