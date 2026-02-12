package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/service"
)

// MatchHandler handles match endpoints
type MatchHandler struct {
	matchService *service.MatchService
}

// NewMatchHandler creates a new MatchHandler
func NewMatchHandler(matchService *service.MatchService) *MatchHandler {
	return &MatchHandler{matchService: matchService}
}

// GetByID handles GET /v1/matches/{id}
func (h *MatchHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	matchID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid match ID")
		return
	}

	match, err := h.matchService.GetByID(r.Context(), matchID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, match)
}

// ListMyMatches handles GET /v1/matches/my
func (h *MatchHandler) ListMyMatches(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	q := r.URL.Query()
	input := service.ListMyMatchesInput{
		CommunityID: q.Get("community_id"),
		OpponentID:  q.Get("opponent_id"),
		Result:      q.Get("result"),
		Page:        queryInt(q.Get("page"), 1),
		PerPage:     queryInt(q.Get("per_page"), 20),
	}

	matches, total, err := h.matchService.ListMyMatches(r.Context(), userID, input)
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

	pagination := service.PaginationInfo{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: totalPages,
	}

	respondPaginated(w, http.StatusOK, matches, pagination)
}

// SubmitResult handles POST /v1/matches/{id}/result
func (h *MatchHandler) SubmitResult(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	matchID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid match ID")
		return
	}

	var input service.SubmitResultInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if input.WinnerID == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "winner_id is required")
		return
	}
	if len(input.Score) == 0 {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "score is required")
		return
	}

	match, err := h.matchService.SubmitResult(r.Context(), userID, matchID, input)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"match_id":      match.ID,
		"result_status": match.ResultStatus,
		"message":       "Result submitted. Waiting for opponent confirmation.",
	})
}

// ConfirmResult handles POST /v1/matches/{id}/confirm
func (h *MatchHandler) ConfirmResult(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	matchID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid match ID")
		return
	}

	var input service.ConfirmResultInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if input.Action == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "action is required (confirm or dispute)")
		return
	}

	match, err := h.matchService.ConfirmResult(r.Context(), userID, matchID, input)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, match)
}

// AdminConfirm handles POST /v1/matches/{id}/admin-confirm
func (h *MatchHandler) AdminConfirm(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	matchID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid match ID")
		return
	}

	var input service.AdminConfirmInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if input.WinnerID == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "winner_id is required")
		return
	}

	match, err := h.matchService.AdminConfirm(r.Context(), userID, matchID, input)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, match)
}
