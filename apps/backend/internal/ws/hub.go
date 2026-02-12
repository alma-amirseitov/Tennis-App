package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

// Hub manages WebSocket connections, rooms, and message broadcasting
type Hub struct {
	// connections maps userID → set of client connections
	connections map[uuid.UUID]map[*Client]struct{}

	// rooms maps chatID → set of userIDs subscribed to this room
	rooms map[uuid.UUID]map[uuid.UUID]struct{}

	// channels
	register   chan *Client
	unregister chan *Client
	broadcast  chan *BroadcastMessage
	joinRoom   chan *RoomAction
	leaveRoom  chan *RoomAction

	// redis for pub/sub across instances
	redis *goredis.Client

	// mutex for concurrent map access
	mu sync.RWMutex

	// context for shutdown
	ctx    context.Context
	cancel context.CancelFunc
}

// Client represents a single WebSocket connection
type Client struct {
	UserID   uuid.UUID
	ConnID   string
	Send     chan []byte
	Hub      *Hub
	Done     chan struct{}
}

// BroadcastMessage represents a message to broadcast to a chat room
type BroadcastMessage struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
	Data     []byte
}

// RoomAction represents a join/leave room action
type RoomAction struct {
	ChatID uuid.UUID
	UserID uuid.UUID
}

// WSMessage is the envelope for WebSocket messages
type WSMessage struct {
	Type     string          `json:"type"`
	ChatID   string          `json:"chat_id,omitempty"`
	Content  string          `json:"content,omitempty"`
	ReplyTo  *string         `json:"reply_to,omitempty"`
	ClientID string          `json:"client_id,omitempty"`
	Data     json.RawMessage `json:"data,omitempty"`
}

// WSResponse is the server → client message envelope
type WSResponse struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

