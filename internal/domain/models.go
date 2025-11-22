package domain

import (
	"time"

	"github.com/google/uuid"
)

// Chirp represents a chirp in the domain
type Chirp struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
	UserID    uuid.UUID
}

// User represents a user in the domain
type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Email          string
	HashedPassword string
	IsChirpyRed    bool
}

// RefreshToken represents a refresh token
type RefreshToken struct {
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt *time.Time
}
