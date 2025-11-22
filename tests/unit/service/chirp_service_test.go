package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/skylarhoughtongithub/chirpy/internal/domain"
	"github.com/skylarhoughtongithub/chirpy/internal/service"
)

// Mock repository for testing
type mockChirpRepository struct {
	createFunc      func(ctx context.Context, body string, userID uuid.UUID) (*domain.Chirp, error)
	getByIDFunc     func(ctx context.Context, id uuid.UUID) (*domain.Chirp, error)
	getAllFunc      func(ctx context.Context) ([]*domain.Chirp, error)
	getByAuthorFunc func(ctx context.Context, authorID uuid.UUID) ([]*domain.Chirp, error)
	deleteFunc      func(ctx context.Context, id uuid.UUID) error
}

func (m *mockChirpRepository) Create(ctx context.Context, body string, userID uuid.UUID) (*domain.Chirp, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, body, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockChirpRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Chirp, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockChirpRepository) GetAll(ctx context.Context) ([]*domain.Chirp, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *mockChirpRepository) GetByAuthor(ctx context.Context, authorID uuid.UUID) ([]*domain.Chirp, error) {
	if m.getByAuthorFunc != nil {
		return m.getByAuthorFunc(ctx, authorID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockChirpRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func TestChirpService_CreateChirp(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		userID      uuid.UUID
		mockSetup   func(*mockChirpRepository)
		wantErr     error
		wantBody    string
		description string
	}{
		{
			name:   "successful chirp creation",
			body:   "Hello world!",
			userID: uuid.New(),
			mockSetup: func(m *mockChirpRepository) {
				m.createFunc = func(ctx context.Context, body string, userID uuid.UUID) (*domain.Chirp, error) {
					return &domain.Chirp{
						ID:     uuid.New(),
						Body:   body,
						UserID: userID,
					}, nil
				}
			},
			wantErr:     nil,
			wantBody:    "Hello world!",
			description: "Should create chirp successfully",
		},
		{
			name:   "chirp too long",
			body:   "This is a very long chirp that exceeds the maximum allowed length of 140 characters. It should fail validation and return an error. Adding more text to ensure it's over the limit.",
			userID: uuid.New(),
			mockSetup: func(m *mockChirpRepository) {
				// Should not be called
			},
			wantErr:     domain.ErrChirpTooLong,
			description: "Should reject chirps over 140 characters",
		},
		{
			name:   "profanity filtering",
			body:   "This is a kerfuffle situation",
			userID: uuid.New(),
			mockSetup: func(m *mockChirpRepository) {
				m.createFunc = func(ctx context.Context, body string, userID uuid.UUID) (*domain.Chirp, error) {
					return &domain.Chirp{
						ID:     uuid.New(),
						Body:   body,
						UserID: userID,
					}, nil
				}
			},
			wantErr:     nil,
			wantBody:    "This is a **** situation",
			description: "Should filter profane words",
		},
		{
			name:   "multiple profane words",
			body:   "The kerfuffle and sharbert were fornax",
			userID: uuid.New(),
			mockSetup: func(m *mockChirpRepository) {
				m.createFunc = func(ctx context.Context, body string, userID uuid.UUID) (*domain.Chirp, error) {
					return &domain.Chirp{
						ID:     uuid.New(),
						Body:   body,
						UserID: userID,
					}, nil
				}
			},
			wantErr:     nil,
			wantBody:    "The **** and **** were ****",
			description: "Should filter multiple profane words",
		},
		{
			name:   "case insensitive profanity",
			body:   "KERFUFFLE and Sharbert",
			userID: uuid.New(),
			mockSetup: func(m *mockChirpRepository) {
				m.createFunc = func(ctx context.Context, body string, userID uuid.UUID) (*domain.Chirp, error) {
					return &domain.Chirp{
						ID:     uuid.New(),
						Body:   body,
						UserID: userID,
					}, nil
				}
			},
			wantErr:     nil,
			wantBody:    "**** and ****",
			description: "Should filter profanity case-insensitively",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &mockChirpRepository{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}
			service := service.NewChirpService(mockRepo)

			// Execute
			chirp, err := service.CreateChirp(context.Background(), tt.body, tt.userID)

			// Assert
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("%s: expected error %v, got %v", tt.description, tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("%s: unexpected error: %v", tt.description, err)
				return
			}

			if chirp.Body != tt.wantBody {
				t.Errorf("%s: expected body %q, got %q", tt.description, tt.wantBody, chirp.Body)
			}

			if chirp.UserID != tt.userID {
				t.Errorf("%s: expected userID %v, got %v", tt.description, tt.userID, chirp.UserID)
			}
		})
	}
}

func TestChirpService_DeleteChirp(t *testing.T) {
	chirpID := uuid.New()
	ownerID := uuid.New()
	otherUserID := uuid.New()

	tests := []struct {
		name        string
		chirpID     uuid.UUID
		userID      uuid.UUID
		mockSetup   func(*mockChirpRepository)
		wantErr     error
		description string
	}{
		{
			name:    "successful deletion by owner",
			chirpID: chirpID,
			userID:  ownerID,
			mockSetup: func(m *mockChirpRepository) {
				m.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*domain.Chirp, error) {
					return &domain.Chirp{
						ID:     chirpID,
						Body:   "Test chirp",
						UserID: ownerID,
					}, nil
				}
				m.deleteFunc = func(ctx context.Context, id uuid.UUID) error {
					return nil
				}
			},
			wantErr:     nil,
			description: "Owner should be able to delete their chirp",
		},
		{
			name:    "unauthorized deletion attempt",
			chirpID: chirpID,
			userID:  otherUserID,
			mockSetup: func(m *mockChirpRepository) {
				m.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*domain.Chirp, error) {
					return &domain.Chirp{
						ID:     chirpID,
						Body:   "Test chirp",
						UserID: ownerID,
					}, nil
				}
			},
			wantErr:     domain.ErrUnauthorized,
			description: "Non-owner should not be able to delete chirp",
		},
		{
			name:    "chirp not found",
			chirpID: chirpID,
			userID:  ownerID,
			mockSetup: func(m *mockChirpRepository) {
				m.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*domain.Chirp, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr:     errors.New("failed to get chirp: not found"),
			description: "Should return error when chirp doesn't exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &mockChirpRepository{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}
			service := service.NewChirpService(mockRepo)

			// Execute
			err := service.DeleteChirp(context.Background(), tt.chirpID, tt.userID)

			// Assert
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("%s: expected error %v, got nil", tt.description, tt.wantErr)
					return
				}
				if tt.wantErr == domain.ErrUnauthorized && err != domain.ErrUnauthorized {
					t.Errorf("%s: expected error %v, got %v", tt.description, tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("%s: unexpected error: %v", tt.description, err)
			}
		})
	}
}

