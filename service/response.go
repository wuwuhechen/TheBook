package service

// QuestionResponse 结构体用于表示用户提交答案后的响应结果
type QuestionResponse struct {
	// Correct 字段用于指示用户提交的答案是否正确
	Correct bool `json:"correct"`

	// Explanation 字段用于提供对用户提交答案的解释或反馈
	Explanation string `json:"explanation"`
}

func NewResponse(correct bool, explanation string) *QuestionResponse {
	return &QuestionResponse{
		Correct:     correct,
		Explanation: explanation,
	}
}

type PracticeResponse struct {
	Total int `json:"total"`

	CorrectCount int `json:"correct_count"`

	WrongCount int `json:"wrong_count"`

	Details []QuestionResponse `json:"details"`
}
