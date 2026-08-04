package service

// Request 包含客户端提交的题目与练习请求数据。
type Request struct {
	// UserID 字段用于存储用户的唯一标识符
	UserID string `json:"user_id" form:"user_id"`

	// QuestionID 字段用于存储问题的唯一标识符
	QuestionID int `json:"question_id" form:"question_id"`

	// Choice 字段用于存储用户选择的答案的索引
	Choice int `json:"choice" form:"choice"`

	// 套题大小字段，用于指定用户希望练习的题目数量
	PracticeSize int `json:"practice_size" form:"practice_size"`

	PracticeID int `json:"practice_id" form:"practice_id"`
}
