package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/skylarhoughtongithub/chirpy/internal/domain"
	"github.com/skylarhoughtongithub/chirpy/internal/repository"
)

// ChirpService handles business logic for chirps
type ChirpService struct {
	chirpRepo repository.ChirpRepository
}

// NewChirpService creates a new chirp service
func NewChirpService(chirpRepo repository.ChirpRepository) *ChirpService {
	return &ChirpService{
		chirpRepo: chirpRepo,
	}
}

// CreateChirp creates a new chirp with validation and cleaning
func (s *ChirpService) CreateChirp(ctx context.Context, body string, userID uuid.UUID) (*domain.Chirp, error) {
	// Business logic validation
	if len(body) > 140 {
		return nil, domain.ErrChirpTooLong
	}

	// Clean profane words (business logic)
	cleanedBody := cleanProfaneWords(body)

	// Delegate to repository
	chirp, err := s.chirpRepo.Create(ctx, cleanedBody, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create chirp: %w", err)
	}

	return chirp, nil
}

// GetChirp retrieves a single chirp by ID
func (s *ChirpService) GetChirp(ctx context.Context, chirpID uuid.UUID) (*domain.Chirp, error) {
	chirp, err := s.chirpRepo.GetByID(ctx, chirpID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chirp: %w", err)
	}
	return chirp, nil
}

// GetAllChirps retrieves all chirps with optional filtering
func (s *ChirpService) GetAllChirps(ctx context.Context, authorID *uuid.UUID, sortDesc bool) ([]*domain.Chirp, error) {
	var chirps []*domain.Chirp
	var err error

	if authorID != nil {
		chirps, err = s.chirpRepo.GetByAuthor(ctx, *authorID)
	} else {
		chirps, err = s.chirpRepo.GetAll(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get chirps: %w", err)
	}

	// Sort in service layer (business logic)
	if sortDesc {
		// Reverse the slice
		for i := len(chirps)/2 - 1; i >= 0; i-- {
			opp := len(chirps) - 1 - i
			chirps[i], chirps[opp] = chirps[opp], chirps[i]
		}
	}

	return chirps, nil
}

// DeleteChirp deletes a chirp if the user owns it
func (s *ChirpService) DeleteChirp(ctx context.Context, chirpID, userID uuid.UUID) error {
	// Get chirp to verify ownership
	chirp, err := s.chirpRepo.GetByID(ctx, chirpID)
	if err != nil {
		return fmt.Errorf("failed to get chirp: %w", err)
	}

	// Business logic: verify ownership
	if chirp.UserID != userID {
		return domain.ErrUnauthorized
	}

	err = s.chirpRepo.Delete(ctx, chirpID)
	if err != nil {
		return fmt.Errorf("failed to delete chirp: %w", err)
	}

	return nil
}

// cleanProfaneWords is internal business logic
func cleanProfaneWords(text string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(text, " ")
	for i, word := range words {
		lowercaseWord := strings.ToLower(word)
		for _, profane := range profaneWords {
			if lowercaseWord == profane {
				words[i] = "****"
				break
			}
		}
	}

	return strings.Join(words, " ")
}
