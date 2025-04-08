package utils

import (
	"encoding/hex"
	"crypto/rand"
	"fmt"
)

func GenerateRandomID(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random ID: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}
