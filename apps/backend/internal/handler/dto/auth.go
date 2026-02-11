package dto

// SendOTPRequest represents request to send OTP
type SendOTPRequest struct {
	Phone string `json:"phone" validate:"required,kz_phone"`
}

// SendOTPResponse represents response with session ID
type SendOTPResponse struct {
	SessionID  string `json:"session_id"`
	ExpiresIn  int    `json:"expires_in"`
	RetryAfter int    `json:"retry_after"`
}

// VerifyOTPRequest represents request to verify OTP
type VerifyOTPRequest struct {
	SessionID string `json:"session_id" validate:"required,uuid"`
	Code      string `json:"code" validate:"required,len=4,numeric"`
}

// VerifyOTPResponse represents response after OTP verification
type VerifyOTPResponse struct {
	IsNew        bool        `json:"is_new"`
	AccessToken  string      `json:"access_token,omitempty"`
	RefreshToken string      `json:"refresh_token,omitempty"`
	TempToken    string      `json:"temp_token,omitempty"`
	User         interface{} `json:"user"`
}

// RefreshTokenRequest represents request to refresh access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents response with new tokens
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ProfileSetupRequest represents request to complete user profile
type ProfileSetupRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Gender    string `json:"gender" validate:"required,oneof=male female"`
	BirthYear int16  `json:"birth_year" validate:"required,min=1940,max=2015"`
	City      string `json:"city" validate:"required,min=2,max=100"`
	District  string `json:"district" validate:"omitempty,min=2,max=100"`
	Language  string `json:"language" validate:"required,oneof=ru kk en"`
}

// ProfileSetupResponse represents response after profile setup
type ProfileSetupResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         interface{} `json:"user"`
}
