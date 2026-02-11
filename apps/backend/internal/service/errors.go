package service

import "errors"

// AppError represents an application error with HTTP status code
type AppError struct {
	Code    string
	Status  int
	Message string
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Code
}

// WithMessage returns a copy of the error with a custom message
func (e *AppError) WithMessage(msg string) *AppError {
	return &AppError{
		Code:    e.Code,
		Status:  e.Status,
		Message: msg,
	}
}

// Auth errors (401)
var (
	ErrUnauthorized = &AppError{Code: "UNAUTHORIZED", Status: 401}
	ErrTokenExpired = &AppError{Code: "TOKEN_EXPIRED", Status: 401}
	ErrTokenInvalid = &AppError{Code: "TOKEN_INVALID", Status: 401}
	ErrTokenRevoked = &AppError{Code: "TOKEN_REVOKED", Status: 401}
)

// OTP errors (400)
var (
	ErrOTPSessionExpired = &AppError{Code: "OTP_SESSION_EXPIRED", Status: 400}
	ErrOTPInvalidCode    = &AppError{Code: "OTP_INVALID_CODE", Status: 400}
	ErrOTPMaxAttempts    = &AppError{Code: "OTP_MAX_ATTEMPTS", Status: 400}
)

// Validation errors (400)
var (
	ErrValidation         = &AppError{Code: "VALIDATION_ERROR", Status: 400}
	ErrInvalidPhoneFormat = &AppError{Code: "INVALID_PHONE_FORMAT", Status: 400}
	ErrInvalidJSON        = &AppError{Code: "INVALID_JSON", Status: 400}
)

// Permission errors (403)
var (
	ErrForbidden          = &AppError{Code: "FORBIDDEN", Status: 403}
	ErrNotCommunityMember = &AppError{Code: "NOT_COMMUNITY_MEMBER", Status: 403}
	ErrInsufficientRole   = &AppError{Code: "INSUFFICIENT_ROLE", Status: 403}
)

// Not Found (404)
var (
	ErrNotFound          = &AppError{Code: "NOT_FOUND", Status: 404}
	ErrUserNotFound      = &AppError{Code: "USER_NOT_FOUND", Status: 404}
	ErrEventNotFound     = &AppError{Code: "EVENT_NOT_FOUND", Status: 404}
	ErrCommunityNotFound = &AppError{Code: "COMMUNITY_NOT_FOUND", Status: 404}
	ErrMatchNotFound     = &AppError{Code: "MATCH_NOT_FOUND", Status: 404}
	ErrChatNotFound      = &AppError{Code: "CHAT_NOT_FOUND", Status: 404}
)

// Conflict (409)
var (
	ErrAlreadyExists       = &AppError{Code: "ALREADY_EXISTS", Status: 409}
	ErrAlreadyMember       = &AppError{Code: "ALREADY_MEMBER", Status: 409}
	ErrAlreadyJoinedEvent  = &AppError{Code: "ALREADY_JOINED_EVENT", Status: 409}
	ErrAlreadyFriends      = &AppError{Code: "ALREADY_FRIENDS", Status: 409}
	ErrProfileAlreadySet   = &AppError{Code: "PROFILE_ALREADY_SET", Status: 409}
	ErrResultAlreadySubmit = &AppError{Code: "RESULT_ALREADY_SUBMITTED", Status: 409}
)

// Rate Limit (429)
var (
	ErrRateLimited    = &AppError{Code: "RATE_LIMITED", Status: 429}
	ErrSMSRateLimited = &AppError{Code: "SMS_RATE_LIMITED", Status: 429}
)

// Server errors (500)
var (
	ErrInternal = &AppError{Code: "INTERNAL_ERROR", Status: 500}
)

// Helper to check if error is AppError
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
