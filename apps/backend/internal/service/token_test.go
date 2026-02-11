package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func setupTestTokenService(t *testing.T) (*TokenService, *redis.Client, func()) {
	// Create test Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // Use test DB
	})

	// Ping to check connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available for testing")
	}

	tokenService := NewTokenService(
		"test-secret-key-32-characters-long",
		15*time.Minute,
		30*24*time.Hour,
		rdb,
	)

	cleanup := func() {
		// Clean up test data
		rdb.FlushDB(ctx)
		rdb.Close()
	}

	return tokenService, rdb, cleanup
}

func TestGenerateAndValidateAccessToken(t *testing.T) {
	service, _, cleanup := setupTestTokenService(t)
	defer cleanup()

	userID := uuid.New()
	role := "user"

	// Generate token
	token, err := service.GenerateAccessToken(userID, role)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}

	// Validate token
	claims, err := service.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken failed: %v", err)
	}

	if claims.UserID != userID.String() {
		t.Errorf("Expected UserID %s, got %s", userID.String(), claims.UserID)
	}

	if claims.Role != role {
		t.Errorf("Expected Role %s, got %s", role, claims.Role)
	}
}

func TestValidateAccessToken_Invalid(t *testing.T) {
	service, _, cleanup := setupTestTokenService(t)
	defer cleanup()

	testCases := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"malformed token", "invalid.token.here"},
		{"wrong secret", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.ValidateAccessToken(tc.token)
			if err == nil {
				t.Error("Expected error for invalid token")
			}
		})
	}
}

func TestGenerateAndValidateRefreshToken(t *testing.T) {
	service, rdb, cleanup := setupTestTokenService(t)
	defer cleanup()

	ctx := context.Background()
	userID := uuid.New()

	// Generate refresh token
	token, err := service.GenerateRefreshToken(ctx, userID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	// Validate refresh token
	claims, err := service.ValidateRefreshToken(ctx, token)
	if err != nil {
		t.Fatalf("ValidateRefreshToken failed: %v", err)
	}

	if claims.UserID != userID.String() {
		t.Errorf("Expected UserID %s, got %s", userID.String(), claims.UserID)
	}

	// Verify JTI is stored in Redis
	key := "refresh:" + claims.JTI
	exists, _ := rdb.Exists(ctx, key).Result()
	if exists == 0 {
		t.Error("JTI should be stored in Redis")
	}
}

func TestRevokeRefreshToken(t *testing.T) {
	service, rdb, cleanup := setupTestTokenService(t)
	defer cleanup()

	ctx := context.Background()
	userID := uuid.New()

	// Generate refresh token
	token, err := service.GenerateRefreshToken(ctx, userID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	claims, _ := service.ValidateRefreshToken(ctx, token)

	// Revoke token
	err = service.RevokeRefreshToken(ctx, claims.JTI)
	if err != nil {
		t.Fatalf("RevokeRefreshToken failed: %v", err)
	}

	// Verify JTI is deleted from Redis
	key := "refresh:" + claims.JTI
	exists, _ := rdb.Exists(ctx, key).Result()
	if exists != 0 {
		t.Error("JTI should be deleted from Redis")
	}

	// Validate should now fail
	_, err = service.ValidateRefreshToken(ctx, token)
	if err != ErrTokenRevoked {
		t.Errorf("Expected ErrTokenRevoked, got: %v", err)
	}
}

func TestRefreshTokens_Valid(t *testing.T) {
	service, _, cleanup := setupTestTokenService(t)
	defer cleanup()

	ctx := context.Background()
	userID := uuid.New()

	// Generate initial refresh token
	oldRefreshToken, err := service.GenerateRefreshToken(ctx, userID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	// Refresh tokens
	newAccessToken, newRefreshToken, err := service.RefreshTokens(ctx, oldRefreshToken)
	if err != nil {
		t.Fatalf("RefreshTokens failed: %v", err)
	}

	if newAccessToken == "" || newRefreshToken == "" {
		t.Error("Expected non-empty tokens")
	}

	// Old refresh token should be revoked
	_, err = service.ValidateRefreshToken(ctx, oldRefreshToken)
	if err != ErrTokenRevoked {
		t.Errorf("Expected ErrTokenRevoked for old token, got: %v", err)
	}

	// New refresh token should be valid
	_, err = service.ValidateRefreshToken(ctx, newRefreshToken)
	if err != nil {
		t.Errorf("New refresh token should be valid, got error: %v", err)
	}
}

func TestRefreshTokens_Expired(t *testing.T) {
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

	service := NewTokenService(
		"test-secret-key-32-characters-long",
		15*time.Minute,
		1*time.Millisecond, // Very short TTL
		rdb,
	)

	userID := uuid.New()

	// Generate refresh token
	refreshToken, err := service.GenerateRefreshToken(ctx, userID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Try to refresh - should fail
	_, _, err = service.RefreshTokens(ctx, refreshToken)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}

func TestRefreshTokens_Reused(t *testing.T) {
	service, _, cleanup := setupTestTokenService(t)
	defer cleanup()

	ctx := context.Background()
	userID := uuid.New()

	// Generate initial refresh token
	refreshToken, err := service.GenerateRefreshToken(ctx, userID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	// Use token once
	_, _, err = service.RefreshTokens(ctx, refreshToken)
	if err != nil {
		t.Fatalf("First RefreshTokens failed: %v", err)
	}

	// Try to reuse same token - should fail with TOKEN_REVOKED
	_, _, err = service.RefreshTokens(ctx, refreshToken)
	if err != ErrTokenRevoked {
		t.Errorf("Expected ErrTokenRevoked for reused token, got: %v", err)
	}
}
