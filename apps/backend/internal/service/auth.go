package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"regexp"
	"time"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

var phoneRegex = regexp.MustCompile(`^\+7[0-9]{10}$`)

// AuthService handles authentication logic
type AuthService struct {
	repo         *repository.Queries
	redis        *redis.Client
	tokenService *TokenService
	environment  string
}

// OTPSession represents OTP session data stored in Redis
type OTPSession struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Attempts int    `json:"attempts"`
}

// AuthResult represents authentication result
type AuthResult struct {
	IsNewUser    bool        `json:"is_new_user"`
	AccessToken  string      `json:"access_token,omitempty"`
	RefreshToken string      `json:"refresh_token,omitempty"`
	TempToken    string      `json:"temp_token,omitempty"`
	User         interface{} `json:"user"`
}

// NewAuthService creates a new AuthService
func NewAuthService(
	repo *repository.Queries,
	redis *redis.Client,
	tokenService *TokenService,
	environment string,
) *AuthService {
	return &AuthService{
		repo:         repo,
		redis:        redis,
		tokenService: tokenService,
		environment:  environment,
	}
}

// SendOTP sends OTP code to phone number
func (s *AuthService) SendOTP(ctx context.Context, phone string) (string, error) {
	// Validate phone format
	if !phoneRegex.MatchString(phone) {
		return "", ErrInvalidPhoneFormat
	}

	// Check rate limits
	if err := s.checkSMSRateLimit(ctx, phone); err != nil {
		return "", err
	}

	// Generate OTP code
	code := s.generateOTP()

	// Create session
	sessionID := uuid.New().String()
	session := OTPSession{
		Phone:    phone,
		Code:     code,
		Attempts: 0,
	}

	// Store session in Redis with 5 min TTL
	sessionData, err := json.Marshal(session)
	if err != nil {
		return "", fmt.Errorf("marshal session: %w", err)
	}

	key := fmt.Sprintf("otp:%s", sessionID)
	if err := s.redis.Set(ctx, key, sessionData, 5*time.Minute).Err(); err != nil {
		return "", fmt.Errorf("store OTP session: %w", err)
	}

	// Send SMS (mock in dev)
	if err := s.sendSMS(phone, code); err != nil {
		slog.Error("failed to send SMS", "phone", maskPhone(phone), "error", err)
		// Don't fail the request if SMS fails, session is still created
	}

	slog.Info("OTP sent", "phone", maskPhone(phone), "session_id", sessionID, "code", code)

	return sessionID, nil
}

// VerifyOTP verifies OTP code and returns auth result
func (s *AuthService) VerifyOTP(ctx context.Context, sessionID, code string) (*AuthResult, error) {
	// Load session from Redis
	key := fmt.Sprintf("otp:%s", sessionID)
	data, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrOTPSessionExpired
	}
	if err != nil {
		return nil, fmt.Errorf("get OTP session: %w", err)
	}

	var session OTPSession
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}

	// Check attempts
	if session.Attempts >= 5 {
		// Delete session
		s.redis.Del(ctx, key)
		return nil, ErrOTPMaxAttempts
	}

	// Verify code
	if session.Code != code {
		// Increment attempts
		session.Attempts++
		sessionData, _ := json.Marshal(session)
		s.redis.Set(ctx, key, sessionData, 5*time.Minute)
		return nil, ErrOTPInvalidCode
	}

	// Code is correct, delete session
	s.redis.Del(ctx, key)

	// Find or create user
	user, err := s.repo.GetUserByPhone(ctx, session.Phone)
	if err == pgx.ErrNoRows {
		// New user - create record
		newUser, err := s.repo.CreateUser(ctx, session.Phone)
		if err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}

		// Convert pgtype.UUID to uuid.UUID
		userUUID, err := uuid.FromBytes(newUser.ID.Bytes[:])
		if err != nil {
			return nil, fmt.Errorf("parse user UUID: %w", err)
		}

		// Generate temp token for profile setup
		tempToken, err := s.tokenService.GenerateAccessToken(userUUID, "user")
		if err != nil {
			return nil, fmt.Errorf("generate temp token: %w", err)
		}

		return &AuthResult{
			IsNewUser: true,
			TempToken: tempToken,
			User: map[string]interface{}{
				"id":                   userUUID.String(),
				"phone":                maskPhone(newUser.Phone),
				"is_profile_complete":  false,
			},
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by phone: %w", err)
	}

	// Convert pgtype.UUID to uuid.UUID
	userUUID, err := uuid.FromBytes(user.ID.Bytes[:])
	if err != nil {
		return nil, fmt.Errorf("parse user UUID: %w", err)
	}

	// Existing user - generate tokens
	accessToken, err := s.tokenService.GenerateAccessToken(userUUID, "user")
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &AuthResult{
		IsNewUser:    false,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: map[string]interface{}{
			"id":                   userUUID.String(),
			"phone":                maskPhone(user.Phone),
			"first_name":           user.FirstName.String,
			"last_name":            user.LastName.String,
			"is_profile_complete":  user.FirstName.Valid && user.LastName.Valid,
		},
	}, nil
}