func TestChirpService_GetAllChirps(t *testing.T) {
	authorID := uuid.New()

	tests := []struct {
		name        string
		authorID    *uuid.UUID
		sortDesc    bool
		mockSetup   func(*mockChirpRepository)
		wantCount   int
		wantErr     bool
		description string
	}{
		{
			name:     "get all chirps ascending",
			authorID: nil,
			sortDesc: false,
			mockSetup: func(m *mockChirpRepository) {
				m.getAllFunc = func(ctx context.Context) ([]*domain.Chirp, error) {
					return []*domain.Chirp{
						{ID: uuid.New(), Body: "First"},
						{ID: uuid.New(), Body: "Second"},
						{ID: uuid.New(), Body: "Third"},
					}, nil
				}
			},
			wantCount:   3,
			wantErr:     false,
			description: "Should return all chirps in ascending order",
		},
		{
			name:     "get chirps by author",
			authorID: &authorID,
			sortDesc: false,
			mockSetup: func(m *mockChirpRepository) {
				m.getByAuthorFunc = func(ctx context.Context, id uuid.UUID) ([]*domain.Chirp, error) {
					return []*domain.Chirp{
						{ID: uuid.New(), Body: "Author chirp", UserID: authorID},
					}, nil
				}
			},
			wantCount:   1,
			wantErr:     false,
			description: "Should return only author's chirps",
		},
		{
			name:     "get all chirps descending",
			authorID: nil,
			sortDesc: true,
			mockSetup: func(m *mockChirpRepository) {
				m.getAllFunc = func(ctx context.Context) ([]*domain.Chirp, error) {
					return []*domain.Chirp{
						{ID: uuid.New(), Body: "First"},
						{ID: uuid.New(), Body: "Second"},
					}, nil
				}
			},
			wantCount:   2,
			wantErr:     false,
			description: "Should return all chirps in descending order",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &mockChirpRepository{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}
			service := service.NewChirpService(mockRepo)

			// Execute
			chirps, err := service.GetAllChirps(context.Background(), tt.authorID, tt.sortDesc)

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Errorf("%s: expected error, got nil", tt.description)
				}
				return
			}

			if err != nil {
				t.Errorf("%s: unexpected error: %v", tt.description, err)
				return
			}

			if len(chirps) != tt.wantCount {
				t.Errorf("%s: expected %d chirps, got %d", tt.description, tt.wantCount, len(chirps))
			}
		})
	}
}
