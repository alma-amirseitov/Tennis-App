package service

import (
	"context"
	"fmt"
	"time"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ChatService handles chat business logic
type ChatService struct {
	repo *repository.Queries
}

// NewChatService creates a new ChatService
func NewChatService(repo *repository.Queries) *ChatService {
	return &ChatService{repo: repo}
}

// ChatResponse represents a chat in API responses
type ChatResponse struct {
	ID        string `json:"id"`
	ChatType  string `json:"chat_type"`
	OtherUser *UserBrief      `json:"other_user,omitempty"`
	Community *CommunityBrief `json:"community,omitempty"`
	Event     *EventBrief     `json:"event,omitempty"`
	Name      *string `json:"name,omitempty"`
	LastMessage *LastMessageInfo `json:"last_message,omitempty"`
	UnreadCount int    `json:"unread_count"`
	IsMuted     bool   `json:"is_muted"`
}

// UserBrief is a minimal user representation for chat
type UserBrief struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// CommunityBrief is a minimal community representation for chat
type CommunityBrief struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	LogoURL *string `json:"logo_url,omitempty"`
}

// EventBrief is a minimal event representation for chat
type EventBrief struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// LastMessageInfo represents the last message in a chat
type LastMessageInfo struct {
	Content   string `json:"content"`
	SenderID  string `json:"sender_id"`
	CreatedAt string `json:"created_at"`
}

// MessageResponse represents a message in API responses
type MessageResponse struct {
	ID        string     `json:"id"`
	ChatID    string     `json:"chat_id"`
	Sender    UserBrief  `json:"sender"`
	Content   string     `json:"content"`
	ReplyTo   *string    `json:"reply_to"`
	CreatedAt string     `json:"created_at"`
	ClientID  string     `json:"client_id,omitempty"`
}

// CreatePersonalChat creates or retrieves an existing personal chat
func (s *ChatService) CreatePersonalChat(ctx context.Context, user1ID, user2ID uuid.UUID) (string, bool, error) {
	if user1ID == user2ID {
		return "", false, ErrValidation.WithMessage("Cannot create chat with yourself")
	}

	// Enforce user1_id < user2_id constraint
	pgUser1 := uuidToPgtype(user1ID)
	pgUser2 := uuidToPgtype(user2ID)
	if user1ID.String() > user2ID.String() {
		pgUser1, pgUser2 = pgUser2, pgUser1
	}

	// Try to find existing chat
	existing, err := s.repo.GetPersonalChat(ctx, repository.GetPersonalChatParams{
		User1ID: pgUser1,
		User2ID: pgUser2,
	})
	if err == nil {
		chatID := pgtypeUUIDToStringRequired(existing.ID)
		return chatID, false, nil
	}
	if err != pgx.ErrNoRows {
		return "", false, fmt.Errorf("get personal chat: %w", err)
	}

	// Create new chat
	chat, err := s.repo.CreatePersonalChat(ctx, repository.CreatePersonalChatParams{
		User1ID: pgUser1,
		User2ID: pgUser2,
	})
	if err != nil {
		return "", false, fmt.Errorf("create personal chat: %w", err)
	}

	chatID := pgtypeUUIDToStringRequired(chat.ID)
	return chatID, true, nil
}

// GetOrCreateCommunityChat gets or creates a community chat
func (s *ChatService) GetOrCreateCommunityChat(ctx context.Context, communityID uuid.UUID) (string, error) {
	pgCommunityID := uuidToPgtype(communityID)

	existing, err := s.repo.GetCommunityChatByCommunityID(ctx, pgCommunityID)
	if err == nil {
		return pgtypeUUIDToStringRequired(existing.ID), nil
	}
	if err != pgx.ErrNoRows {
		return "", fmt.Errorf("get community chat: %w", err)
	}

	// Get community name for chat
	comm, err := s.repo.GetCommunityBasicInfo(ctx, pgCommunityID)
	if err != nil {
		return "", ErrCommunityNotFound
	}

	chat, err := s.repo.CreateCommunityChat(ctx, repository.CreateCommunityChatParams{
		CommunityID: pgCommunityID,
		Name:        pgtype.Text{String: comm.Name, Valid: true},
	})
	if err != nil {
		return "", fmt.Errorf("create community chat: %w", err)
	}

	return pgtypeUUIDToStringRequired(chat.ID), nil
}

