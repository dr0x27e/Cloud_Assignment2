package webhooks

import (
	"Assignment2/internal/constants"
	"Assignment2/internal/utils"
	"fmt" // TEMP
	"log"
)


func Service(event string, country string) {
	log.Println("Recived POST request")

	// Validating EVENT:
	if !constants.ValidEvents[event] {
		log.Println("Invalid Event")
		return
	}

	// Fetching sub collection:
	webhooks, err := utils.GetSubCollection(event, country)
	if err != nil {
		log.Println("Error fetchin subcollection")
		return
	}
		
	for _, hook := range webhooks {
		fmt.Println(hook)
	}
}
