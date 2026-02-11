package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiter creates a rate limiting middleware
func RateLimiter(redis *redis.Client, limit int, window time.Duration, keyFunc func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			key := keyFunc(r)
			if key == "" {
				// Skip rate limiting if no key
				next.ServeHTTP(w, r)
				return
			}

			rateLimitKey := fmt.Sprintf("rate:%s", key)

			// Increment counter
			count, err := redis.Incr(ctx, rateLimitKey).Result()
			if err != nil {
				// If Redis fails, allow the request (fail open)
				next.ServeHTTP(w, r)
				return
			}

			// Set expiry on first request
			if count == 1 {
				redis.Expire(ctx, rateLimitKey, window)
			}

			// Get TTL for reset time
			ttl, _ := redis.TTL(ctx, rateLimitKey).Result()
			resetTime := time.Now().Add(ttl)

			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(max(0, limit-int(count))))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

			// Check if limit exceeded
			if count > int64(limit) {
				retryAfter := int(ttl.Seconds())
				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				respondError(w, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// IPKeyFunc returns a key function that uses client IP
func IPKeyFunc(prefix string) func(*http.Request) string {
	return func(r *http.Request) string {
		ip := getClientIP(r)
		return fmt.Sprintf("%s:%s", prefix, ip)
	}
}

// UserKeyFunc returns a key function that uses authenticated user ID
func UserKeyFunc(prefix string) func(*http.Request) string {
	return func(r *http.Request) string {
		userID := GetUserID(r.Context())
		if userID == "" {
			// Fall back to IP if not authenticated
			return IPKeyFunc(prefix)(r)
		}
		return fmt.Sprintf("%s:%s", prefix, userID)
	}
}

// getClientIP extracts client IP from request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (proxy/load balancer)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
