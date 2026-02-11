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

// QuizHandler handles quiz endpoints
type QuizHandler struct {
	quizService *service.QuizService
	validator   *validator.Validator
}

// NewQuizHandler creates a new QuizHandler
func NewQuizHandler(quizService *service.QuizService, validator *validator.Validator) *QuizHandler {
	return &QuizHandler{
		quizService: quizService,
		validator:   validator,
	}
}

// GetQuestions handles GET /v1/quiz
func (h *QuizHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	questions := h.quizService.GetQuestions()

	respondJSON(w, http.StatusOK, dto.QuizQuestionsResponse{
		Questions: questions,
	})
}

// SubmitAnswers handles POST /v1/quiz
func (h *QuizHandler) SubmitAnswers(w http.ResponseWriter, r *http.Request) {
	var req dto.QuizSubmitRequest
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

	// Submit answers and calculate NTRP
	result, err := h.quizService.SubmitAnswers(r.Context(), userID, req.ToServiceAnswers())
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, dto.QuizSubmitResponse{
		NTRPLevel:     result.NTRPLevel,
		InitialRating: result.InitialRating,
	})
}
