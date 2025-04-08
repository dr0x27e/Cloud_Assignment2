package api

import (
	"Assignment2/internal/utils"
	"mime/multipart"
	"net/http"
	"bytes"
	"io"
)


func Prediction_Python(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form data (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form data (To Big): "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form (Empty?): "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file into a buffer
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the multipart form for the Python API
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", "image.jpg")
	if err != nil {
		http.Error(w, "Failed to create form file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = part.Write(fileBytes)
	if err != nil {
		http.Error(w, "Failed to write file to form: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Close()

	// Send the request to the PYTHON API:
	req, err := http.NewRequest("POST", "http://10.212.170.29:5000/predict", &buf)
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to contact Python API: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Python API returned error: "+resp.Status, http.StatusInternalServerError)
		return
	}

	// Read and forward the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read Python API response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Invoke webhooks:
	utils.Invoke("PREDICT", "", "Image Predicted.")

	// Forward the Python API's JSON response directly to the client
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
