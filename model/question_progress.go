package model

import "time"

type QuestionProgress struct {
	UserID            uint      `json:"user_id"`
	CurrentQuestionID int       `json:"current_question_id"`
	UpdatedAt         time.Time `json:"updated_at"`
}
