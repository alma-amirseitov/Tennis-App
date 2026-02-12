package service

import (
	"context"
	"log/slog"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/config"
	"github.com/google/uuid"
)

// FirebaseService is a thin abstraction over push delivery.
// In development it runs in mock mode and only logs pushes.
type FirebaseService struct {
	logger *slog.Logger
	isMock bool
}

// NewFirebaseService creates a new FirebaseService.
// For now, we always run in mock mode. Real FCM integration can be
// implemented later without changing the service API.
func NewFirebaseService(logger *slog.Logger, cfg *config.Config) *FirebaseService {
	isMock := cfg.Environment != "production" || cfg.FirebaseCredentials == ""

	if isMock {
		logger.Info("FirebaseService initialized in mock mode",
			"environment", cfg.Environment,
			"credentials", cfg.FirebaseCredentials != "")
	}

	return &FirebaseService{
		logger: logger,
		isMock: isMock,
	}
}

// SendPush sends a push notification to a single device token.
func (f *FirebaseService) SendPush(_ context.Context, deviceToken, title, body string, data map[string]any) error {
	if f == nil {
		return nil
	}

	// Mock implementation: just log. Real FCM integration can replace this.
	f.logger.Info("mock push notification",
		"device_token", deviceToken,
		"title", title,
		"body", body,
		"data", data,
	)
	return nil
}

// SendToTopic sends a push notification to a topic.
func (f *FirebaseService) SendToTopic(_ context.Context, topic, title, body string, data map[string]any) error {
	if f == nil {
		return nil
	}

	// Mock implementation: just log.
	f.logger.Info("mock topic notification",
		"topic", topic,
		"title", title,
		"body", body,
		"data", data,
	)
	return nil
}

// SendToUser is a convenience helper that, for now, only logs.
// In the future this can look up the user's device tokens and fan-out.
func (f *FirebaseService) SendToUser(_ context.Context, userID uuid.UUID, title, body string, data map[string]any) error {
	if f == nil {
		return nil
	}

	f.logger.Info("mock user notification",
		"user_id", userID,
		"title", title,
		"body", body,
		"data", data,
	)
	return nil
}
