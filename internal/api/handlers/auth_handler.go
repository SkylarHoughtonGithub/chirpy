package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/skylarhoughtongithub/chirpy/internal/api/response"
	"github.com/skylarhoughtongithub/chirpy/internal/auth"
	"github.com/skylarhoughtongithub/chirpy/internal/domain"
	"github.com/skylarhoughtongithub/chirpy/internal/service"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles POST /api/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	loginResp, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	resp := struct {
		ID           string `json:"id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		Email        string `json:"email"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		ID:           loginResp.User.ID.String(),
		CreatedAt:    loginResp.User.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    loginResp.User.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Email:        loginResp.User.Email,
		IsChirpyRed:  loginResp.User.IsChirpyRed,
		Token:        loginResp.AccessToken,
		RefreshToken: loginResp.RefreshToken,
	}

	response.JSON(w, http.StatusOK, resp)
}

// Refresh handles POST /api/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Malformed token")
		return
	}

	accessToken, err := h.authService.RefreshAccessToken(r.Context(), refreshToken)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: accessToken,
	})
}

// Revoke handles POST /api/revoke
func (h *AuthHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Malformed token")
		return
	}

	err = h.authService.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper methods

func (h *AuthHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrInvalidCredentials:
		response.Error(w, http.StatusUnauthorized, "Incorrect email or password")
	case domain.ErrInvalidToken:
		response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
	case domain.ErrUserNotFound:
		response.Error(w, http.StatusUnauthorized, "User not found")
	default:
		response.Error(w, http.StatusInternalServerError, "Internal server error")
	}
}
