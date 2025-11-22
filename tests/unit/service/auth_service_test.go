package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/skylarhoughtongithub/chirpy/internal/auth"
	"github.com/skylarhoughtongithub/chirpy/internal/domain"
	"github.com/skylarhoughtongithub/chirpy/internal/service"
)

// Mock refresh token repository
type mockRefreshTokenRepository struct {
	createFunc      func(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error
	getUserByTokenFunc func(ctx context.Context, token string) (*domain.User, error)
	revokeFunc      func(ctx context.Context, token string) error
}

func (m *mockRefreshTokenRepository) Create(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, token, userID, expiresAt)
	}
	return errors.New("not implemented")
}

func (m *mockRefreshTokenRepository) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	if m.getUserByTokenFunc != nil {
		return m.getUserByTokenFunc(ctx, token)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	if m.revokeFunc != nil {
		return m.revokeFunc(ctx, token)
	}
	return errors.New("not implemented")
}

func TestAuthService_Login(t *testing.T) {
	// Create a test user with hashed password
	testPassword := "password123"
	hashedPassword, _ := auth.HashPassword(testPassword)
	testUser := &domain.User{
		ID:             uuid.New(),
		Email:          "test@example.com",
		HashedPassword: hashedPassword,
		IsChirpyRed:    false,
	}

	tests := []struct {
		name          string
		email         string
		password      string
		mockUserRepo  func(*mockUserRepository)
		mockTokenRepo func(*mockRefreshTokenRepository)
		wantErr       error
		description   string
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: testPassword,
			mockUserRepo: func(m *mockUserRepository) {
				m.getByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return testUser, nil
				}
			},
			mockTokenRepo: func(m *mockRefreshTokenRepository) {
				m.createFunc = func(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error {
					return nil
				}
			},
			wantErr:     nil,
			description: "Should login successfully with correct credentials",
		},
		{
			name:     "user not found",
			email:    "nonexistent@example.com",
			password: testPassword,
			mockUserRepo: func(m *mockUserRepository) {
				m.getByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return nil, errors.New("user not found")
				}
			},
			mockTokenRepo: func(m *mockRefreshTokenRepository) {},
			wantErr:       domain.ErrInvalidCredentials,
			description:   "Should return invalid credentials for non-existent user",
		},
		{
			name:     "wrong password",
			email:    "test@example.com",
			password: "wrongpassword",
			mockUserRepo: func(m *mockUserRepository) {
				m.getByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return testUser, nil
				}
			},
			mockTokenRepo: func(m *mockRefreshTokenRepository) {},
			wantErr:       domain.ErrInvalidCredentials,
			description:   "Should return invalid credentials for wrong password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUserRepo := &mockUserRepository{}
			mockTokenRepo := &mockRefreshTokenRepository{}
			if tt.mockUserRepo != nil {
				tt.mockUserRepo(mockUserRepo)
			}
			if tt.mockTokenRepo != nil {
				tt.mockTokenRepo(mockTokenRepo)
			}

			service := service.NewAuthService(mockUserRepo, mockTokenRepo, "test-secret")

			// Execute
			loginResp, err := service.Login(context.Background(), tt.email, tt.password)

			// Assert
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("%s: expected error %v, got nil", tt.description, tt.wantErr)
					return
				}
				if err != tt.wantErr {
					t.Errorf("%s: expected error %v, got %v", tt.description, tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("%s: unexpected error: %v", tt.description, err)
				return
			}

			// Verify tokens were created
			if loginResp.AccessToken == "" {
				t.Errorf("%s: access token is empty", tt.description)
			}

			if loginResp.RefreshToken == "" {
				t.Errorf("%s: refresh token is empty", tt.description)
			}

			if loginResp.User.Email != tt.email {
				t.Errorf("%s: expected email %q, got %q", tt.description, tt.email, loginResp.User.Email)
			}
		})
	}
}

