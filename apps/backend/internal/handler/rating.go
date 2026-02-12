package handler

import (
	"net/http"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/service"
)

// RatingHandler handles rating and leaderboard endpoints
type RatingHandler struct {
	ratingService *service.RatingService
}

// NewRatingHandler creates a new RatingHandler
func NewRatingHandler(ratingService *service.RatingService) *RatingHandler {
	return &RatingHandler{ratingService: ratingService}
}

// GetGlobalLeaderboard handles GET /v1/rating/global
func (h *RatingHandler) GetGlobalLeaderboard(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	input := service.ListLeaderboardInput{
		MinGames: queryInt(q.Get("min_games"), 0),
		Page:     queryInt(q.Get("page"), 1),
		PerPage:  queryInt(q.Get("per_page"), 20),
	}

	entries, total, err := h.ratingService.GetGlobalLeaderboard(r.Context(), input)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	page := input.Page
	if page < 1 {
		page = 1
	}
	perPage := input.PerPage
	if perPage < 1 || perPage > 50 {
		perPage = 20
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"leaderboard": entries,
		"pagination": map[string]any{
			"page":        page,
			"per_page":    perPage,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetCommunityLeaderboard handles GET /v1/rating/community/{id}
func (h *RatingHandler) GetCommunityLeaderboard(w http.ResponseWriter, r *http.Request) {
	communityID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid community ID")
		return
	}

	q := r.URL.Query()
	input := service.ListLeaderboardInput{
		MinGames: queryInt(q.Get("min_games"), 0),
		Page:     queryInt(q.Get("page"), 1),
		PerPage:  queryInt(q.Get("per_page"), 20),
	}

	entries, err := h.ratingService.GetCommunityLeaderboard(r.Context(), communityID, input)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"leaderboard": entries,
	})
}

// GetMyRating handles GET /v1/rating/me
func (h *RatingHandler) GetMyRating(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	rating, err := h.ratingService.GetMyRating(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, rating)
}

// GetRatingHistory handles GET /v1/rating/history
func (h *RatingHandler) GetRatingHistory(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	q := r.URL.Query()
	input := service.ListRatingHistoryInput{
		CommunityID: q.Get("community_id"),
		Period:      q.Get("period"),
		Page:        queryInt(q.Get("page"), 1),
		PerPage:     queryInt(q.Get("per_page"), 50),
	}

	history, err := h.ratingService.GetMyRatingHistory(r.Context(), userID, input)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"history": history,
	})
}

// GetMyStats handles GET /v1/rating/stats
func (h *RatingHandler) GetMyStats(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	stats, err := h.ratingService.GetMyStats(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, stats)
}