// NewHub creates a new WebSocket hub
func NewHub(redis *goredis.Client) *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Hub{
		connections: make(map[uuid.UUID]map[*Client]struct{}),
		rooms:       make(map[uuid.UUID]map[uuid.UUID]struct{}),
		register:    make(chan *Client, 256),
		unregister:  make(chan *Client, 256),
		broadcast:   make(chan *BroadcastMessage, 256),
		joinRoom:    make(chan *RoomAction, 256),
		leaveRoom:   make(chan *RoomAction, 256),
		redis:       redis,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Run starts the hub's main event loop
func (h *Hub) Run() {
	// Start Redis subscriber in background
	go h.subscribeRedis()

	for {
		select {
		case <-h.ctx.Done():
			return

		case client := <-h.register:
			h.mu.Lock()
			if h.connections[client.UserID] == nil {
				h.connections[client.UserID] = make(map[*Client]struct{})
			}
			h.connections[client.UserID][client] = struct{}{}
			h.mu.Unlock()
			slog.Debug("client registered", "user_id", client.UserID, "conn_id", client.ConnID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.connections[client.UserID]; ok {
				delete(clients, client)
				if len(clients) == 0 {
					delete(h.connections, client.UserID)
					// Remove user from all rooms
					for chatID, users := range h.rooms {
						delete(users, client.UserID)
						if len(users) == 0 {
							delete(h.rooms, chatID)
						}
					}
				}
			}
			close(client.Send)
			h.mu.Unlock()
			slog.Debug("client unregistered", "user_id", client.UserID, "conn_id", client.ConnID)

		case action := <-h.joinRoom:
			h.mu.Lock()
			if h.rooms[action.ChatID] == nil {
				h.rooms[action.ChatID] = make(map[uuid.UUID]struct{})
			}
			h.rooms[action.ChatID][action.UserID] = struct{}{}
			h.mu.Unlock()

		case action := <-h.leaveRoom:
			h.mu.Lock()
			if users, ok := h.rooms[action.ChatID]; ok {
				delete(users, action.UserID)
				if len(users) == 0 {
					delete(h.rooms, action.ChatID)
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.broadcastToRoom(msg)
			// Publish to Redis for multi-instance support
			h.publishRedis(msg)
		}
	}
}

// Stop shuts down the hub
func (h *Hub) Stop() {
	h.cancel()
}

// Register registers a new client
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// JoinRoom adds a user to a chat room
func (h *Hub) JoinRoom(chatID, userID uuid.UUID) {
	h.joinRoom <- &RoomAction{ChatID: chatID, UserID: userID}
}

// LeaveRoom removes a user from a chat room
func (h *Hub) LeaveRoom(chatID, userID uuid.UUID) {
	h.leaveRoom <- &RoomAction{ChatID: chatID, UserID: userID}
}

// BroadcastToChat sends a message to all users in a chat room
func (h *Hub) BroadcastToChat(chatID, senderID uuid.UUID, data []byte) {
	h.broadcast <- &BroadcastMessage{
		ChatID:   chatID,
		SenderID: senderID,
		Data:     data,
	}
}

// SendToUser sends a message directly to a specific user
func (h *Hub) SendToUser(userID uuid.UUID, data []byte) {
	h.mu.RLock()
	clients, ok := h.connections[userID]
	h.mu.RUnlock()
	if !ok {
		return
	}
	for client := range clients {
		select {
		case client.Send <- data:
		default:
			// Client's send buffer is full, skip
		}
	}
}

// broadcastToRoom sends a message to all connected users in a room
func (h *Hub) broadcastToRoom(msg *BroadcastMessage) {
	h.mu.RLock()
	users, ok := h.rooms[msg.ChatID]
	if !ok {
		h.mu.RUnlock()
		return
	}
	// Copy user IDs to avoid holding the lock while sending
	userIDs := make([]uuid.UUID, 0, len(users))
	for uid := range users {
		userIDs = append(userIDs, uid)
	}
	h.mu.RUnlock()

	for _, uid := range userIDs {
		h.mu.RLock()
		clients, ok := h.connections[uid]
		h.mu.RUnlock()
		if !ok {
			continue
		}
		for client := range clients {
			select {
			case client.Send <- msg.Data:
			default:
				// Skip if buffer full
			}
		}
	}
}

// Redis pub/sub methods

const redisChannelPrefix = "ws:channel:"

type redisPubSubMessage struct {
	ChatID   string `json:"chat_id"`
	SenderID string `json:"sender_id"`
	Data     json.RawMessage `json:"data"`
}

// publishRedis publishes a broadcast message to Redis for other instances
func (h *Hub) publishRedis(msg *BroadcastMessage) {
	if h.redis == nil {
		return
	}

	pubMsg := redisPubSubMessage{
		ChatID:   msg.ChatID.String(),
		SenderID: msg.SenderID.String(),
		Data:     msg.Data,
	}
	payload, err := json.Marshal(pubMsg)
	if err != nil {
		slog.Error("failed to marshal redis pub/sub message", "error", err)
		return
	}

	channel := redisChannelPrefix + msg.ChatID.String()
	if err := h.redis.Publish(h.ctx, channel, payload).Err(); err != nil {
		slog.Error("failed to publish to redis", "channel", channel, "error", err)
	}
}

// subscribeRedis listens for messages from other instances via Redis pub/sub
func (h *Hub) subscribeRedis() {
	if h.redis == nil {
		return
	}

	// Subscribe to pattern for all chat channels
	pubsub := h.redis.PSubscribe(h.ctx, redisChannelPrefix+"*")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-h.ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			var pubMsg redisPubSubMessage
			if err := json.Unmarshal([]byte(msg.Payload), &pubMsg); err != nil {
				slog.Error("failed to unmarshal redis pub/sub message", "error", err)
				continue
			}

			chatID, err := uuid.Parse(pubMsg.ChatID)
			if err != nil {
				continue
			}

			// Broadcast locally (don't re-publish to Redis)
			h.broadcastToRoom(&BroadcastMessage{
				ChatID: chatID,
				Data:   pubMsg.Data,
			})
		}
	}
}

// IsOnline checks if a user has any active connections
func (h *Hub) IsOnline(userID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients, ok := h.connections[userID]
	return ok && len(clients) > 0
}

// SetOnlineStatus sets the user's online status in Redis
func (h *Hub) SetOnlineStatus(userID uuid.UUID) {
	if h.redis == nil {
		return
	}
	key := "online:" + userID.String()
	h.redis.Set(h.ctx, key, "1", 5*time.Minute)
}

// ClearOnlineStatus removes the user's online status from Redis
func (h *Hub) ClearOnlineStatus(userID uuid.UUID) {
	if h.redis == nil {
		return
	}
	key := "online:" + userID.String()
	h.redis.Del(h.ctx, key)
}
