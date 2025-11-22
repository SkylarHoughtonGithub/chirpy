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

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *service.UserService
	jwtSecret   string
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

// Create handles POST /api/users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.Email, req.Password)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, toUserResponse(user))
}

// Update handles PUT /api/users
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Authenticate user
	userID, err := h.authenticateRequest(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid or missing token")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), userID, req.Email, req.Password)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, toUserResponse(user))
}

// Helper methods

func (h *UserHandler) authenticateRequest(r *http.Request) (uuid.UUID, error) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return uuid.Nil, err
	}

	return auth.ValidateJWT(tokenString, h.jwtSecret)
}

func (h *UserHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrInvalidInput:
		response.Error(w, http.StatusBadRequest, "Email and password are required")
	case domain.ErrUserAlreadyExists:
		response.Error(w, http.StatusConflict, "User already exists")
	case domain.ErrUserNotFound:
		response.Error(w, http.StatusNotFound, "User not found")
	default:
		response.Error(w, http.StatusInternalServerError, "Internal server error")
	}
}

// Response DTOs

type userResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func toUserResponse(user *domain.User) userResponse {
	return userResponse{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
}
