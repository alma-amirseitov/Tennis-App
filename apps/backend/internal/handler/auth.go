package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/handler/dto"
	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/handler/middleware"
	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/pkg/validator"
	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/service"
	"github.com/google/uuid"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService  *service.AuthService
	tokenService *service.TokenService
	validator    *validator.Validator
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *service.AuthService, tokenService *service.TokenService, validator *validator.Validator) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		tokenService: tokenService,
		validator:    validator,
	}
}

// SendOTP handles POST /v1/auth/otp/send
func (h *AuthHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	var req dto.SendOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondValidationError(w, err)
		return
	}

	sessionID, err := h.authService.SendOTP(r.Context(), req.Phone)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, dto.SendOTPResponse{
		SessionID:  sessionID,
		ExpiresIn:  300,
		RetryAfter: 60,
	})
}

// VerifyOTP handles POST /v1/auth/otp/verify
func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondValidationError(w, err)
		return
	}

	result, err := h.authService.VerifyOTP(r.Context(), req.SessionID, req.Code)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, dto.VerifyOTPResponse{
		IsNew:        result.IsNewUser,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TempToken:    result.TempToken,
		User:         result.User,
	})
}

// RefreshToken handles POST /v1/auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondValidationError(w, err)
		return
	}

	accessToken, refreshToken, err := h.tokenService.RefreshTokens(r.Context(), req.RefreshToken)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// ProfileSetup handles POST /v1/auth/profile/setup
func (h *AuthHandler) ProfileSetup(w http.ResponseWriter, r *http.Request) {
	var req dto.ProfileSetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondValidationError(w, err)
		return
	}

	// Get user ID from context (set by Auth middleware)
	userIDStr := middleware.GetUserID(r.Context())
	if userIDStr == "" {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID")
		return
	}

	// Call service
	input := service.ProfileSetupInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Gender:    req.Gender,
		BirthYear: req.BirthYear,
		City:      req.City,
		District:  req.District,
		Language:  req.Language,
	}

	result, err := h.authService.ProfileSetup(r.Context(), userID, input)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, dto.ProfileSetupResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User:         result.User,
	})
}
