package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/service"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func setupTestMiddleware(t *testing.T) (*service.TokenService, func()) {
	// Create test Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available for testing")
	}

	tokenService := service.NewTokenService(
		"test-secret-key-32-characters-long",
		15*time.Minute,
		30*24*time.Hour,
		rdb,
	)

	cleanup := func() {
		rdb.FlushDB(ctx)
		rdb.Close()
	}

	return tokenService, cleanup
}

func TestAuth_ValidToken(t *testing.T) {
	tokenService, cleanup := setupTestMiddleware(t)
	defer cleanup()

	userID := uuid.New()
	role := "user"

	// Generate valid token
	token, err := tokenService.GenerateAccessToken(userID, role)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	// Create test handler that checks context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserID := GetUserID(r.Context())
		gotRole := GetUserRole(r.Context())

		if gotUserID != userID.String() {
			t.Errorf("Expected UserID %s, got %s", userID.String(), gotUserID)
		}

		if gotRole != role {
			t.Errorf("Expected Role %s, got %s", role, gotRole)
		}

		w.WriteHeader(http.StatusOK)
	})

	// Wrap with Auth middleware
	handler := Auth(tokenService)(testHandler)

	// Create request with Authorization header
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestAuth_MissingHeader(t *testing.T) {
	tokenService, cleanup := setupTestMiddleware(t)
	defer cleanup()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	handler := Auth(tokenService)(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rec.Code)
	}
}

func TestAuth_InvalidHeaderFormat(t *testing.T) {
	tokenService, cleanup := setupTestMiddleware(t)
	defer cleanup()

	testCases := []struct {
		name   string
		header string
	}{
		{"no bearer prefix", "token123"},
		{"wrong prefix", "Basic token123"},
		{"only bearer", "Bearer"},
		{"empty bearer", "Bearer "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("Handler should not be called")
			})

			handler := Auth(tokenService)(testHandler)

			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tc.header)

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("Expected status 401, got %d", rec.Code)
			}
		})
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	tokenService, cleanup := setupTestMiddleware(t)
	defer cleanup()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	handler := Auth(tokenService)(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rec.Code)
	}
}

func TestAuth_ExpiredToken(t *testing.T) {
	// Create service with very short TTL
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1,
	})
	defer rdb.Close()

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available for testing")
	}
	defer rdb.FlushDB(ctx)

	tokenService := service.NewTokenService(
		"test-secret-key-32-characters-long",
		1*time.Millisecond, // Very short TTL
		30*24*time.Hour,
		rdb,
	)

	userID := uuid.New()
	token, _ := tokenService.GenerateAccessToken(userID, "user")

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	handler := Auth(tokenService)(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rec.Code)
	}
}

func TestGetUserID_NoContext(t *testing.T) {
	ctx := context.Background()
	userID := GetUserID(ctx)
	if userID != "" {
		t.Errorf("Expected empty string for missing context, got %s", userID)
	}
}

func TestGetUserRole_NoContext(t *testing.T) {
	ctx := context.Background()
	role := GetUserRole(ctx)
	if role != "" {
		t.Errorf("Expected empty string for missing context, got %s", role)
	}
}
