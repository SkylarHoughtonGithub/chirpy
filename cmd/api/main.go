package main

import (
	"database/sql"
	"log"
	"net/http"
	"sync/atomic"

	_ "github.com/lib/pq"
	"github.com/skylarhoughtongithub/chirpy/internal/api/handlers"
	"github.com/skylarhoughtongithub/chirpy/internal/api/middleware"
	"github.com/skylarhoughtongithub/chirpy/internal/api/routes"
	"github.com/skylarhoughtongithub/chirpy/internal/config"
	"github.com/skylarhoughtongithub/chirpy/internal/database"
	"github.com/skylarhoughtongithub/chirpy/internal/repository"
	"github.com/skylarhoughtongithub/chirpy/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ping database to verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize database queries
	dbQueries := database.New(db)

	// Initialize repositories
	chirpRepo := repository.NewChirpRepository(dbQueries)
	userRepo := repository.NewUserRepository(dbQueries)
	tokenRepo := repository.NewRefreshTokenRepository(dbQueries)

	// Initialize services
	chirpService := service.NewChirpService(chirpRepo)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, tokenRepo, cfg.Auth.JWTSecret)

	// Initialize metrics (for admin)
	metrics := &atomic.Int32{}

	// Initialize handlers
	chirpHandler := handlers.NewChirpHandler(chirpService, cfg.Auth.JWTSecret)
	userHandler := handlers.NewUserHandler(userService, cfg.Auth.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)
	webhookHandler := handlers.NewWebhookHandler(userService, cfg.Auth.PolkaKey)
	adminHandler := handlers.NewAdminHandler(dbQueries, metrics, cfg.Server.Platform)

	// Setup router
	mux := http.NewServeMux()

	// Setup file server with metrics middleware
	fileserverHandler := middleware.MetricsMiddleware(
		metrics,
		http.StripPrefix("/app", http.FileServer(http.Dir("."))),
	)

	// Register all routes
	router := routes.NewRouter(
		chirpHandler,
		userHandler,
		authHandler,
		webhookHandler,
		adminHandler,
	)
	router.Setup(mux, fileserverHandler)

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: middleware.LoggingMiddleware(mux), // Wrap with logging
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(srv.ListenAndServe())
}
