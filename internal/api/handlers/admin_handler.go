package handlers

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/skylarhoughtongithub/chirpy/internal/database"
)

// AdminHandler handles admin HTTP requests
type AdminHandler struct {
	db       *database.Queries
	metrics  *atomic.Int32
	platform string
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(db *database.Queries, metrics *atomic.Int32, platform string) *AdminHandler {
	return &AdminHandler{
		db:       db,
		metrics:  metrics,
		platform: platform,
	}
}

// Metrics handles GET /admin/metrics
func (h *AdminHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	html := fmt.Sprintf(`
<html>
<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>
</html>
	`, h.metrics.Load())

	w.Write([]byte(html))
}

// Reset handles POST /admin/reset
func (h *AdminHandler) Reset(w http.ResponseWriter, r *http.Request) {
	// Only allow reset in development
	if h.platform != "dev" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Delete all users (cascades to chirps and tokens)
	err := h.db.DeleteAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Could not reset database", http.StatusInternalServerError)
		return
	}

	// Reset metrics
	h.metrics.Store(0)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database reset and hits reset to 0"))
}
