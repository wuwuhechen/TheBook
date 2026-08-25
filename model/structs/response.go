package structs

// QuestionResponse 描述单题答案是否正确及其解析。
type QuestionResponse struct {
	// Correct 字段用于指示用户提交的答案是否正确
	Correct bool `json:"correct"`

	// Explanation 字段用于提供对用户提交答案的解释或反馈
	Explanation string `json:"explanation"`
}

// NewResponse 根据判题结果创建 QuestionResponse。
func NewResponse(correct bool, explanation string) *QuestionResponse {
	return &QuestionResponse{
		Correct:     correct,
		Explanation: explanation,
	}
}

// PracticeResponse 汇总一套练习的判题结果。
type PracticeResponse struct {
	PracticeID int `json:"practice_id"`

	Total int `json:"total"`

	CorrectCount int `json:"correct_count"`

	WrongCount int `json:"wrong_count"`

	Details []QuestionResponse `json:"details"`
}
