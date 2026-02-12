package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/config"
)

type seedUser struct {
	Phone     string
	FirstName string
	LastName  string
	Rating    float64
	Level     float64
}

func main() {
	// Load .env for local development (ignore error if not present).
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}

	users := []seedUser{
		{Phone: "+77070000001", FirstName: "Аян", LastName: "Тестовый", Rating: 1000, Level: 2.5},
		{Phone: "+77070000002", FirstName: "Марина", LastName: "Тестовая", Rating: 1150, Level: 3.0},
		{Phone: "+77070000003", FirstName: "Руслан", LastName: "Тестовый", Rating: 1250, Level: 3.5},
		{Phone: "+77070000004", FirstName: "Айгерим", LastName: "Тестовая", Rating: 1350, Level: 4.0},
		{Phone: "+77070000005", FirstName: "Марат", LastName: "Тестовый", Rating: 1500, Level: 4.5},
	}

	for _, u := range users {
		_, err := pool.Exec(ctx, `
			INSERT INTO users (
				phone,
				phone_verified,
				first_name,
				last_name,
				city,
				district,
				ntrp_level,
				level_label,
				global_rating,
				global_games_count,
				status,
				is_profile_complete,
				quiz_completed
			) VALUES (
				$1,
				TRUE,
				$2,
				$3,
				'Астана',
				'Есильский',
				$4,
				'Тестовый игрок',
				$5,
				0,
				'active',
				TRUE,
				TRUE
			)
			ON CONFLICT (phone) DO NOTHING;
		`, u.Phone, u.FirstName, u.LastName, u.Level, u.Rating)
		if err != nil {
			slog.Error("failed to insert seed user", "phone", u.Phone, "error", err)
		} else {
			slog.Info("seeded user", "phone", u.Phone, "first_name", u.FirstName)
		}
	}
}
