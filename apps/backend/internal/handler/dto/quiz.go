package dto

import "github.com/alma-amirseitov/Tennis-App/apps/backend/internal/service"

// QuizQuestionsResponse represents response with quiz questions
type QuizQuestionsResponse struct {
	Questions []service.Question `json:"questions"`
}

// QuizAnswer represents a single answer in quiz submission
type QuizAnswer struct {
	QuestionID string `json:"question_id" validate:"required"`
	AnswerID   string `json:"answer_id" validate:"required"`
}

// QuizSubmitRequest represents request to submit quiz answers
type QuizSubmitRequest struct {
	Answers []QuizAnswer `json:"answers" validate:"required,min=1,dive"`
}

// QuizSubmitResponse represents response after quiz submission
type QuizSubmitResponse struct {
	NTRPLevel     float64 `json:"ntrp_level"`
	InitialRating int     `json:"initial_rating"`
}

// ToServiceAnswers converts DTOs to service layer answers
func (r *QuizSubmitRequest) ToServiceAnswers() []service.Answer {
	answers := make([]service.Answer, len(r.Answers))
	for i, a := range r.Answers {
		answers[i] = service.Answer{
			QuestionID: a.QuestionID,
			AnswerID:   a.AnswerID,
		}
	}
	return answers
}
