package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/service"
)

// NotificationHandler handles notification-related endpoints.
type NotificationHandler struct {
	notificationService *service.NotificationService
}

// NewNotificationHandler creates a new NotificationHandler.
func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// List handles GET /v1/notifications
func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	q := r.URL.Query()
	page := queryInt(q.Get("page"), 1)
	perPage := queryInt(q.Get("per_page"), 20)

	items, pagination, err := h.notificationService.List(r.Context(), userID, page, perPage)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondPaginated(w, http.StatusOK, items, *pagination)
}

// MarkRead handles POST /v1/notifications/read
// Body: { "ids": ["uuid1", "uuid2"] } or { "read_all": true }
func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var body struct {
		IDs     []string `json:"ids"`
		ReadAll bool     `json:"read_all"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if body.ReadAll {
		if err := h.notificationService.MarkAllAsRead(r.Context(), userID); err != nil {
			handleServiceError(w, err)
			return
		}
		respondJSON(w, http.StatusOK, map[string]any{"status": "ok"})
		return
	}

	if len(body.IDs) == 0 {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "ids or read_all is required")
		return
	}

	for _, idStr := range body.IDs {
		id, err := parseUUIDString(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid notification ID")
			return
		}
		if err := h.notificationService.MarkAsRead(r.Context(), userID, id); err != nil {
			handleServiceError(w, err)
			return
		}
	}

	respondJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

// GetUnreadCount handles GET /v1/notifications/unread-count
func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	count, err := h.notificationService.GetUnreadCount(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"count": count,
	})
}

// Delete handles DELETE /v1/notifications/{id}
func (h *NotificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	notificationID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid notification ID")
		return
	}

	if err := h.notificationService.Delete(r.Context(), userID, notificationID); err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
	})
}
