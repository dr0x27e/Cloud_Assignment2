package utils

import (
	"net/http"
	"strconv"
	"bytes"
	"fmt"
	"log"
	"io"
)

// Webhook invocation function:
func Invoke(event string, country string , payload string) {
	log.Println("Attempting invocation...")
	fmt.Println("Event   : " + event)
	fmt.Println("Country : " + country)
	fmt.Println("Payload : " + payload)

	collection, err := GetSubCollection(event, country)
	if err != nil {
		log.Println("Failed to fetch subcollection: " + err.Error())
		return
	}

	for _, doc := range collection {
		go invoke_post(doc, payload)
	}
}


// Function to send requests:
func invoke_post(document map[string]interface{}, payload string) {
	url := document["Url"].(string)
	event := document["Event"].(string)

	// Creating a new POST request:
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte("Event "+
	event+" occured. Payload: "+payload)))
	if err != nil {
		log.Println("Error during request creation: " + err.Error())
		return
	}

	// Perform invocation:
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error in HTTP request: " + err.Error())
		return
	}

	// Reading resposne:
	response, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Somethign went wrong with invocation response: " + err.Error())
		return
	}

	log.Println("Webhook " + url + " Invoked. Recieved status code: " + 
		strconv.Itoa(res.StatusCode) + " and body : " + string(response))
}
