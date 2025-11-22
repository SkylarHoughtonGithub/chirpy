package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/skylarhoughtongithub/chirpy/internal/domain"
	"github.com/skylarhoughtongithub/chirpy/internal/service"
)

// Mock user repository
type mockUserRepository struct {
	createFunc           func(ctx context.Context, email, hashedPassword string) (*domain.User, error)
	getByEmailFunc       func(ctx context.Context, email string) (*domain.User, error)
	getByIDFunc          func(ctx context.Context, id uuid.UUID) (*domain.User, error)
	updateFunc           func(ctx context.Context, id uuid.UUID, email, hashedPassword string) (*domain.User, error)
	upgradeToChirpyFunc  func(ctx context.Context, id uuid.UUID) error
	deleteAllFunc        func(ctx context.Context) error
}

func (m *mockUserRepository) Create(ctx context.Context, email, hashedPassword string) (*domain.User, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, email, hashedPassword)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(ctx, email)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) Update(ctx context.Context, id uuid.UUID, email, hashedPassword string) (*domain.User, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, email, hashedPassword)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) UpgradeToChirpyRed(ctx context.Context, id uuid.UUID) error {
	if m.upgradeToChirpyFunc != nil {
		return m.upgradeToChirpyFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockUserRepository) DeleteAll(ctx context.Context) error {
	if m.deleteAllFunc != nil {
		return m.deleteAllFunc(ctx)
	}
	return errors.New("not implemented")
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		mockSetup   func(*mockUserRepository)
		wantErr     error
		description string
	}{
		{
			name:     "successful user creation",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(m *mockUserRepository) {
				m.createFunc = func(ctx context.Context, email, hashedPassword string) (*domain.User, error) {
					return &domain.User{
						ID:             uuid.New(),
						Email:          email,
						HashedPassword: hashedPassword,
						IsChirpyRed:    false,
					}, nil
				}
			},
			wantErr:     nil,
			description: "Should create user with hashed password",
		},
		{
			name:     "empty email",
			email:    "",
			password: "password123",
			mockSetup: func(m *mockUserRepository) {
				// Should not be called
			},
			wantErr:     domain.ErrInvalidInput,
			description: "Should reject empty email",
		},
		{
			name:     "empty password",
			email:    "test@example.com",
			password: "",
			mockSetup: func(m *mockUserRepository) {
				// Should not be called
			},
			wantErr:     domain.ErrInvalidInput,
			description: "Should reject empty password",
		},
		{
			name:     "duplicate email",
			email:    "existing@example.com",
			password: "password123",
			mockSetup: func(m *mockUserRepository) {
				m.createFunc = func(ctx context.Context, email, hashedPassword string) (*domain.User, error) {
					return nil, errors.New("duplicate key")
				}
			},
			wantErr:     errors.New("failed to create user: duplicate key"),
			description: "Should handle duplicate email error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &mockUserRepository{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}
			service := service.NewUserService(mockRepo)

			// Execute
			user, err := service.CreateUser(context.Background(), tt.email, tt.password)

			// Assert
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("%s: expected error %v, got nil", tt.description, tt.wantErr)
					return
				}
				if tt.wantErr == domain.ErrInvalidInput && err != domain.ErrInvalidInput {
					t.Errorf("%s: expected error %v, got %v", tt.description, tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("%s: unexpected error: %v", tt.description, err)
				return
			}

			if user.Email != tt.email {
				t.Errorf("%s: expected email %q, got %q", tt.description, tt.email, user.Email)
			}

			// Verify password was hashed (should not be plain text)
			if user.HashedPassword == tt.password {
				t.Errorf("%s: password was not hashed", tt.description)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name        string
		userID      uuid.UUID
		email       string
		password    string
		mockSetup   func(*mockUserRepository)
		wantErr     error
		description string
	}{
		{
			name:     "successful update",
			userID:   userID,
			email:    "updated@example.com",
			password: "newpassword123",
			mockSetup: func(m *mockUserRepository) {
				m.updateFunc = func(ctx context.Context, id uuid.UUID, email, hashedPassword string) (*domain.User, error) {
					return &domain.User{
						ID:             id,
						Email:          email,
						HashedPassword: hashedPassword,
						IsChirpyRed:    false,
					}, nil
				}
			},
			wantErr:     nil,
			description: "Should update user successfully",
		},
		{
			name:     "empty email",
			userID:   userID,
			email:    "",
			password: "password123",
			mockSetup: func(m *mockUserRepository) {
				// Should not be called
			},
			wantErr:     domain.ErrInvalidInput,
			description: "Should reject empty email",
		},
		{
			name:     "empty password",
			userID:   userID,
			email:    "test@example.com",
			password: "",
			mockSetup: func(m *mockUserRepository) {
				// Should not be called
			},
			wantErr:     domain.ErrInvalidInput,
			description: "Should reject empty password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &mockUserRepository{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}
			service := service.NewUserService(mockRepo)

			// Execute
			user, err := service.UpdateUser(context.Background(), tt.userID, tt.email, tt.password)

			// Assert
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("%s: expected error %v, got nil", tt.description, tt.wantErr)
					return
				}
				if tt.wantErr == domain.ErrInvalidInput && err != domain.ErrInvalidInput {
					t.Errorf("%s: expected error %v, got %v", tt.description, tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("%s: unexpected error: %v", tt.description, err)
				return
			}

			if user.Email != tt.email {
				t.Errorf("%s: expected email %q, got %q", tt.description, tt.email, user.Email)
			}

			if user.ID != tt.userID {
				t.Errorf("%s: expected userID %v, got %v", tt.description, tt.userID, user.ID)
			}
		})
	}
}

func TestUserService_UpgradeToChirpyRed(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name        string
		userID      uuid.UUID
		mockSetup   func(*mockUserRepository)
		wantErr     bool
		description string
	}{
		{
			name:   "successful upgrade",
			userID: userID,
			mockSetup: func(m *mockUserRepository) {
				m.upgradeToChirpyFunc = func(ctx context.Context, id uuid.UUID) error {
					return nil
				}
			},
			wantErr:     false,
			description: "Should upgrade user to Chirpy Red",
		},
		{
			name:   "user not found",
			userID: userID,
			mockSetup: func(m *mockUserRepository) {
				m.upgradeToChirpyFunc = func(ctx context.Context, id uuid.UUID) error {
					return errors.New("user not found")
				}
			},
			wantErr:     true,
			description: "Should return error if user doesn't exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &mockUserRepository{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}
			service := service.NewUserService(mockRepo)

			// Execute
			err := service.UpgradeToChirpyRed(context.Background(), tt.userID)

			// Assert
			if tt.wantErr && err == nil {
				t.Errorf("%s: expected error, got nil", tt.description)
				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("%s: unexpected error: %v", tt.description, err)
			}
		})
	}
}
