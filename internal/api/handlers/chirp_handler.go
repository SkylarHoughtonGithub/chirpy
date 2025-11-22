package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/skylarhoughtongithub/chirpy/internal/api/response"
	"github.com/skylarhoughtongithub/chirpy/internal/auth"
	"github.com/skylarhoughtongithub/chirpy/internal/domain"
	"github.com/skylarhoughtongithub/chirpy/internal/service"
)

// ChirpHandler handles chirp-related HTTP requests
type ChirpHandler struct {
	chirpService *service.ChirpService
	jwtSecret    string
}

// NewChirpHandler creates a new chirp handler
func NewChirpHandler(chirpService *service.ChirpService, jwtSecret string) *ChirpHandler {
	return &ChirpHandler{
		chirpService: chirpService,
		jwtSecret:    jwtSecret,
	}
}

// CreateChirp handles POST /api/chirps
func (h *ChirpHandler) CreateChirp(w http.ResponseWriter, r *http.Request) {
	// Extract and validate JWT
	userID, err := h.authenticateRequest(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid or missing token")
		return
	}

	// Parse request body
	var req struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Call service layer
	chirp, err := h.chirpService.CreateChirp(r.Context(), req.Body, userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// Return response
	response.JSON(w, http.StatusCreated, toChirpResponse(chirp))
}

// GetChirp handles GET /api/chirps/{chirpID}
func (h *ChirpHandler) GetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := h.chirpService.GetChirp(r.Context(), chirpID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, toChirpResponse(chirp))
}

// GetAllChirps handles GET /api/chirps
func (h *ChirpHandler) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	var authorID *uuid.UUID
	if authorIDStr := r.URL.Query().Get("author_id"); authorIDStr != "" {
		id, err := uuid.Parse(authorIDStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
		authorID = &id
	}

	sortDesc := r.URL.Query().Get("sort") == "desc"

	// Call service layer
	chirps, err := h.chirpService.GetAllChirps(r.Context(), authorID, sortDesc)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// Convert to response DTOs
	chirpResponses := make([]chirpResponse, len(chirps))
	for i, chirp := range chirps {
		chirpResponses[i] = toChirpResponse(chirp)
	}

	response.JSON(w, http.StatusOK, chirpResponses)
}

// DeleteChirp handles DELETE /api/chirps/{chirpID}
func (h *ChirpHandler) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	// Authenticate user
	userID, err := h.authenticateRequest(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid or missing token")
		return
	}

	// Parse chirp ID
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	// Call service layer
	err = h.chirpService.DeleteChirp(r.Context(), chirpID, userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper methods

func (h *ChirpHandler) authenticateRequest(r *http.Request) (uuid.UUID, error) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return uuid.Nil, err
	}

	return auth.ValidateJWT(tokenString, h.jwtSecret)
}

func (h *ChirpHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrChirpTooLong:
		response.Error(w, http.StatusBadRequest, "Chirp is too long")
	case domain.ErrChirpNotFound:
		response.Error(w, http.StatusNotFound, "Chirp not found")
	case domain.ErrUnauthorized:
		response.Error(w, http.StatusForbidden, "You are not authorized to delete this chirp")
	default:
		response.Error(w, http.StatusInternalServerError, "Internal server error")
	}
}

// Response DTOs

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func toChirpResponse(chirp *domain.Chirp) chirpResponse {
	return chirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: chirp.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
}
