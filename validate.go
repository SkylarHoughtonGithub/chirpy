package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type ChirpValidationRequest struct {
	Body string `json:"body"`
}

type ChirpValidationResponse struct {
	Valid bool   `json:"valid,omitempty"`
	Error string `json:"error,omitempty"`
}

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var validationReq ChirpValidationRequest
	err = json.Unmarshal(body, &validationReq)
	if err != nil {
		respondWithError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(validationReq.Body) > 140 {
		respondWithError(w, "Chirp is too long", http.StatusBadRequest)
		return
	}

	respondWithSuccess(w)
}

func respondWithError(w http.ResponseWriter, errorMsg string, statusCode int) {
	response := ChirpValidationResponse{
		Error: errorMsg,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}

func respondWithSuccess(w http.ResponseWriter) {
	response := ChirpValidationResponse{
		Valid: true,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
