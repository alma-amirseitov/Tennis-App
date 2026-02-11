package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// TokenService handles JWT token generation and validation
type TokenService struct {
	secret          string
	accessTTL       time.Duration
	refreshTTL      time.Duration
	redis           *redis.Client
}

// TokenClaims represents JWT claims
type TokenClaims struct {
	UserID string `json:"sub"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// RefreshTokenClaims represents refresh token claims
type RefreshTokenClaims struct {
	UserID string `json:"sub"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

// NewTokenService creates a new TokenService
func NewTokenService(secret string, accessTTL, refreshTTL time.Duration, redis *redis.Client) *TokenService {
	return &TokenService{
		secret:     secret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		redis:      redis,
	}
}

// GenerateAccessToken generates a new access token
func (s *TokenService) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID: userID.String(),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// GenerateRefreshToken generates a new refresh token and stores JTI in Redis
func (s *TokenService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	now := time.Now()
	jti := uuid.New().String()

	claims := RefreshTokenClaims{
		UserID: userID.String(),
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("sign refresh token: %w", err)
	}

	// Store JTI in Redis with TTL
	key := fmt.Sprintf("refresh:%s", jti)
	if err := s.redis.Set(ctx, key, userID.String(), s.refreshTTL).Err(); err != nil {
		return "", fmt.Errorf("store refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateAccessToken validates an access token and returns claims
func (s *TokenService) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// ValidateRefreshToken validates refresh token and checks JTI in Redis
func (s *TokenService) ValidateRefreshToken(ctx context.Context, tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	// Check if JTI exists in Redis (not revoked)
	key := fmt.Sprintf("refresh:%s", claims.JTI)
	exists, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("check refresh token: %w", err)
	}
	if exists == 0 {
		return nil, ErrTokenRevoked
	}

	return claims, nil
}

// RevokeRefreshToken removes JTI from Redis (one-time use)
func (s *TokenService) RevokeRefreshToken(ctx context.Context, jti string) error {
	key := fmt.Sprintf("refresh:%s", jti)
	if err := s.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

// RevokeAllUserTokens removes all refresh tokens for a user
func (s *TokenService) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	// Scan for all refresh tokens belonging to this user
	// Note: This is a simplified version. In production, you might want to maintain
	// a separate set of JTIs per user for faster revocation.
	pattern := "refresh:*"
	iter := s.redis.Scan(ctx, 0, pattern, 0).Iterator()

	var keysToDelete []string
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := s.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		if val == userID.String() {
			keysToDelete = append(keysToDelete, key)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("scan refresh tokens: %w", err)
	}

	if len(keysToDelete) > 0 {
		if err := s.redis.Del(ctx, keysToDelete...).Err(); err != nil {
			return fmt.Errorf("delete user tokens: %w", err)
		}
	}

	return nil
}

// RefreshTokens validates refresh token and issues new token pair
func (s *TokenService) RefreshTokens(ctx context.Context, refreshTokenString string) (string, string, error) {
	// Validate refresh token
	claims, err := s.ValidateRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return "", "", err
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return "", "", fmt.Errorf("parse user ID: %w", err)
	}

	// Revoke old refresh token (one-time use)
	if err := s.RevokeRefreshToken(ctx, claims.JTI); err != nil {
		return "", "", fmt.Errorf("revoke old token: %w", err)
	}

	// Generate new token pair
	accessToken, err := s.GenerateAccessToken(userID, "user")
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	newRefreshToken, err := s.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	return accessToken, newRefreshToken, nil
}
