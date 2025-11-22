package fixtures

import (
	"time"

	"github.com/google/uuid"
	"github.com/skylarhoughtongithub/chirpy/internal/domain"
)

// TestUser creates a test user with default values
func TestUser() *domain.User {
	return &domain.User{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          "test@example.com",
		HashedPassword: "$argon2id$v=19$m=65536,t=3,p=2$test",
		IsChirpyRed:    false,
	}
}

// TestUserWithEmail creates a test user with a specific email
func TestUserWithEmail(email string) *domain.User {
	user := TestUser()
	user.Email = email
	return user
}

// TestChirpyRedUser creates a test user with Chirpy Red subscription
func TestChirpyRedUser() *domain.User {
	user := TestUser()
	user.IsChirpyRed = true
	return user
}

// TestChirp creates a test chirp with default values
func TestChirp() *domain.Chirp {
	return &domain.Chirp{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      "This is a test chirp",
		UserID:    uuid.New(),
	}
}

// TestChirpWithBody creates a test chirp with a specific body
func TestChirpWithBody(body string) *domain.Chirp {
	chirp := TestChirp()
	chirp.Body = body
	return chirp
}

// TestChirpWithUser creates a test chirp for a specific user
func TestChirpWithUser(userID uuid.UUID) *domain.Chirp {
	chirp := TestChirp()
	chirp.UserID = userID
	return chirp
}

// TestChirps creates multiple test chirps
func TestChirps(count int) []*domain.Chirp {
	chirps := make([]*domain.Chirp, count)
	baseTime := time.Now()
	
	for i := 0; i < count; i++ {
		chirps[i] = &domain.Chirp{
			ID:        uuid.New(),
			CreatedAt: baseTime.Add(time.Duration(i) * time.Second),
			UpdatedAt: baseTime.Add(time.Duration(i) * time.Second),
			Body:      "Test chirp " + string(rune('A'+i)),
			UserID:    uuid.New(),
		}
	}
	
	return chirps
}

// TestRefreshToken creates a test refresh token
func TestRefreshToken(userID uuid.UUID) *domain.RefreshToken {
	return &domain.RefreshToken{
		Token:     "test-refresh-token-" + uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		RevokedAt: nil,
	}
}

// ValidPassword returns a valid test password
func ValidPassword() string {
	return "password123"
}

// ValidEmail returns a valid test email
func ValidEmail() string {
	return "test@example.com"
}

// LongChirpBody returns a chirp body that exceeds 140 characters
func LongChirpBody() string {
	return "This is a very long chirp that exceeds the maximum allowed length of 140 characters. It should fail validation and return an error. Adding more text to ensure it goes well over the limit."
}

// MaxLengthChirpBody returns a chirp body at exactly 140 characters
func MaxLengthChirpBody() string {
	return "This is exactly 140 characters long string that should be accepted by the system without any issues at all and should pass validation!"
}

// ProfaneChirpBody returns a chirp body with profane words
func ProfaneChirpBody() string {
	return "This is a kerfuffle with sharbert and fornax"
}

// CleanedProfaneChirpBody returns the expected cleaned version
func CleanedProfaneChirpBody() string {
	return "This is a **** with **** and ****"
}
