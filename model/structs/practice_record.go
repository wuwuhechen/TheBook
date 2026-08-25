package structs

import "time"

// PracticeRecord 保存用户一次套题提交后的汇总结果和逐题答题情况。
type PracticeRecord struct {
	ID             int            `json:"id"`
	UserID         uint           `json:"user_id"`
	PracticeID     int            `json:"practice_id"`
	TotalQuestions int            `json:"total_questions"`
	CorrectCount   int            `json:"correct_count"`
	WrongCount     int            `json:"wrong_count"`
	StartTime      time.Time      `json:"start_time"`
	SubmitTime     time.Time      `json:"submit_time"`
	Answers        []AnswerRecord `json:"answers"`
}

// AnswerRecord 保存一道题的原题 ID、用户答案和判题结果。
type AnswerRecord struct {
	QuestionID int  `json:"question_id"`
	Answer     int  `json:"answer"`
	Answered   bool `json:"answered"`
	Correct    bool `json:"correct"`
}
