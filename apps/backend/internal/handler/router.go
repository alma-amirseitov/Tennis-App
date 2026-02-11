package handler

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/handler/middleware"
)

func NewRouter(logger *slog.Logger, db *pgxpool.Pool, redis *goredis.Client) *chi.Mux {
	r := chi.NewRouter()

	// Middleware chain
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recovery(logger))
	r.Use(cors.Handler(middleware.CORS()))

	// Health check
	health := NewHealthHandler(db, redis)
	r.Get("/health", health.Check)

	// API v1 routes
	r.Route("/v1", func(r chi.Router) {
		// TODO: add routes
	})

	return r
}
