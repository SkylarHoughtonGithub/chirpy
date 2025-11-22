package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/skylarhoughtongithub/chirpy/internal/api/response"
	"github.com/skylarhoughtongithub/chirpy/internal/auth"
	"github.com/skylarhoughtongithub/chirpy/internal/service"
)

// WebhookHandler handles webhook HTTP requests
type WebhookHandler struct {
	userService *service.UserService
	polkaKey    string
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(userService *service.UserService, polkaKey string) *WebhookHandler {
	return &WebhookHandler{
		userService: userService,
		polkaKey:    polkaKey,
	}
}

// HandlePolkaWebhook handles POST /api/polka/webhooks
func (h *WebhookHandler) HandlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	// Verify API key
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Missing or malformed API key")
		return
	}

	if apiKey != h.polkaKey {
		response.Error(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	// Parse webhook payload
	var req struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Only handle user.upgraded event
	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Upgrade user to Chirpy Red
	err = h.userService.UpgradeToChirpyRed(r.Context(), req.Data.UserID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "User not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
