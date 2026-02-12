package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/service"
	"github.com/google/uuid"
)

// ChatHandler handles chat REST endpoints
type ChatHandler struct {
	chatService *service.ChatService
}

// NewChatHandler creates a new ChatHandler
func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

// CreatePersonalChat handles POST /v1/chats/personal
func (h *ChatHandler) CreatePersonalChat(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var req struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if req.UserID == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "user_id is required")
		return
	}

	otherUserID, err := parseUUIDString(req.UserID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid user_id")
		return
	}

	chatID, isNew, err := h.chatService.CreatePersonalChat(r.Context(), userID, otherUserID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"chat_id": chatID,
		"is_new":  isNew,
	})
}

// ListChats handles GET /v1/chats
func (h *ChatHandler) ListChats(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	chats, err := h.chatService.ListMyChats(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, chats)
}

// GetMessages handles GET /v1/chats/:id/messages
func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	chatID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid chat ID")
		return
	}

	q := r.URL.Query()
	before := q.Get("before")
	limit := queryInt(q.Get("limit"), 50)

	messages, hasMore, err := h.chatService.GetMessages(r.Context(), userID, chatID, before, limit)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Shape aligns with mobile client: { "data": { "data": [...], "has_more": bool } }
	respondJSON(w, http.StatusOK, map[string]any{
		"data":     messages,
		"has_more": hasMore,
	})
}

// SendMessage handles POST /v1/chats/:id/messages (REST fallback)
func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	chatID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid chat ID")
		return
	}

	var req struct {
		Content   string  `json:"content"`
		ReplyToID *string `json:"reply_to_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if req.Content == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "content is required")
		return
	}

	msg, err := h.chatService.SendMessage(r.Context(), userID, chatID, req.Content, req.ReplyToID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, msg)
}

// MarkAsRead handles POST /v1/chats/:id/read
func (h *ChatHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	chatID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid chat ID")
		return
	}

	if err := h.chatService.MarkAsRead(r.Context(), userID, chatID); err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// UpdateMuted handles PATCH /v1/chats/:id/mute
func (h *ChatHandler) UpdateMuted(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	chatID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid chat ID")
		return
	}

	var req struct {
		IsMuted bool `json:"is_muted"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if err := h.chatService.UpdateMuted(r.Context(), userID, chatID, req.IsMuted); err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// GetUnreadCount handles GET /v1/chats/unread-count
func (h *ChatHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	count, err := h.chatService.GetUnreadCount(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"total_unread": count,
	})
}

// parseUUIDString parses a UUID string
func parseUUIDString(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
