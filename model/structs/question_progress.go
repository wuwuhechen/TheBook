package structs

import "time"

type QuestionProgress struct {
	UserID            uint      `json:"user_id"`
	CurrentQuestionID int       `json:"current_question_id"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// NewQuestionProgress 创建指定用户的初始答题进度。
func NewQuestionProgress(userID uint) *QuestionProgress {
	return &QuestionProgress{
		UserID:            userID,
		CurrentQuestionID: 0,
		UpdatedAt:         time.Now(),
	}
}
