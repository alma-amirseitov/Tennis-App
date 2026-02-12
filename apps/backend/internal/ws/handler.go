package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"nhooyr.io/websocket"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/service"
)

const (
	// Rate limit: 60 messages per minute per user
	msgRateLimit    = 60
	msgRateWindow   = time.Minute
	writeTimeout    = 10 * time.Second
	pongWait        = 60 * time.Second
	pingInterval    = 30 * time.Second
	maxMessageSize  = 4096
)

// Handler handles WebSocket upgrade and message routing
type Handler struct {
	hub          *Hub
	chatService  *service.ChatService
	tokenService *service.TokenService
	redis        *goredis.Client
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub, chatService *service.ChatService, tokenService *service.TokenService, redis *goredis.Client) *Handler {
	return &Handler{
		hub:          hub,
		chatService:  chatService,
		tokenService: tokenService,
		redis:        redis,
	}
}

// ServeHTTP handles the WebSocket upgrade request
// GET /v1/ws?token={jwt}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract and validate JWT from query param
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"Missing token"}}`, http.StatusUnauthorized)
		return
	}

	claims, err := h.tokenService.ValidateAccessToken(token)
	if err != nil {
		http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"Invalid token"}}`, http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"Invalid user ID"}}`, http.StatusUnauthorized)
		return
	}

	// Accept WebSocket connection
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Allow any origin in development
	})
	if err != nil {
		slog.Error("websocket accept failed", "error", err)
		return
	}

	// Set read limit
	conn.SetReadLimit(maxMessageSize)

	// Create client
	client := &Client{
		UserID: userID,
		ConnID: uuid.New().String(),
		Send:   make(chan []byte, 256),
		Hub:    h.hub,
		Done:   make(chan struct{}),
	}

	// Register client
	h.hub.Register(client)
	h.hub.SetOnlineStatus(userID)

	// Create a cancellable context for this connection
	ctx, cancel := context.WithCancel(r.Context())

	var wg sync.WaitGroup
	wg.Add(2)

	// Start reader and writer goroutines
	go func() {
		defer wg.Done()
		h.readPump(ctx, cancel, conn, client)
	}()
	go func() {
		defer wg.Done()
		h.writePump(ctx, cancel, conn, client)
	}()

	wg.Wait()

	// Cleanup
	h.hub.Unregister(client)
	h.hub.ClearOnlineStatus(userID)
	conn.Close(websocket.StatusNormalClosure, "")
}

// readPump reads messages from the WebSocket connection
func (h *Handler) readPump(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, client *Client) {
	defer cancel()

	for {
		_, msgBytes, err := conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) != -1 {
				slog.Debug("websocket closed", "user_id", client.UserID, "status", websocket.CloseStatus(err))
			} else {
				slog.Debug("websocket read error", "user_id", client.UserID, "error", err)
			}
			return
		}

		// Parse message
		var msg WSMessage
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			h.sendError(client, "INVALID_JSON", "Invalid message format")
			continue
		}

		// Route message by type
		switch msg.Type {
		case "message":
			h.handleMessage(ctx, client, msg)
		case "typing":
			h.handleTyping(ctx, client, msg)
		case "read":
			h.handleRead(ctx, client, msg)
		case "join_room":
			h.handleJoinRoom(ctx, client, msg)
		case "leave_room":
			h.handleLeaveRoom(ctx, client, msg)
		case "ping":
			h.handlePing(client)
		default:
			h.sendError(client, "UNKNOWN_TYPE", "Unknown message type: "+msg.Type)
		}
	}
}

// writePump writes messages from the send channel to the WebSocket connection
func (h *Handler) writePump(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, client *Client) {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return

		case message, ok := <-client.Send:
			if !ok {
				// Channel closed
				conn.Close(websocket.StatusNormalClosure, "")
				return
			}
			writeCtx, writeCancel := context.WithTimeout(ctx, writeTimeout)
			err := conn.Write(writeCtx, websocket.MessageText, message)
			writeCancel()
			if err != nil {
				slog.Debug("websocket write error", "user_id", client.UserID, "error", err)
				return
			}

		case <-ticker.C:
			// Send ping to keep connection alive
			writeCtx, writeCancel := context.WithTimeout(ctx, writeTimeout)
			err := conn.Ping(writeCtx)
			writeCancel()
			if err != nil {
				slog.Debug("websocket ping error", "user_id", client.UserID, "error", err)
				return
			}
			// Refresh online status
			h.hub.SetOnlineStatus(client.UserID)
		}
	}
}

// handleMessage processes incoming chat messages
func (h *Handler) handleMessage(ctx context.Context, client *Client, msg WSMessage) {
	// Rate limit check
	if !h.checkRateLimit(ctx, client.UserID) {
		h.sendError(client, "RATE_LIMITED", "Message rate limit exceeded (60/min)")
		return
	}

	if msg.ChatID == "" {
		h.sendError(client, "VALIDATION_ERROR", "chat_id is required")
		return
	}
	if msg.Content == "" {
		h.sendError(client, "VALIDATION_ERROR", "content is required")
		return
	}

	chatID, err := uuid.Parse(msg.ChatID)
	if err != nil {
		h.sendError(client, "INVALID_ID", "Invalid chat_id")
		return
	}

	// Send via chat service (handles access check, DB insert, etc.)
	msgResp, err := h.chatService.SendMessage(ctx, client.UserID, chatID, msg.Content, msg.ReplyTo)
	if err != nil {
		if appErr, ok := service.IsAppError(err); ok {
			h.sendError(client, appErr.Code, appErr.Error())
		} else {
			h.sendError(client, "INTERNAL_ERROR", "Failed to send message")
			slog.Error("ws: send message failed", "error", err)
		}
		return
	}

	// Include client_id in response for sender matching
	msgResp.ClientID = msg.ClientID

	// Build server → client response
	response := WSResponse{
		Type: "message",
		Data: msgResp,
	}
	data, _ := json.Marshal(response)

	// Broadcast to chat room
	h.hub.BroadcastToChat(chatID, client.UserID, data)
}

// handleTyping broadcasts a typing indicator
func (h *Handler) handleTyping(ctx context.Context, client *Client, msg WSMessage) {
	if msg.ChatID == "" {
		return
	}

	chatID, err := uuid.Parse(msg.ChatID)
	if err != nil {
		return
	}

	// Get user info for typing notification
	response := WSResponse{
		Type: "typing",
		Data: map[string]string{
			"chat_id":    msg.ChatID,
			"user_id":    client.UserID.String(),
		},
	}
	data, _ := json.Marshal(response)

	h.hub.BroadcastToChat(chatID, client.UserID, data)
}

// handleRead processes read receipts
func (h *Handler) handleRead(ctx context.Context, client *Client, msg WSMessage) {
	if msg.ChatID == "" {
		return
	}

	chatID, err := uuid.Parse(msg.ChatID)
	if err != nil {
		return
	}

	// Mark as read in DB
	if err := h.chatService.MarkAsRead(ctx, client.UserID, chatID); err != nil {
		slog.Error("ws: mark as read failed", "error", err)
		return
	}

	// Broadcast read receipt to chat room
	response := WSResponse{
		Type: "read",
		Data: map[string]string{
			"chat_id":      msg.ChatID,
			"user_id":      client.UserID.String(),
			"last_read_at": time.Now().Format(time.RFC3339),
		},
	}
	data, _ := json.Marshal(response)

	h.hub.BroadcastToChat(chatID, client.UserID, data)
}

// handleJoinRoom adds the user to a chat room for real-time updates
func (h *Handler) handleJoinRoom(ctx context.Context, client *Client, msg WSMessage) {
	if msg.ChatID == "" {
		h.sendError(client, "VALIDATION_ERROR", "chat_id is required")
		return
	}

	chatID, err := uuid.Parse(msg.ChatID)
	if err != nil {
		h.sendError(client, "INVALID_ID", "Invalid chat_id")
		return
	}

	// Verify user has access to this chat
	isMember, err := h.chatService.IsUserInChat(ctx, client.UserID, chatID)
	if err != nil || !isMember {
		h.sendError(client, "FORBIDDEN", "You are not a member of this chat")
		return
	}

	h.hub.JoinRoom(chatID, client.UserID)
}

// handleLeaveRoom removes the user from a chat room
func (h *Handler) handleLeaveRoom(ctx context.Context, client *Client, msg WSMessage) {
	if msg.ChatID == "" {
		return
	}

	chatID, err := uuid.Parse(msg.ChatID)
	if err != nil {
		return
	}

	h.hub.LeaveRoom(chatID, client.UserID)
}

// handlePing responds with a pong
func (h *Handler) handlePing(client *Client) {
	response := WSResponse{Type: "pong"}
	data, _ := json.Marshal(response)
	select {
	case client.Send <- data:
	default:
	}
}

// sendError sends an error message to the client
func (h *Handler) sendError(client *Client, code, message string) {
	response := WSResponse{
		Type: "error",
		Data: map[string]string{
			"code":    code,
			"message": message,
		},
	}
	data, _ := json.Marshal(response)
	select {
	case client.Send <- data:
	default:
	}
}

// checkRateLimit checks if the user is within the message rate limit
func (h *Handler) checkRateLimit(ctx context.Context, userID uuid.UUID) bool {
	if h.redis == nil {
		return true // Skip rate limiting if Redis is not available
	}

	key := "rate:ws_msg:" + userID.String()
	count, err := h.redis.Incr(ctx, key).Result()
	if err != nil {
		slog.Error("ws: rate limit check failed", "error", err)
		return true // Allow on error
	}

	if count == 1 {
		h.redis.Expire(ctx, key, msgRateWindow)
	}

	return count <= msgRateLimit
}
