package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewHealthHandler(db *pgxpool.Pool, redis *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, redis: redis}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	dbStatus := "connected"
	if err := h.db.Ping(ctx); err != nil {
		dbStatus = "disconnected"
	}

	redisStatus := "connected"
	if err := h.redis.Ping(ctx).Err(); err != nil {
		redisStatus = "disconnected"
	}

	status := http.StatusOK
	statusText := "ok"
	if dbStatus != "connected" || redisStatus != "connected" {
		status = http.StatusServiceUnavailable
		statusText = "degraded"
	}

	respondJSON(w, status, map[string]any{
		"status":   statusText,
		"version":  "0.1.0",
		"database": dbStatus,
		"redis":    redisStatus,
	})
}
