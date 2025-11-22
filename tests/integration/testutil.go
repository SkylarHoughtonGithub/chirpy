package integration

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/skylarhoughtongithub/chirpy/internal/database"
)

// TestDB wraps a test database connection
type TestDB struct {
	DB      *sql.DB
	Queries *database.Queries
}

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()

	// Get test database URL from environment
	dbURL := os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/chirpy_test?sslmode=disable"
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Ping to verify connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	queries := database.New(db)

	return &TestDB{
		DB:      db,
		Queries: queries,
	}
}

// Cleanup cleans up the test database
func (tdb *TestDB) Cleanup(t *testing.T) {
	t.Helper()

	// Clean all tables
	ctx := context.Background()
	
	// Delete in correct order due to foreign keys
	_, err := tdb.DB.ExecContext(ctx, "DELETE FROM chirps")
	if err != nil {
		t.Logf("Warning: Failed to clean chirps table: %v", err)
	}

	_, err = tdb.DB.ExecContext(ctx, "DELETE FROM refresh_tokens")
	if err != nil {
		t.Logf("Warning: Failed to clean refresh_tokens table: %v", err)
	}

	_, err = tdb.DB.ExecContext(ctx, "DELETE FROM users")
	if err != nil {
		t.Logf("Warning: Failed to clean users table: %v", err)
	}
}

// Close closes the database connection
func (tdb *TestDB) Close() {
	tdb.DB.Close()
}

// CreateTestUser creates a test user
func (tdb *TestDB) CreateTestUser(t *testing.T, email, hashedPassword string) *database.User {
	t.Helper()

	user, err := tdb.Queries.CreateUser(context.Background(), database.CreateUserParams{
		Email:          email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return &user
}

// CreateTestChirp creates a test chirp
func (tdb *TestDB) CreateTestChirp(t *testing.T, body string, userID uuid.UUID) *database.Chirp {
	t.Helper()

	chirp, err := tdb.Queries.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   body,
		UserID: userID,
	})
	if err != nil {
		t.Fatalf("Failed to create test chirp: %v", err)
	}

	return &chirp
}

// AssertDatabaseState helps verify database state
type AssertDatabaseState struct {
	t  *testing.T
	db *TestDB
}

// NewAssertDatabaseState creates a new assertion helper
func NewAssertDatabaseState(t *testing.T, db *TestDB) *AssertDatabaseState {
	return &AssertDatabaseState{t: t, db: db}
}

// UserExists checks if a user exists by email
func (a *AssertDatabaseState) UserExists(email string) bool {
	a.t.Helper()

	_, err := a.db.Queries.GetUserByEmail(context.Background(), email)
	return err == nil
}

// ChirpExists checks if a chirp exists by ID
func (a *AssertDatabaseState) ChirpExists(id uuid.UUID) bool {
	a.t.Helper()

	_, err := a.db.Queries.GetChirp(context.Background(), id)
	return err == nil
}

// ChirpCount returns the total number of chirps
func (a *AssertDatabaseState) ChirpCount() int {
	a.t.Helper()

	chirps, err := a.db.Queries.GetAllChirps(context.Background())
	if err != nil {
		a.t.Fatalf("Failed to count chirps: %v", err)
	}

	return len(chirps)
}

// UserCount returns the total number of users
func (a *AssertDatabaseState) UserCount() int {
	a.t.Helper()

	var count int
	err := a.db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		a.t.Fatalf("Failed to count users: %v", err)
	}

	return count
}
