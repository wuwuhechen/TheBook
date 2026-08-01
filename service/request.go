package service

type Request struct {
	// UserID 字段用于存储用户的唯一标识符
	UserID string `json:"user_id"`

	// QuestionID 字段用于存储问题的唯一标识符
	QuestionID int `json:"question_id"`

	// Choice 字段用于存储用户选择的答案的索引
	Choice int `json:"choice"`
}
