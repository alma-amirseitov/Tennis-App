package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/repository"
	"github.com/google/uuid"
)

// NotificationService handles creation and delivery of in-app notifications.
type NotificationService struct {
	repo     *repository.Queries
	logger   *slog.Logger
	firebase *FirebaseService
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(repo *repository.Queries, logger *slog.Logger, firebase *FirebaseService) *NotificationService {
	return &NotificationService{
		repo:     repo,
		logger:   logger,
		firebase: firebase,
	}
}

// Notification represents a notification in API responses.
type Notification struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	Data      map[string]any `json:"data,omitempty"`
	IsRead    bool           `json:"is_read"`
	CreatedAt string         `json:"created_at"`
	ExpiresAt *string        `json:"expires_at,omitempty"`
}

// Create creates a notification in the database and (optionally) sends a push.
func (s *NotificationService) Create(
	ctx context.Context,
	userID uuid.UUID,
	notificationType string,
	title string,
	body string,
	data map[string]any,
) (*Notification, error) {
	if notificationType == "" {
		return nil, ErrValidation.WithMessage("notification type is required")
	}
	if title == "" {
		return nil, ErrValidation.WithMessage("title is required")
	}
	if body == "" {
		return nil, ErrValidation.WithMessage("body is required")
	}

	pgUserID := uuidToPgtype(userID)

	var rawData []byte
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("marshal notification data: %w", err)
		}
		rawData = b
	}

	dbNotif, err := s.repo.CreateNotification(ctx, repository.CreateNotificationParams{
		UserID: pgUserID,
		Type:   repository.NotificationType(notificationType),
		Title:  title,
		Body:   body,
		Data:   rawData,
	})
	if err != nil {
		return nil, fmt.Errorf("create notification: %w", err)
	}

	notif := buildNotification(dbNotif)

	// Best-effort push delivery; log but don't fail the request on push error.
	if s.firebase != nil {
		if err := s.firebase.SendToUser(ctx, userID, title, body, data); err != nil {
			s.logger.Warn("failed to send push notification",
				"user_id", userID,
				"type", notificationType,
				"error", err,
			)
		}
	}

	return notif, nil
}

// List returns notifications for a user with pagination.
func (s *NotificationService) List(ctx context.Context, userID uuid.UUID, page, perPage int) ([]Notification, *PaginationInfo, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	rows, err := s.repo.ListNotifications(ctx, repository.ListNotificationsParams{
		UserID:       uuidToPgtype(userID),
		ResultLimit:  int32(perPage),
		ResultOffset: int32(offset),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("list notifications: %w", err)
	}

	total64, err := s.repo.CountNotifications(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, nil, fmt.Errorf("count notifications: %w", err)
	}
	total := int(total64)
	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	result := make([]Notification, 0, len(rows))
	for _, n := range rows {
		result = append(result, *buildNotification(n))
	}

	pagination := &PaginationInfo{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}

	return result, pagination, nil
}

// MarkAsRead marks a single notification as read.
func (s *NotificationService) MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error {
	if err := s.repo.MarkNotificationRead(ctx, repository.MarkNotificationReadParams{
		ID:     uuidToPgtype(notificationID),
		UserID: uuidToPgtype(userID),
	}); err != nil {
		return fmt.Errorf("mark notification read: %w", err)
	}
	return nil
}

// MarkAllAsRead marks all notifications for a user as read.
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.MarkAllNotificationsRead(ctx, uuidToPgtype(userID)); err != nil {
		return fmt.Errorf("mark all notifications read: %w", err)
	}
	return nil
}

// GetUnreadCount returns the total unread notifications for a user.
func (s *NotificationService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := s.repo.GetUnreadNotificationCount(ctx, uuidToPgtype(userID))
	if err != nil {
		return 0, fmt.Errorf("get unread notifications count: %w", err)
	}
	return int(count), nil
}

// Delete deletes a notification for a user.
func (s *NotificationService) Delete(ctx context.Context, userID, notificationID uuid.UUID) error {
	if err := s.repo.DeleteNotification(ctx, repository.DeleteNotificationParams{
		ID:     uuidToPgtype(notificationID),
		UserID: uuidToPgtype(userID),
	}); err != nil {
		return fmt.Errorf("delete notification: %w", err)
	}
	return nil
}

// buildNotification converts a repository.Notification into service Notification struct.
func buildNotification(n repository.Notification) *Notification {
	id := pgtypeUUIDToStringRequired(n.ID)

	var data map[string]any
	if len(n.Data) > 0 {
		if err := json.Unmarshal(n.Data, &data); err != nil {
			// If data is invalid JSON, ignore it rather than failing the whole response.
			data = nil
		}
	}

	createdAt := ""
	if n.CreatedAt.Valid {
		createdAt = n.CreatedAt.Time.Format(time.RFC3339)
	}

	var expiresAt *string
	if n.ExpiresAt.Valid {
		s := n.ExpiresAt.Time.Format(time.RFC3339)
		expiresAt = &s
	}

	return &Notification{
		ID:        id,
		Type:      string(n.Type),
		Title:     n.Title,
		Body:      n.Body,
		Data:      data,
		IsRead:    n.IsRead.Bool,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}
}