// checkSMSRateLimit checks if phone number exceeded SMS rate limits
func (s *AuthService) checkSMSRateLimit(ctx context.Context, phone string) error {
	// Check hourly limit (3 per hour)
	hourKey := fmt.Sprintf("sms_rate:%s:hour", phone)
	hourCount, err := s.redis.Incr(ctx, hourKey).Result()
	if err != nil {
		return fmt.Errorf("increment hour counter: %w", err)
	}
	if hourCount == 1 {
		s.redis.Expire(ctx, hourKey, time.Hour)
	}
	if hourCount > 3 {
		return ErrSMSRateLimited.WithMessage("Too many SMS requests. Try again in 1 hour")
	}

	// Check daily limit (10 per day)
	dayKey := fmt.Sprintf("sms_rate:%s:day", phone)
	dayCount, err := s.redis.Incr(ctx, dayKey).Result()
	if err != nil {
		return fmt.Errorf("increment day counter: %w", err)
	}
	if dayCount == 1 {
		s.redis.Expire(ctx, dayKey, 24*time.Hour)
	}
	if dayCount > 10 {
		return ErrSMSRateLimited.WithMessage("Too many SMS requests. Try again tomorrow")
	}

	return nil
}

// generateOTP generates a 4-digit OTP code
func (s *AuthService) generateOTP() string {
	// In dev mode, always return 1234
	if s.environment == "development" {
		return "1234"
	}

	// Generate random 4-digit code
	n, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		slog.Error("failed to generate random OTP", "error", err)
		return "1234" // fallback
	}

	return fmt.Sprintf("%04d", n.Int64())
}

// sendSMS sends SMS via provider (mocked in dev)
func (s *AuthService) sendSMS(phone, code string) error {
	if s.environment == "development" {
		slog.Info("SMS mock", "phone", maskPhone(phone), "code", code)
		return nil
	}

	// TODO: Integrate with real SMS provider (SMSC.kz)
	// For now, just log
	slog.Warn("SMS sending not implemented", "phone", maskPhone(phone))
	return nil
}

// ProfileSetupInput represents profile setup data
type ProfileSetupInput struct {
	FirstName string
	LastName  string
	Gender    string
	BirthYear int16
	City      string
	District  string
	Language  string
}

// ProfileSetup completes user profile and upgrades temp token to full tokens
func (s *AuthService) ProfileSetup(ctx context.Context, userID uuid.UUID, input ProfileSetupInput) (*AuthResult, error) {
	// Check if profile already complete
	user, err := s.repo.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if user.FirstName.Valid && user.LastName.Valid {
		return nil, ErrProfileAlreadySet
	}

	// Convert gender string to pgtype
	var genderType repository.NullGenderType
	if input.Gender == "male" || input.Gender == "female" {
		genderType = repository.NullGenderType{
			GenderType: repository.GenderType(input.Gender),
			Valid:      true,
		}
	}

	// Update user profile
	updatedUser, err := s.repo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:                pgtype.UUID{Bytes: userID, Valid: true},
		FirstName:         pgtype.Text{String: input.FirstName, Valid: true},
		LastName:          pgtype.Text{String: input.LastName, Valid: true},
		Gender:            genderType,
		BirthYear:         pgtype.Int2{Int16: input.BirthYear, Valid: true},
		City:              pgtype.Text{String: input.City, Valid: true},
		District:          pgtype.Text{String: input.District, Valid: true},
		Language:          pgtype.Text{String: input.Language, Valid: true},
		IsProfileComplete: pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	// Generate full access and refresh tokens (upgrade from temp token)
	accessToken, err := s.tokenService.GenerateAccessToken(userID, "user")
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &AuthResult{
		IsNewUser:    false,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: map[string]interface{}{
			"id":                   userID.String(),
			"phone":                maskPhone(updatedUser.Phone),
			"first_name":           updatedUser.FirstName.String,
			"last_name":            updatedUser.LastName.String,
			"gender":               input.Gender,
			"birth_year":           updatedUser.BirthYear.Int16,
			"city":                 updatedUser.City.String,
			"district":             updatedUser.District.String,
			"language":             updatedUser.Language.String,
			"is_profile_complete":  true,
		},
	}, nil
}

// maskPhone masks phone number for logging
func maskPhone(phone string) string {
	if len(phone) < 8 {
		return phone
	}
	return phone[:5] + "***" + phone[len(phone)-4:]
}
