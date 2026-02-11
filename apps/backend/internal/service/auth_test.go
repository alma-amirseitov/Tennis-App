package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func setupTestAuthService(t *testing.T) (*AuthService, *redis.Client, func()) {
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

	// Create token service
	tokenService := NewTokenService(
		"test-secret-key-32-characters-long",
		15*time.Minute,
		30*24*time.Hour,
		rdb,
	)

	// Create auth service (with nil repo for these tests)
	authService := NewAuthService(
		nil, // repo
		rdb,
		tokenService,
		"development",
	)

	cleanup := func() {
		// Clean up test data
		rdb.FlushDB(ctx)
		rdb.Close()
	}

	return authService, rdb, cleanup
}

func TestSendOTP_ValidPhone(t *testing.T) {
	service, _, cleanup := setupTestAuthService(t)
	defer cleanup()

	ctx := context.Background()
	phone := "+77071234567"

	sessionID, err := service.SendOTP(ctx, phone)
	if err != nil {
		t.Fatalf("SendOTP failed: %v", err)
	}

	if sessionID == "" {
		t.Error("Expected non-empty session ID")
	}

	// Verify UUID format
	if _, err := uuid.Parse(sessionID); err != nil {
		t.Errorf("Session ID is not a valid UUID: %v", err)
	}
}

func TestSendOTP_InvalidPhone(t *testing.T) {
	service, _, cleanup := setupTestAuthService(t)
	defer cleanup()

	ctx := context.Background()

	testCases := []struct {
		name  string
		phone string
	}{
		{"missing prefix", "7071234567"},
		{"wrong prefix", "+87071234567"},
		{"too short", "+770712345"},
		{"too long", "+770712345678"},
		{"non-numeric", "+7707abc4567"},
		{"empty", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.SendOTP(ctx, tc.phone)
			if err != ErrInvalidPhoneFormat {
				t.Errorf("Expected ErrInvalidPhoneFormat, got: %v", err)
			}
		})
	}
}

func TestSendOTP_RateLimit(t *testing.T) {
	service, _, cleanup := setupTestAuthService(t)
	defer cleanup()

	ctx := context.Background()
	phone := "+77071234567"

	// Send 3 OTPs (should succeed)
	for i := 0; i < 3; i++ {
		_, err := service.SendOTP(ctx, phone)
		if err != nil {
			t.Fatalf("OTP %d failed: %v", i+1, err)
		}
	}

	// 4th should fail
	_, err := service.SendOTP(ctx, phone)
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Code != ErrSMSRateLimited.Code {
		t.Errorf("Expected ErrSMSRateLimited, got: %v", err)
	}
}

func TestGenerateOTP_DevMode(t *testing.T) {
	service, _, cleanup := setupTestAuthService(t)
	defer cleanup()

	// In dev mode, should always return 1234
	code := service.generateOTP()
	if code != "1234" {
		t.Errorf("Expected dev OTP '1234', got: %s", code)
	}
}

func TestVerifyOTP_CorrectCode(t *testing.T) {
	_, rdb, cleanup := setupTestAuthService(t)
	defer cleanup()

	ctx := context.Background()

	// Create mock session
	sessionID := uuid.New().String()

	// Store in Redis manually
	key := "otp:" + sessionID
	sessionJSON := `{"phone":"+77071234567","code":"1234","attempts":0}`
	rdb.Set(ctx, key, sessionJSON, 5*time.Minute)

	// Verify session exists before test
	exists, _ := rdb.Exists(ctx, key).Result()
	if exists == 0 {
		t.Fatal("Session should exist before verification")
	}

	// Note: VerifyOTP will panic because repo is nil (no DB connection in unit tests).
	// This test verifies that the OTP session logic (code matching + session deletion)
	// works correctly. The repo call happens AFTER session deletion, so we can verify
	// Redis state was cleaned up by recovering from the panic.
	// In a real test, you'd use a test database or mock the repository.
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic due to nil repo, but no panic occurred")
			}
		}()
		service, _, _ := setupTestAuthService(t)
		service.VerifyOTP(ctx, sessionID, "1234")
	}()

	// Session should be deleted after correct code (deletion happens before repo call)
	exists, _ = rdb.Exists(ctx, key).Result()
	if exists != 0 {
		t.Error("Session should be deleted after correct code")
	}
}

func TestVerifyOTP_WrongCode(t *testing.T) {
	service, rdb, cleanup := setupTestAuthService(t)
	defer cleanup()

	ctx := context.Background()

	// Create mock session
	sessionID := uuid.New().String()
	key := "otp:" + sessionID
	sessionJSON := `{"phone":"+77071234567","code":"1234","attempts":0}`
	rdb.Set(ctx, key, sessionJSON, 5*time.Minute)

	// Try wrong code
	_, err := service.VerifyOTP(ctx, sessionID, "9999")
	if err != ErrOTPInvalidCode {
		t.Errorf("Expected ErrOTPInvalidCode, got: %v", err)
	}

	// Session should still exist with incremented attempts
	exists, _ := rdb.Exists(ctx, key).Result()
	if exists == 0 {
		t.Error("Session should still exist after wrong code")
	}
}

func TestVerifyOTP_MaxAttempts(t *testing.T) {
	service, rdb, cleanup := setupTestAuthService(t)
	defer cleanup()

	ctx := context.Background()

	// Create session with 5 attempts already
	sessionID := uuid.New().String()
	key := "otp:" + sessionID
	sessionJSON := `{"phone":"+77071234567","code":"1234","attempts":5}`
	rdb.Set(ctx, key, sessionJSON, 5*time.Minute)

	// Try to verify
	_, err := service.VerifyOTP(ctx, sessionID, "1234")
	if err != ErrOTPMaxAttempts {
		t.Errorf("Expected ErrOTPMaxAttempts, got: %v", err)
	}

	// Session should be deleted
	exists, _ := rdb.Exists(ctx, key).Result()
	if exists != 0 {
		t.Error("Session should be deleted after max attempts")
	}
}

func TestVerifyOTP_ExpiredSession(t *testing.T) {
	service, _, cleanup := setupTestAuthService(t)
	defer cleanup()

	ctx := context.Background()

	// Try with non-existent session
	_, err := service.VerifyOTP(ctx, uuid.New().String(), "1234")
	if err != ErrOTPSessionExpired {
		t.Errorf("Expected ErrOTPSessionExpired, got: %v", err)
	}
}

func TestMaskPhone(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"+77071234567", "+7707***4567"},
		{"+77001234567", "+7700***4567"},
		{"short", "short"},
		{"", ""},
	}

	for _, tc := range testCases {
		result := maskPhone(tc.input)
		if result != tc.expected {
			t.Errorf("maskPhone(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}
