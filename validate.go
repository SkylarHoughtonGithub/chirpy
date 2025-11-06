package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type ChirpValidationRequest struct {
	Body string `json:"body"`
}

type ChirpValidationResponse struct {
	CleanedBody string `json:"cleaned_body,omitempty"`
	Error       string `json:"error,omitempty"`
}

// Helper function to respond with an error
func respondWithError(w http.ResponseWriter, code int, msg string) {
	response := ChirpValidationResponse{
		Error: msg,
	}

	// Use json.Marshal to ensure clean JSON output
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}

// Helper function to respond with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Use json.Marshal to ensure clean JSON output
	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}

// Function to clean profane words
func cleanProfaneWords(text string) string {
	// List of profane words to replace
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	// Split the text into words
	words := strings.Split(text, " ")

	// Replace profane words
	for i, word := range words {
		// Convert to lowercase for case-insensitive matching
		// Preserve original casing of surrounding words
		lowercaseWord := strings.ToLower(word)

		for _, profane := range profaneWords {
			if lowercaseWord == profane {
				words[i] = "****"
				break
			}
		}
	}

	// Rejoin the words
	return strings.Join(words, " ")
}

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	// Ensure this is a POST request
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error reading request body")
		return
	}

	// Parse the JSON request
	var validationReq ChirpValidationRequest
	err = json.Unmarshal(body, &validationReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate chirp length
	if len(validationReq.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// Clean the chirp body
	cleanedBody := cleanProfaneWords(validationReq.Body)

	// Prepare and send the response
	response := ChirpValidationResponse{
		CleanedBody: cleanedBody,
	}
	respondWithJSON(w, http.StatusOK, response)
}
