package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/skylarhoughtongithub/chirpy/internal/database"
	"github.com/skylarhoughtongithub/chirpy/internal/domain"
)

// ChirpRepository defines the interface for chirp data access
type ChirpRepository interface {
	Create(ctx context.Context, body string, userID uuid.UUID) (*domain.Chirp, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Chirp, error)
	GetAll(ctx context.Context) ([]*domain.Chirp, error)
	GetByAuthor(ctx context.Context, authorID uuid.UUID) ([]*domain.Chirp, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// chirpRepository is the concrete implementation
type chirpRepository struct {
	db *database.Queries
}

// NewChirpRepository creates a new chirp repository
func NewChirpRepository(db *database.Queries) ChirpRepository {
	return &chirpRepository{db: db}
}

// Create inserts a new chirp
func (r *chirpRepository) Create(ctx context.Context, body string, userID uuid.UUID) (*domain.Chirp, error) {
	dbChirp, err := r.db.CreateChirp(ctx, database.CreateChirpParams{
		Body:   body,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return toDomainChirp(dbChirp), nil
}

// GetByID retrieves a chirp by ID
func (r *chirpRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Chirp, error) {
	dbChirp, err := r.db.GetChirp(ctx, id)
	if err != nil {
		return nil, err
	}

	return toDomainChirp(dbChirp), nil
}

// GetAll retrieves all chirps
func (r *chirpRepository) GetAll(ctx context.Context) ([]*domain.Chirp, error) {
	dbChirps, err := r.db.GetAllChirps(ctx)
	if err != nil {
		return nil, err
	}

	chirps := make([]*domain.Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		chirps[i] = toDomainChirp(dbChirp)
	}

	return chirps, nil
}

// GetByAuthor retrieves chirps by author
func (r *chirpRepository) GetByAuthor(ctx context.Context, authorID uuid.UUID) ([]*domain.Chirp, error) {
	dbChirps, err := r.db.GetChirpsByAuthor(ctx, authorID)
	if err != nil {
		return nil, err
	}

	chirps := make([]*domain.Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		chirps[i] = toDomainChirp(dbChirp)
	}

	return chirps, nil
}

// Delete removes a chirp
func (r *chirpRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.DeleteChirp(ctx, id)
}

// toDomainChirp converts database model to domain model
func toDomainChirp(dbChirp database.Chirp) *domain.Chirp {
	return &domain.Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
}
