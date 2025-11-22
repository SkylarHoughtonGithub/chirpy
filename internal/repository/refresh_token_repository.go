package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/skylarhoughtongithub/chirpy/internal/database"
	"github.com/skylarhoughtongithub/chirpy/internal/domain"
)

// RefreshTokenRepository defines the interface for refresh token data access
type RefreshTokenRepository interface {
	Create(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error
	GetUserByToken(ctx context.Context, token string) (*domain.User, error)
	Revoke(ctx context.Context, token string) error
}

// refreshTokenRepository is the concrete implementation
type refreshTokenRepository struct {
	db *database.Queries
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *database.Queries) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create inserts a new refresh token
func (r *refreshTokenRepository) Create(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error {
	return r.db.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})
}

// GetUserByToken retrieves a user by their refresh token
func (r *refreshTokenRepository) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	dbUser, err := r.db.GetUserFromRefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return toDomainUser(dbUser), nil
}

// Revoke revokes a refresh token
func (r *refreshTokenRepository) Revoke(ctx context.Context, token string) error {
	return r.db.RevokeRefreshToken(ctx, token)
}
