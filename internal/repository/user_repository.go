package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/skylarhoughtongithub/chirpy/internal/database"
	"github.com/skylarhoughtongithub/chirpy/internal/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, email, hashedPassword string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, id uuid.UUID, email, hashedPassword string) (*domain.User, error)
	UpgradeToChirpyRed(ctx context.Context, id uuid.UUID) error
	DeleteAll(ctx context.Context) error
}

// userRepository is the concrete implementation
type userRepository struct {
	db *database.Queries
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.Queries) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user
func (r *userRepository) Create(ctx context.Context, email, hashedPassword string) (*domain.User, error) {
	dbUser, err := r.db.CreateUser(ctx, database.CreateUserParams{
		Email:          email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return nil, err
	}

	return toDomainUser(dbUser), nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	dbUser, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return toDomainUser(dbUser), nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	dbUser, err := r.db.GetUserByEmail(ctx, "") // Note: You may need to add GetUserByID to SQLC
	if err != nil {
		return nil, err
	}

	return toDomainUser(dbUser), nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, id uuid.UUID, email, hashedPassword string) (*domain.User, error) {
	dbUser, err := r.db.UpdateUser(ctx, database.UpdateUserParams{
		ID:             id,
		Email:          email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return nil, err
	}

	return toDomainUser(dbUser), nil
}

// UpgradeToChirpyRed upgrades a user to Chirpy Red
func (r *userRepository) UpgradeToChirpyRed(ctx context.Context, id uuid.UUID) error {
	return r.db.UpgradeUserToChirpyRed(ctx, id)
}

// DeleteAll deletes all users (for testing/reset)
func (r *userRepository) DeleteAll(ctx context.Context) error {
	return r.db.DeleteAllUsers(ctx)
}

// toDomainUser converts database model to domain model
func toDomainUser(dbUser database.User) *domain.User {
	return &domain.User{
		ID:             dbUser.ID,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
		Email:          dbUser.Email,
		HashedPassword: dbUser.HashedPassword,
		IsChirpyRed:    dbUser.IsChirpyRed,
	}
}
