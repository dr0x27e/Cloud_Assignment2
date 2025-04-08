package api

import(
	"Assignment2/internal/constants"
	"encoding/json"
	"net/http"
	"fmt"
	"io"
)


func Weather(input []float64) (map[string]interface{}, error) {
	// Creating URL
	OpenMeteo_API := fmt.Sprintf(
		"%s?latitude=%f&longitude=%f&current=temperature_2m,precipitation",
		constants.OpenMeteo_API,
		input[0],
		input[1])
	
	// Fetching data
	resp, err := http.Get(OpenMeteo_API)
	if err != nil {
		return map[string]interface{}{}, 
		fmt.Errorf("error fetching the url")
	}
	defer resp.Body.Close()
	
	// Reading body:
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{}, 
		fmt.Errorf("failed to read the response body: %w", err) 
	}

	// Parse JSON into map interface
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return map[string]interface{}{}, 
		fmt.Errorf("Failed to Parse json")
	}

	// Fetching current data:
	current := data["current"].(map[string]interface{})

	result  := map[string]interface{}{
		"temperature"  : current["temperature_2m"],
		"precipitation": current["precipitation"],
	}

	return result, nil
}


