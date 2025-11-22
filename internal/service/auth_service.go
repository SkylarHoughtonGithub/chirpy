package service

import (
	"context"
	"fmt"
	"time"

	"github.com/skylarhoughtongithub/chirpy/internal/auth"
	"github.com/skylarhoughtongithub/chirpy/internal/domain"
	"github.com/skylarhoughtongithub/chirpy/internal/repository"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.RefreshTokenRepository
	jwtSecret string
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.RefreshTokenRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtSecret: jwtSecret,
	}
}

// LoginResponse contains tokens and user info
type LoginResponse struct {
	User         *domain.User
	AccessToken  string
	RefreshToken string
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Check password
	match, err := auth.CheckPasswordHash(password, user.HashedPassword)
	if err != nil || !match {
		return nil, domain.ErrInvalidCredentials
	}

	// Create access token (1 hour)
	accessToken, err := auth.MakeJWT(user.ID, s.jwtSecret, time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Create refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	// Save refresh token (60 days)
	err = s.tokenRepo.Create(ctx, refreshToken, user.ID, time.Now().Add(60*24*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshAccessToken creates a new access token from a refresh token
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	// Get user from refresh token
	user, err := s.tokenRepo.GetUserByToken(ctx, refreshToken)
	if err != nil {
		return "", domain.ErrInvalidToken
	}

	// Create new access token
	accessToken, err := auth.MakeJWT(user.ID, s.jwtSecret, time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to create access token: %w", err)
	}

	return accessToken, nil
}

// RevokeRefreshToken revokes a refresh token
func (s *AuthService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	err := s.tokenRepo.Revoke(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}
