package api

import(
	"encoding/json"
	"net/http"
	"io"
)

// Stub handler for test endpoint in python api.
func Python_stub(w http.ResponseWriter, r *http.Request) {
	// Send GET request to Python API:
	resp, err := http.Get("http://localhost:5000/test")
	if err != nil {
		http.Error(w, "Failed to contact Python API: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Python API returned error: "+resp.Status, http.StatusInternalServerError)
		return
	}

	// Read and parse the response:
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read Python API response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var result struct {
		Message string `json:"message"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		http.Error(w, "Failed to parse Python API response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response back to server:
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": result.Message})
}
