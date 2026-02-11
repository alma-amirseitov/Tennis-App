package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/service"
)

// Context keys for user data
const (
	userIDKey   contextKey = "user_id"
	userRoleKey contextKey = "user_role"
)

// Auth validates JWT token and injects user data into context
func Auth(tokenService *service.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization header")
				return
			}

			// Check Bearer prefix
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization header format")
				return
			}

			tokenString := parts[1]

			// Validate token
			claims, err := tokenService.ValidateAccessToken(tokenString)
			if err != nil {
				if err == service.ErrTokenInvalid {
					respondError(w, http.StatusUnauthorized, "TOKEN_INVALID", "Invalid token")
				} else if err == service.ErrTokenExpired {
					respondError(w, http.StatusUnauthorized, "TOKEN_EXPIRED", "Token expired")
				} else {
					respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication failed")
				}
				return
			}

			// Inject user data into context
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, userRoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves user ID from context
func GetUserID(ctx context.Context) string {
	userID, _ := ctx.Value(userIDKey).(string)
	return userID
}

// GetUserRole retrieves user role from context
func GetUserRole(ctx context.Context) string {
	role, _ := ctx.Value(userRoleKey).(string)
	return role
}

// respondError sends an error response (helper for middleware)
func respondError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error":{"code":"` + code + `","message":"` + message + `"}}`))
}