// GetOrCreateEventChat gets or creates an event chat
func (s *ChatService) GetOrCreateEventChat(ctx context.Context, eventID uuid.UUID) (string, error) {
	pgEventID := uuidToPgtype(eventID)

	existing, err := s.repo.GetEventChatByEventID(ctx, pgEventID)
	if err == nil {
		return pgtypeUUIDToStringRequired(existing.ID), nil
	}
	if err != pgx.ErrNoRows {
		return "", fmt.Errorf("get event chat: %w", err)
	}

	// Get event title for chat name
	evt, err := s.repo.GetEventBasicInfo(ctx, pgEventID)
	if err != nil {
		return "", ErrEventNotFound
	}

	chat, err := s.repo.CreateEventChat(ctx, repository.CreateEventChatParams{
		EventID: pgEventID,
		Name:    pgtype.Text{String: evt.Title, Valid: true},
	})
	if err != nil {
		return "", fmt.Errorf("create event chat: %w", err)
	}

	return pgtypeUUIDToStringRequired(chat.ID), nil
}

// ListMyChats returns the user's chats with last message and unread count
func (s *ChatService) ListMyChats(ctx context.Context, userID uuid.UUID) ([]ChatResponse, error) {
	rows, err := s.repo.ListMyChats(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, fmt.Errorf("list my chats: %w", err)
	}

	chats := make([]ChatResponse, 0, len(rows))
	for _, row := range rows {
		chat := ChatResponse{
			ID:          pgtypeUUIDToStringRequired(row.ID),
			ChatType:    string(row.ChatType),
			UnreadCount: int(row.UnreadCount),
			IsMuted:     row.IsMuted,
		}

		if row.Name.Valid {
			chat.Name = &row.Name.String
		}

		// Last message
		if row.LastMessagePreview.Valid {
			chat.LastMessage = &LastMessageInfo{
				Content: row.LastMessagePreview.String,
			}
			if row.LastMessageAt.Valid {
				chat.LastMessage.CreatedAt = row.LastMessageAt.Time.Format(time.RFC3339)
			}
		}

		// Enrich based on chat type
		switch row.ChatType {
		case repository.ChatTypePersonal:
			// Get the other user's info
			if row.OtherUserID != nil {
				otherUIDBytes, ok := row.OtherUserID.([16]byte)
				if ok {
					otherPgUUID := pgtype.UUID{Bytes: otherUIDBytes, Valid: true}
					if userInfo, err := s.repo.GetUserBasicInfo(ctx, otherPgUUID); err == nil {
						chat.OtherUser = &UserBrief{
							ID:        pgtypeUUIDToStringRequired(userInfo.ID),
							FirstName: userInfo.FirstName.String,
							LastName:  userInfo.LastName.String,
						}
						if userInfo.AvatarUrl.Valid {
							chat.OtherUser.AvatarURL = &userInfo.AvatarUrl.String
						}
					}
				}
			}
		case repository.ChatTypeCommunity:
			if row.CommunityID.Valid {
				if comm, err := s.repo.GetCommunityBasicInfo(ctx, row.CommunityID); err == nil {
					chat.Community = &CommunityBrief{
						ID:   pgtypeUUIDToStringRequired(comm.ID),
						Name: comm.Name,
					}
					if comm.LogoUrl.Valid {
						chat.Community.LogoURL = &comm.LogoUrl.String
					}
				}
			}
		case repository.ChatTypeEvent:
			if row.EventID.Valid {
				if evt, err := s.repo.GetEventBasicInfo(ctx, row.EventID); err == nil {
					chat.Event = &EventBrief{
						ID:    pgtypeUUIDToStringRequired(evt.ID),
						Title: evt.Title,
					}
				}
			}
		}

		chats = append(chats, chat)
	}

	return chats, nil
}

