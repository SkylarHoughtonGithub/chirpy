package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/skylarhoughtongithub/chirpy/internal/auth"
	"github.com/skylarhoughtongithub/chirpy/internal/domain"
	"github.com/skylarhoughtongithub/chirpy/internal/repository"
)

// UserService handles user business logic
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with hashed password
func (s *UserService) CreateUser(ctx context.Context, email, password string) (*domain.User, error) {
	// Validate input
	if email == "" || password == "" {
		return nil, domain.ErrInvalidInput
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := s.userRepo.Create(ctx, email, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// UpdateUser updates a user's email and password
func (s *UserService) UpdateUser(ctx context.Context, userID uuid.UUID, email, password string) (*domain.User, error) {
	// Validate input
	if email == "" || password == "" {
		return nil, domain.ErrInvalidInput
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user
	user, err := s.userRepo.Update(ctx, userID, email, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpgradeToChirpyRed upgrades a user to premium status
func (s *UserService) UpgradeToChirpyRed(ctx context.Context, userID uuid.UUID) error {
	err := s.userRepo.UpgradeToChirpyRed(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to upgrade user: %w", err)
	}
	return nil
}
