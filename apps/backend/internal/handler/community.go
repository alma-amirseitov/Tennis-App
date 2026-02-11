package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/service"
)

// CommunityHandler handles community endpoints
type CommunityHandler struct {
	communityService *service.CommunityService
}

// NewCommunityHandler creates a new CommunityHandler
func NewCommunityHandler(communityService *service.CommunityService) *CommunityHandler {
	return &CommunityHandler{communityService: communityService}
}

// Create handles POST /v1/communities
func (h *CommunityHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var input service.CreateCommunityInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if input.Name == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name is required")
		return
	}
	if input.CommunityType == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "community_type is required")
		return
	}

	community, err := h.communityService.Create(r.Context(), userID, input)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, community)
}

// List handles GET /v1/communities
func (h *CommunityHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	input := service.ListCommunitiesInput{
		Type:        q.Get("type"),
		AccessLevel: q.Get("access_level"),
		District:    q.Get("district"),
		Query:       q.Get("q"),
		Sort:        q.Get("sort"),
		Page:        queryInt(q.Get("page"), 1),
		PerPage:     queryInt(q.Get("per_page"), 20),
	}

	if q.Get("verified_only") == "true" {
		input.VerifiedOnly = true
	}

	communities, pagination, err := h.communityService.List(r.Context(), input)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondPaginated(w, http.StatusOK, communities, *pagination)
}

// GetByID handles GET /v1/communities/:id
func (h *CommunityHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	communityID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid community ID")
		return
	}

	community, err := h.communityService.GetByID(r.Context(), userID, communityID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, community)
}

// Join handles POST /v1/communities/:id/join
func (h *CommunityHandler) Join(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	communityID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid community ID")
		return
	}

	var body struct {
		Message string `json:"message"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	result, err := h.communityService.Join(r.Context(), userID, communityID, body.Message)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// Leave handles POST /v1/communities/:id/leave
func (h *CommunityHandler) Leave(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	communityID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid community ID")
		return
	}

	if err := h.communityService.Leave(r.Context(), userID, communityID); err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Successfully left community"})
}

// ListMembers handles GET /v1/communities/:id/members
func (h *CommunityHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	communityID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid community ID")
		return
	}

	q := r.URL.Query()
	input := service.ListMembersInput{
		CommunityID: communityID,
		Role:        q.Get("role"),
		Status:      q.Get("status"),
		Query:       q.Get("q"),
		Sort:        q.Get("sort"),
		Page:        queryInt(q.Get("page"), 1),
		PerPage:     queryInt(q.Get("per_page"), 20),
	}

	members, pagination, err := h.communityService.ListMembers(r.Context(), input)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondPaginated(w, http.StatusOK, members, *pagination)
}

// UpdateMemberRole handles PATCH /v1/communities/:id/members/:userId
func (h *CommunityHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	actorID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	communityID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid community ID")
		return
	}

	targetUserIDStr := chi.URLParam(r, "userId")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
		return
	}

	var body struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if body.Role == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "role is required")
		return
	}

	if err := h.communityService.UpdateMemberRole(r.Context(), actorID, communityID, targetUserID, body.Role); err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Role updated"})
}

// ReviewRequest handles POST /v1/communities/:id/members/:userId/review
func (h *CommunityHandler) ReviewRequest(w http.ResponseWriter, r *http.Request) {
	actorID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	communityID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid community ID")
		return
	}

	targetUserIDStr := chi.URLParam(r, "userId")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
		return
	}

	var body struct {
		Action string `json:"action"` // "approve" or "reject"
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	approve := body.Action == "approve"

	if err := h.communityService.ReviewRequest(r.Context(), actorID, communityID, targetUserID, approve); err != nil {
		handleServiceError(w, err)
		return
	}

	msg := "Request rejected"
	if approve {
		msg = "Request approved"
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": msg})
}

// ListMyCommunities handles GET /v1/communities/my
func (h *CommunityHandler) ListMyCommunities(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	communities, err := h.communityService.ListMyCommunities(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, communities)
}

// parseUUIDParam extracts and parses a UUID URL parameter
func parseUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	s := chi.URLParam(r, param)
	return uuid.Parse(s)
}

// parseQueryFloat parses a float query parameter
func parseQueryFloat(s string) *float64 {
	if s == "" {
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &f
}
