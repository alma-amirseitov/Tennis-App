package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/repository"
)

const communityRoleKey contextKey = "community_role"

// RequireCommunityRole creates middleware that checks if the user has one of the required roles
// in the community specified by the :id URL parameter
func RequireCommunityRole(queries *repository.Queries, roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDStr := GetUserID(r.Context())
			if userIDStr == "" {
				respondUnauthorized(w)
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				respondUnauthorized(w)
				return
			}

			communityIDStr := chi.URLParam(r, "id")
			communityID, err := uuid.Parse(communityIDStr)
			if err != nil {
				respondJSON(w, http.StatusBadRequest, "INVALID_ID", "Invalid community ID")
				return
			}

			member, err := queries.GetCommunityMember(r.Context(), repository.GetCommunityMemberParams{
				CommunityID: pgtype.UUID{Bytes: communityID, Valid: true},
				UserID:      pgtype.UUID{Bytes: userID, Valid: true},
			})
			if err != nil {
				respondJSON(w, http.StatusForbidden, "NOT_COMMUNITY_MEMBER", "You are not a member of this community")
				return
			}

			if member.Status.MemberStatus != repository.MemberStatusActive {
				respondJSON(w, http.StatusForbidden, "NOT_ACTIVE_MEMBER", "Your membership is not active")
				return
			}

			// Check if user has one of the required roles
			hasRole := false
			for _, requiredRole := range roles {
				if string(member.Role.CommunityRole) == requiredRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				respondJSON(w, http.StatusForbidden, "INSUFFICIENT_ROLE", "You don't have the required role")
				return
			}

			// Set community role in context
			ctx := context.WithValue(r.Context(), communityRoleKey, string(member.Role.CommunityRole))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetCommunityRole returns the community role from context
func GetCommunityRole(ctx context.Context) string {
	role, _ := ctx.Value(communityRoleKey).(string)
	return role
}

func respondUnauthorized(w http.ResponseWriter) {
	respondJSON(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
}

func respondJSON(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error":{"code":"` + code + `","message":"` + message + `"}}`))
}
