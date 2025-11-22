package routes

import (
	"net/http"

	"github.com/skylarhoughtongithub/chirpy/internal/api/handlers"
)

// Router holds all handlers and sets up routes
type Router struct {
	chirpHandler   *handlers.ChirpHandler
	userHandler    *handlers.UserHandler
	authHandler    *handlers.AuthHandler
	webhookHandler *handlers.WebhookHandler
	adminHandler   *handlers.AdminHandler
}

// NewRouter creates a new router with all handlers
func NewRouter(
	chirpHandler *handlers.ChirpHandler,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	webhookHandler *handlers.WebhookHandler,
	adminHandler *handlers.AdminHandler,
) *Router {
	return &Router{
		chirpHandler:   chirpHandler,
		userHandler:    userHandler,
		authHandler:    authHandler,
		webhookHandler: webhookHandler,
		adminHandler:   adminHandler,
	}
}

// Setup configures all application routes
func (rt *Router) Setup(mux *http.ServeMux, fileserverHandler http.Handler) {
	// Static files
	mux.Handle("/app/", fileserverHandler)

	// Health check
	mux.HandleFunc("GET /api/healthz", handlers.Healthz)

	// Authentication routes
	mux.HandleFunc("POST /api/login", rt.authHandler.Login)
	mux.HandleFunc("POST /api/refresh", rt.authHandler.Refresh)
	mux.HandleFunc("POST /api/revoke", rt.authHandler.Revoke)

	// User routes
	mux.HandleFunc("POST /api/users", rt.userHandler.Create)
	mux.HandleFunc("PUT /api/users", rt.userHandler.Update)

	// Chirp routes
	mux.HandleFunc("POST /api/chirps", rt.chirpHandler.CreateChirp)
	mux.HandleFunc("GET /api/chirps", rt.chirpHandler.GetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", rt.chirpHandler.GetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", rt.chirpHandler.DeleteChirp)

	// Webhook routes
	mux.HandleFunc("POST /api/polka/webhooks", rt.webhookHandler.HandlePolkaWebhook)

	// Admin routes
	mux.HandleFunc("GET /admin/metrics", rt.adminHandler.Metrics)
	mux.HandleFunc("POST /admin/reset", rt.adminHandler.Reset)
}
