package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/handler/middleware"
	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/service"
)

// UserHandler handles user endpoints
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetMe handles GET /v1/users/me
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	profile, err := h.userService.GetProfile(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, profile)
}

// UpdateMe handles PATCH /v1/users/me
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var input service.UpdateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	profile, err := h.userService.UpdateProfile(r.Context(), userID, input)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, profile)
}

// UploadAvatar handles POST /v1/users/me/avatar
func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Parse multipart form (max 6 MB to account for overhead)
	if err := r.ParseMultipartForm(6 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Could not parse multipart form")
		return
	}

	file, _, err := r.FormFile("avatar")
	if err != nil {
		respondError(w, http.StatusBadRequest, "MISSING_FILE", "Avatar file is required")
		return
	}
	defer file.Close()

	// Read file data
	fileData, err := io.ReadAll(io.LimitReader(file, 6<<20))
	if err != nil {
		respondError(w, http.StatusBadRequest, "READ_ERROR", "Could not read file")
		return
	}

	avatarURL, err := h.userService.UploadAvatar(r.Context(), userID, fileData)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"avatar_url": avatarURL})
}

// GetUser handles GET /v1/users/:id
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	currentUserID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	targetIDStr := chi.URLParam(r, "id")
	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
		return
	}

	profile, err := h.userService.GetPublicProfile(r.Context(), currentUserID, targetID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, profile)
}

// SearchUsers handles GET /v1/users/search
func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	input := service.SearchUsersInput{
		Query:    q.Get("q"),
		Gender:   q.Get("gender"),
		District: q.Get("district"),
		Sort:     q.Get("sort"),
		Page:     queryInt(q.Get("page"), 1),
		PerPage:  queryInt(q.Get("per_page"), 20),
	}

	if v := q.Get("min_level"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			input.MinLevel = &f
		}
	}
	if v := q.Get("max_level"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			input.MaxLevel = &f
		}
	}

	result, err := h.userService.SearchUsers(r.Context(), input)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondPaginated(w, http.StatusOK, result.Users, result.Pagination)
}

// getUserUUID extracts and parses user ID from context
func getUserUUID(r *http.Request) (uuid.UUID, error) {
	userIDStr := middleware.GetUserID(r.Context())
	if userIDStr == "" {
		return uuid.Nil, ErrNoUserID
	}
	return uuid.Parse(userIDStr)
}

func queryInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 {
		return defaultVal
	}
	return v
}

var ErrNoUserID = &struct{ error }{error: nil}

// respondPaginated sends a paginated JSON response
func respondPaginated(w http.ResponseWriter, status int, data any, pagination service.PaginationInfo) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"data":       data,
		"pagination": pagination,
	})
}