func TestAuthService_RefreshAccessToken(t *testing.T) {
	testUser := &domain.User{
		ID:             uuid.New(),
		Email:          "test@example.com",
		HashedPassword: "hashed",
		IsChirpyRed:    false,
	}

	tests := []struct {
		name          string
		refreshToken  string
		mockUserRepo  func(*mockUserRepository)
		mockTokenRepo func(*mockRefreshTokenRepository)
		wantErr       error
		description   string
	}{
		{
			name:         "successful token refresh",
			refreshToken: "valid-token",
			mockUserRepo: func(m *mockUserRepository) {},
			mockTokenRepo: func(m *mockRefreshTokenRepository) {
				m.getUserByTokenFunc = func(ctx context.Context, token string) (*domain.User, error) {
					return testUser, nil
				}
			},
			wantErr:     nil,
			description: "Should create new access token from valid refresh token",
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid-token",
			mockUserRepo: func(m *mockUserRepository) {},
			mockTokenRepo: func(m *mockRefreshTokenRepository) {
				m.getUserByTokenFunc = func(ctx context.Context, token string) (*domain.User, error) {
					return nil, errors.New("token not found")
				}
			},
			wantErr:     domain.ErrInvalidToken,
			description: "Should reject invalid refresh token",
		},
		{
			name:         "expired refresh token",
			refreshToken: "expired-token",
			mockUserRepo: func(m *mockUserRepository) {},
			mockTokenRepo: func(m *mockRefreshTokenRepository) {
				m.getUserByTokenFunc = func(ctx context.Context, token string) (*domain.User, error) {
					return nil, errors.New("token expired")
				}
			},
			wantErr:     domain.ErrInvalidToken,
			description: "Should reject expired refresh token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUserRepo := &mockUserRepository{}
			mockTokenRepo := &mockRefreshTokenRepository{}
			if tt.mockUserRepo != nil {
				tt.mockUserRepo(mockUserRepo)
			}
			if tt.mockTokenRepo != nil {
				tt.mockTokenRepo(mockTokenRepo)
			}

			service := service.NewAuthService(mockUserRepo, mockTokenRepo, "test-secret")

			// Execute
			accessToken, err := service.RefreshAccessToken(context.Background(), tt.refreshToken)

			// Assert
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("%s: expected error %v, got nil", tt.description, tt.wantErr)
					return
				}
				if err != tt.wantErr {
					t.Errorf("%s: expected error %v, got %v", tt.description, tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("%s: unexpected error: %v", tt.description, err)
				return
			}

			if accessToken == "" {
				t.Errorf("%s: access token is empty", tt.description)
			}

			// Verify the access token is valid
			userID, err := auth.ValidateJWT(accessToken, "test-secret")
			if err != nil {
				t.Errorf("%s: access token is invalid: %v", tt.description, err)
			}

			if userID != testUser.ID {
				t.Errorf("%s: expected userID %v, got %v", tt.description, testUser.ID, userID)
			}
		})
	}
}

func TestAuthService_RevokeRefreshToken(t *testing.T) {
	tests := []struct {
		name          string
		refreshToken  string
		mockTokenRepo func(*mockRefreshTokenRepository)
		wantErr       bool
		description   string
	}{
		{
			name:         "successful revocation",
			refreshToken: "valid-token",
			mockTokenRepo: func(m *mockRefreshTokenRepository) {
				m.revokeFunc = func(ctx context.Context, token string) error {
					return nil
				}
			},
			wantErr:     false,
			description: "Should revoke token successfully",
		},
		{
			name:         "token not found",
			refreshToken: "nonexistent-token",
			mockTokenRepo: func(m *mockRefreshTokenRepository) {
				m.revokeFunc = func(ctx context.Context, token string) error {
					return errors.New("token not found")
				}
			},
			wantErr:     true,
			description: "Should handle non-existent token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUserRepo := &mockUserRepository{}
			mockTokenRepo := &mockRefreshTokenRepository{}
			if tt.mockTokenRepo != nil {
				tt.mockTokenRepo(mockTokenRepo)
			}

			service := service.NewAuthService(mockUserRepo, mockTokenRepo, "test-secret")

			// Execute
			err := service.RevokeRefreshToken(context.Background(), tt.refreshToken)

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