// GetMessages returns messages in a chat with cursor-based pagination
func (s *ChatService) GetMessages(ctx context.Context, userID, chatID uuid.UUID, beforeCursor string, limit int) ([]MessageResponse, bool, error) {
	// Check user is in chat
	isMember, err := s.repo.IsUserInChat(ctx, repository.IsUserInChatParams{
		ChatID: uuidToPgtype(chatID),
		UserID: uuidToPgtype(userID),
	})
	if err != nil {
		return nil, false, fmt.Errorf("check chat membership: %w", err)
	}
	if !isMember {
		return nil, false, ErrForbidden.WithMessage("You are not a member of this chat")
	}

	if limit <= 0 || limit > 100 {
		limit = 50
	}

	params := repository.GetMessagesParams{
		ChatID:      uuidToPgtype(chatID),
		ResultLimit: int32(limit + 1), // Fetch one extra to determine has_more
	}

	if beforeCursor != "" {
		cursorUUID, err := uuid.Parse(beforeCursor)
		if err == nil {
			params.BeforeID = uuidToPgtype(cursorUUID)
		}
	}

	rows, err := s.repo.GetMessages(ctx, params)
	if err != nil {
		return nil, false, fmt.Errorf("get messages: %w", err)
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	messages := make([]MessageResponse, 0, len(rows))
	for _, row := range rows {
		msg := MessageResponse{
			ID:     pgtypeUUIDToStringRequired(row.ID),
			ChatID: pgtypeUUIDToStringRequired(row.ChatID),
			Sender: UserBrief{
				ID: pgtypeUUIDToStringRequired(row.SenderID),
			},
			Content:   row.Content,
			CreatedAt: row.CreatedAt.Time.Format(time.RFC3339),
		}
		if row.SenderFirstName.Valid {
			msg.Sender.FirstName = row.SenderFirstName.String
		}
		if row.SenderLastName.Valid {
			msg.Sender.LastName = row.SenderLastName.String
		}
		if row.SenderAvatarUrl.Valid {
			msg.Sender.AvatarURL = &row.SenderAvatarUrl.String
		}
		if row.ReplyToID.Valid {
			replyID := pgtypeUUIDToStringRequired(row.ReplyToID)
			msg.ReplyTo = &replyID
		}
		messages = append(messages, msg)
	}

	return messages, hasMore, nil
}

// SendMessage sends a message to a chat
func (s *ChatService) SendMessage(ctx context.Context, userID, chatID uuid.UUID, content string, replyToID *string) (*MessageResponse, error) {
	// Check user is in chat
	isMember, err := s.repo.IsUserInChat(ctx, repository.IsUserInChatParams{
		ChatID: uuidToPgtype(chatID),
		UserID: uuidToPgtype(userID),
	})
	if err != nil {
		return nil, fmt.Errorf("check chat membership: %w", err)
	}
	if !isMember {
		return nil, ErrForbidden.WithMessage("You are not a member of this chat")
	}

	if content == "" {
		return nil, ErrValidation.WithMessage("Message content cannot be empty")
	}
	if len(content) > 4000 {
		return nil, ErrValidation.WithMessage("Message content too long (max 4000 characters)")
	}

	params := repository.CreateMessageParams{
		ChatID:   uuidToPgtype(chatID),
		SenderID: uuidToPgtype(userID),
		Content:  content,
	}
	if replyToID != nil && *replyToID != "" {
		replyUUID, err := uuid.Parse(*replyToID)
		if err == nil {
			params.ReplyToID = uuidToPgtype(replyUUID)
		}
	}

	msg, err := s.repo.CreateMessage(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	// Update chat's last message
	preview := content
	if len(preview) > 100 {
		preview = preview[:100] + "..."
	}
	_ = s.repo.UpdateChatLastMessage(ctx, repository.UpdateChatLastMessageParams{
		LastMessageAt: msg.CreatedAt,
		Preview:       pgtype.Text{String: preview, Valid: true},
		ID:            uuidToPgtype(chatID),
	})

	// Get sender info
	senderInfo, _ := s.repo.GetUserBasicInfo(ctx, uuidToPgtype(userID))

	resp := &MessageResponse{
		ID:     pgtypeUUIDToStringRequired(msg.ID),
		ChatID: pgtypeUUIDToStringRequired(msg.ChatID),
		Sender: UserBrief{
			ID: pgtypeUUIDToStringRequired(msg.SenderID),
		},
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt.Time.Format(time.RFC3339),
	}
	if senderInfo.FirstName.Valid {
		resp.Sender.FirstName = senderInfo.FirstName.String
	}
	if senderInfo.LastName.Valid {
		resp.Sender.LastName = senderInfo.LastName.String
	}
	if senderInfo.AvatarUrl.Valid {
		resp.Sender.AvatarURL = &senderInfo.AvatarUrl.String
	}
	if msg.ReplyToID.Valid {
		replyID := pgtypeUUIDToStringRequired(msg.ReplyToID)
		resp.ReplyTo = &replyID
	}

	return resp, nil
}

// MarkAsRead marks messages in a chat as read
func (s *ChatService) MarkAsRead(ctx context.Context, userID, chatID uuid.UUID) error {
	// Check user is in chat
	isMember, err := s.repo.IsUserInChat(ctx, repository.IsUserInChatParams{
		ChatID: uuidToPgtype(chatID),
		UserID: uuidToPgtype(userID),
	})
	if err != nil {
		return fmt.Errorf("check chat membership: %w", err)
	}
	if !isMember {
		return ErrForbidden.WithMessage("You are not a member of this chat")
	}

	return s.repo.UpsertChatReadStatus(ctx, repository.UpsertChatReadStatusParams{
		ChatID:     uuidToPgtype(chatID),
		UserID:     uuidToPgtype(userID),
		LastReadAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
}

// UpdateMuted updates the muted status of a chat
func (s *ChatService) UpdateMuted(ctx context.Context, userID, chatID uuid.UUID, isMuted bool) error {
	isMember, err := s.repo.IsUserInChat(ctx, repository.IsUserInChatParams{
		ChatID: uuidToPgtype(chatID),
		UserID: uuidToPgtype(userID),
	})
	if err != nil {
		return fmt.Errorf("check chat membership: %w", err)
	}
	if !isMember {
		return ErrForbidden.WithMessage("You are not a member of this chat")
	}

	return s.repo.UpdateChatMuted(ctx, repository.UpdateChatMutedParams{
		ChatID:  uuidToPgtype(chatID),
		UserID:  uuidToPgtype(userID),
		IsMuted: pgtype.Bool{Bool: isMuted, Valid: true},
	})
}

// GetUnreadCount returns the total unread count across all chats
func (s *ChatService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := s.repo.GetTotalUnreadCount(ctx, uuidToPgtype(userID))
	if err != nil {
		return 0, fmt.Errorf("get total unread count: %w", err)
	}
	return int(count), nil
}

// GetChatMembers returns the user IDs of all members in a chat
func (s *ChatService) GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	pgChatID := uuidToPgtype(chatID)

	chat, err := s.repo.GetChatByID(ctx, pgChatID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrChatNotFound
		}
		return nil, fmt.Errorf("get chat: %w", err)
	}

	var memberIDs []uuid.UUID

	switch chat.ChatType {
	case repository.ChatTypePersonal:
		rows, err := s.repo.GetChatMembersForPersonal(ctx, pgChatID)
		if err != nil {
			return nil, fmt.Errorf("get personal chat members: %w", err)
		}
		for _, row := range rows {
			if row.User1ID.Valid {
				id, _ := uuid.FromBytes(row.User1ID.Bytes[:])
				memberIDs = append(memberIDs, id)
			}
			if row.User2ID.Valid {
				id, _ := uuid.FromBytes(row.User2ID.Bytes[:])
				memberIDs = append(memberIDs, id)
			}
		}
	case repository.ChatTypeCommunity:
		pgIDs, err := s.repo.GetChatMembersForCommunity(ctx, pgChatID)
		if err != nil {
			return nil, fmt.Errorf("get community chat members: %w", err)
		}
		for _, pgID := range pgIDs {
			if pgID.Valid {
				id, _ := uuid.FromBytes(pgID.Bytes[:])
				memberIDs = append(memberIDs, id)
			}
		}
	case repository.ChatTypeEvent:
		pgIDs, err := s.repo.GetChatMembersForEvent(ctx, pgChatID)
		if err != nil {
			return nil, fmt.Errorf("get event chat members: %w", err)
		}
		for _, pgID := range pgIDs {
			if pgID.Valid {
				id, _ := uuid.FromBytes(pgID.Bytes[:])
				memberIDs = append(memberIDs, id)
			}
		}
	}

	return memberIDs, nil
}

// IsUserInChat checks if a user is a member of a chat
func (s *ChatService) IsUserInChat(ctx context.Context, userID, chatID uuid.UUID) (bool, error) {
	return s.repo.IsUserInChat(ctx, repository.IsUserInChatParams{
		ChatID: uuidToPgtype(chatID),
		UserID: uuidToPgtype(userID),
	})
}
