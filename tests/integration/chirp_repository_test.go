package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/skylarhoughtongithub/chirpy/internal/auth"
	"github.com/skylarhoughtongithub/chirpy/internal/repository"
)

func TestChirpRepository_Create(t *testing.T) {
	// Setup test database
	testDB := SetupTestDB(t)
	defer testDB.Close()
	defer testDB.Cleanup(t)

	// Create test user
	hashedPassword, _ := auth.HashPassword("password123")
	user := testDB.CreateTestUser(t, "test@example.com", hashedPassword)

	// Create repository
	repo := repository.NewChirpRepository(testDB.Queries)

	tests := []struct {
		name        string
		body        string
		userID      uuid.UUID
		wantErr     bool
		description string
	}{
		{
			name:        "create valid chirp",
			body:        "Hello, world!",
			userID:      user.ID,
			wantErr:     false,
			description: "Should create chirp successfully",
		},
		{
			name:        "create chirp with maximum length",
			body:        "This is exactly 140 characters long string that should be accepted by the system without any issues at all and should pass validation tests perfectly!",
			userID:      user.ID,
			wantErr:     false,
			description: "Should accept chirp at maximum length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chirp, err := repo.Create(context.Background(), tt.body, tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("%s: expected error, got nil", tt.description)
				}
				return
			}

			if err != nil {
				t.Fatalf("%s: unexpected error: %v", tt.description, err)
			}

			// Verify chirp was created
			if chirp.ID == uuid.Nil {
				t.Errorf("%s: chirp ID is nil", tt.description)
			}

			if chirp.Body != tt.body {
				t.Errorf("%s: expected body %q, got %q", tt.description, tt.body, chirp.Body)
			}

			if chirp.UserID != tt.userID {
				t.Errorf("%s: expected userID %v, got %v", tt.description, tt.userID, chirp.UserID)
			}

			// Verify it's in the database
			assert := NewAssertDatabaseState(t, testDB)
			if !assert.ChirpExists(chirp.ID) {
				t.Errorf("%s: chirp not found in database", tt.description)
			}
		})
	}
}

func TestChirpRepository_GetByID(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer testDB.Close()
	defer testDB.Cleanup(t)

	hashedPassword, _ := auth.HashPassword("password123")
	user := testDB.CreateTestUser(t, "test@example.com", hashedPassword)
	dbChirp := testDB.CreateTestChirp(t, "Test chirp", user.ID)

	repo := repository.NewChirpRepository(testDB.Queries)

	// Test: Get existing chirp
	t.Run("get existing chirp", func(t *testing.T) {
		chirp, err := repo.GetByID(context.Background(), dbChirp.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if chirp.ID != dbChirp.ID {
			t.Errorf("expected ID %v, got %v", dbChirp.ID, chirp.ID)
		}

		if chirp.Body != dbChirp.Body {
			t.Errorf("expected body %q, got %q", dbChirp.Body, chirp.Body)
		}
	})

	// Test: Get non-existent chirp
	t.Run("get non-existent chirp", func(t *testing.T) {
		_, err := repo.GetByID(context.Background(), uuid.New())
		if err == nil {
			t.Error("expected error for non-existent chirp, got nil")
		}
	})
}

func TestChirpRepository_GetAll(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer testDB.Close()
	defer testDB.Cleanup(t)

	hashedPassword, _ := auth.HashPassword("password123")
	user := testDB.CreateTestUser(t, "test@example.com", hashedPassword)

	// Create multiple chirps
	testDB.CreateTestChirp(t, "First chirp", user.ID)
	testDB.CreateTestChirp(t, "Second chirp", user.ID)
	testDB.CreateTestChirp(t, "Third chirp", user.ID)

	repo := repository.NewChirpRepository(testDB.Queries)

	// Test
	chirps, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(chirps) != 3 {
		t.Errorf("expected 3 chirps, got %d", len(chirps))
	}

	// Verify all chirps belong to the user
	for _, chirp := range chirps {
		if chirp.UserID != user.ID {
			t.Errorf("expected userID %v, got %v", user.ID, chirp.UserID)
		}
	}
}

func TestChirpRepository_GetByAuthor(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer testDB.Close()
	defer testDB.Cleanup(t)

	hashedPassword, _ := auth.HashPassword("password123")
	user1 := testDB.CreateTestUser(t, "user1@example.com", hashedPassword)
	user2 := testDB.CreateTestUser(t, "user2@example.com", hashedPassword)

	// Create chirps for both users
	testDB.CreateTestChirp(t, "User 1 chirp 1", user1.ID)
	testDB.CreateTestChirp(t, "User 1 chirp 2", user1.ID)
	testDB.CreateTestChirp(t, "User 2 chirp", user2.ID)

	repo := repository.NewChirpRepository(testDB.Queries)

	// Test: Get chirps by user1
	chirps, err := repo.GetByAuthor(context.Background(), user1.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(chirps) != 2 {
		t.Errorf("expected 2 chirps for user1, got %d", len(chirps))
	}

	for _, chirp := range chirps {
		if chirp.UserID != user1.ID {
			t.Errorf("expected all chirps to belong to user1, got chirp with userID %v", chirp.UserID)
		}
	}
}

func TestChirpRepository_Delete(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer testDB.Close()
	defer testDB.Cleanup(t)

	hashedPassword, _ := auth.HashPassword("password123")
	user := testDB.CreateTestUser(t, "test@example.com", hashedPassword)
	dbChirp := testDB.CreateTestChirp(t, "Test chirp", user.ID)

	repo := repository.NewChirpRepository(testDB.Queries)
	assert := NewAssertDatabaseState(t, testDB)

	// Verify chirp exists
	if !assert.ChirpExists(dbChirp.ID) {
		t.Fatal("chirp should exist before deletion")
	}

	// Delete chirp
	err := repo.Delete(context.Background(), dbChirp.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify chirp no longer exists
	if assert.ChirpExists(dbChirp.ID) {
		t.Error("chirp should not exist after deletion")
	}
}
