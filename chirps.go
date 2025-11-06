package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/skylarhoughtongithub/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error reading request body")
		return
	}

	type createChirpRequest struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	var req createChirpRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := cleanProfaneWords(req.Body)

	params := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: req.UserID,
	}

	dbChirp, err := cfg.DB.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp")
		return
	}

	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirpResponse := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: dbChirp.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirpResponse)
}
