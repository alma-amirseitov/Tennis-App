package service

import (
	"context"
	"fmt"

	"github.com/alma-amirseitov/Tennis-App/apps/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// QuizService handles quiz logic
type QuizService struct {
	repo *repository.Queries
}

// NewQuizService creates a new QuizService
func NewQuizService(repo *repository.Queries) *QuizService {
	return &QuizService{repo: repo}
}

// Question represents a quiz question
type Question struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Options []Option `json:"options"`
}

// Option represents an answer option
type Option struct {
	ID     string  `json:"id"`
	Text   string  `json:"text"`
	Weight float64 `json:"-"` // Hidden from JSON response
}

// Answer represents a user's answer
type Answer struct {
	QuestionID string `json:"question_id"`
	AnswerID   string `json:"answer_id"`
}

// QuizResult represents quiz submission result
type QuizResult struct {
	NTRPLevel     float64 `json:"ntrp_level"`
	InitialRating int     `json:"initial_rating"`
}

// GetQuestions returns hardcoded quiz questions
func (s *QuizService) GetQuestions() []Question {
	return []Question{
		{
			ID:   "q1",
			Text: "Как долго вы играете в теннис?",
			Options: []Option{
				{ID: "q1_a1", Text: "Меньше 6 месяцев", Weight: 1.0},
				{ID: "q1_a2", Text: "6 месяцев - 1 год", Weight: 1.5},
				{ID: "q1_a3", Text: "1-2 года", Weight: 2.0},
				{ID: "q1_a4", Text: "2-3 года", Weight: 2.5},
				{ID: "q1_a5", Text: "Больше 3 лет", Weight: 3.0},
			},
		},
		{
			ID:   "q2",
			Text: "Как часто вы играете?",
			Options: []Option{
				{ID: "q2_a1", Text: "Раз в месяц или реже", Weight: 1.0},
				{ID: "q2_a2", Text: "2-3 раза в месяц", Weight: 1.5},
				{ID: "q2_a3", Text: "1-2 раза в неделю", Weight: 2.0},
				{ID: "q2_a4", Text: "3-4 раза в неделю", Weight: 2.5},
				{ID: "q2_a5", Text: "Почти каждый день", Weight: 3.0},
			},
		},
		{
			ID:   "q3",
			Text: "Оцените свою физическую форму",
			Options: []Option{
				{ID: "q3_a1", Text: "Новичок, часто устаю", Weight: 1.0},
				{ID: "q3_a2", Text: "Средняя выносливость", Weight: 1.5},
				{ID: "q3_a3", Text: "Хорошая форма", Weight: 2.0},
				{ID: "q3_a4", Text: "Отличная выносливость", Weight: 2.5},
				{ID: "q3_a5", Text: "Профессиональный уровень", Weight: 3.0},
			},
		},
		{
			ID:   "q4",
			Text: "Насколько уверенно вы контролируете мяч?",
			Options: []Option{
				{ID: "q4_a1", Text: "Часто промахиваюсь", Weight: 1.0},
				{ID: "q4_a2", Text: "Могу поддерживать простой розыгрыш", Weight: 1.5},
				{ID: "q4_a3", Text: "Стабильные удары с обеих сторон", Weight: 2.5},
				{ID: "q4_a4", Text: "Могу менять направление и темп", Weight: 3.5},
				{ID: "q4_a5", Text: "Полный контроль, варьирую кручение", Weight: 4.5},
			},
		},
		{
			ID:   "q5",
			Text: "Участвовали ли вы в турнирах?",
			Options: []Option{
				{ID: "q5_a1", Text: "Никогда", Weight: 1.0},
				{ID: "q5_a2", Text: "Только любительские", Weight: 2.0},
				{ID: "q5_a3", Text: "Клубные турниры", Weight: 3.0},
				{ID: "q5_a4", Text: "Региональные турниры", Weight: 4.0},
				{ID: "q5_a5", Text: "Профессиональный уровень", Weight: 5.0},
			},
		},
	}
}

// SubmitAnswers calculates NTRP level and initial rating based on answers
func (s *QuizService) SubmitAnswers(ctx context.Context, userID uuid.UUID, answers []Answer) (*QuizResult, error) {
	// Get questions to access weights
	questions := s.GetQuestions()
	questionMap := make(map[string]Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	// Calculate total weight
	totalWeight := 0.0
	answeredCount := 0

	for _, answer := range answers {
		question, ok := questionMap[answer.QuestionID]
		if !ok {
			return nil, fmt.Errorf("invalid question_id %s", answer.QuestionID)
		}

		// Find answer option and its weight
		found := false
		for _, option := range question.Options {
			if option.ID == answer.AnswerID {
				totalWeight += option.Weight
				answeredCount++
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("invalid answer_id %s for question %s", answer.AnswerID, answer.QuestionID)
		}
	}

	if answeredCount == 0 {
		return nil, fmt.Errorf("no valid answers provided")
	}

	// Calculate average weight
	avgWeight := totalWeight / float64(answeredCount)

	// Map weight to NTRP level (1.0 - 5.0)
	ntrpLevel := calculateNTRPLevel(avgWeight)

	// Map NTRP to initial rating (100 - 1000)
	initialRating := ntrpLevelToRating(ntrpLevel)

	// Convert NTRP to pgtype.Numeric
	ntrpNumeric := pgtype.Numeric{}
	if err := ntrpNumeric.Scan(fmt.Sprintf("%.1f", ntrpLevel)); err != nil {
		return nil, fmt.Errorf("convert NTRP to numeric: %w", err)
	}

	// Convert rating to pgtype.Numeric
	ratingNumeric := pgtype.Numeric{}
	if err := ratingNumeric.Scan(fmt.Sprintf("%d", initialRating)); err != nil {
		return nil, fmt.Errorf("convert rating to numeric: %w", err)
	}

	// Update user's NTRP level in database
	_, err := s.repo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:               pgtype.UUID{Bytes: userID, Valid: true},
		NtrpLevel:        ntrpNumeric,
		QuizCompleted:    pgtype.Bool{Bool: true, Valid: true},
		GlobalRating:     ratingNumeric,
		GlobalGamesCount: pgtype.Int4{Int32: 0, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("update user NTRP level: %w", err)
	}

	return &QuizResult{
		NTRPLevel:     ntrpLevel,
		InitialRating: initialRating,
	}, nil
}

// calculateNTRPLevel maps average weight to NTRP level
func calculateNTRPLevel(avgWeight float64) float64 {
	switch {
	case avgWeight < 1.3:
		return 1.0 // Beginner
	case avgWeight < 1.7:
		return 1.5
	case avgWeight < 2.2:
		return 2.0 // Novice
	case avgWeight < 2.7:
		return 2.5
	case avgWeight < 3.2:
		return 3.0 // Intermediate
	case avgWeight < 3.7:
		return 3.5
	case avgWeight < 4.2:
		return 4.0 // Advanced
	case avgWeight < 4.7:
		return 4.5
	default:
		return 5.0 // Professional
	}
}

// ntrpLevelToRating maps NTRP level to initial rating
func ntrpLevelToRating(ntrp float64) int {
	// Linear mapping: NTRP 1.0 = 100, NTRP 5.0 = 1000
	// Formula: rating = 100 + (ntrp - 1.0) * 225
	rating := 100 + int((ntrp-1.0)*225)

	// Clamp to valid range
	if rating < 100 {
		return 100
	}
	if rating > 1000 {
		return 1000
	}

	return rating
}
