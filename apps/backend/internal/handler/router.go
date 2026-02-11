package handler

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/config"
	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/handler/middleware"
	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/pkg/validator"
	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/repository"
	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/service"
)

func NewRouter(logger *slog.Logger, db *pgxpool.Pool, redis *goredis.Client, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// Middleware chain
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recovery(logger))
	r.Use(cors.Handler(middleware.CORS()))

	// Health check
	health := NewHealthHandler(db, redis)
	r.Get("/health", health.Check)

	// Initialize services
	queries := repository.New(db)
	tokenService := service.NewTokenService(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL, redis)
	authService := service.NewAuthService(queries, redis, tokenService, cfg.Environment)
	quizService := service.NewQuizService(queries)

	// Initialize validator
	v := validator.New()

	// Initialize handlers
	authHandler := NewAuthHandler(authService, tokenService, v)
	quizHandler := NewQuizHandler(quizService, v)

	// API v1 routes
	r.Route("/v1", func(r chi.Router) {
		// Public routes (no auth required)
		r.Group(func(r chi.Router) {
			// Rate limit for auth endpoints
			r.Use(middleware.RateLimiter(redis, 100, time.Minute, middleware.IPKeyFunc("api_general")))

			// Auth endpoints
			r.Post("/auth/otp/send", authHandler.SendOTP)
			r.Post("/auth/otp/verify", authHandler.VerifyOTP)
			r.Post("/auth/refresh", authHandler.RefreshToken)
		})

		// Protected routes (auth required)
		r.Group(func(r chi.Router) {
			// Auth middleware
			r.Use(middleware.Auth(tokenService))

			// Rate limit for authenticated users
			r.Use(middleware.RateLimiter(redis, 100, time.Minute, middleware.UserKeyFunc("api_user")))

			// Profile setup (temp_token can access this)
			r.Post("/auth/profile/setup", authHandler.ProfileSetup)

			// Quiz endpoints
			r.Get("/quiz", quizHandler.GetQuestions)
			r.Post("/quiz", quizHandler.SubmitAnswers)

			// TODO: Add other protected routes here
			// Example:
			// r.Get("/users/me", userHandler.GetMe)
		})
	})

	return r
}
