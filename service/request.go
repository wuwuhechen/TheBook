package service

type Request struct {
	UserID     string `json:"user_id"`
	QuestionID int    `json:"question_id"`
	Choice     int    `json:"choice"`
}
