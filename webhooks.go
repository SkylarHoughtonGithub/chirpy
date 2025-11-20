package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error reading request body")
		return
	}

	type webhookRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		}
	}

	var req webhookRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.DB.UpgradeUserToChirpyRed(r.Context(), req.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
