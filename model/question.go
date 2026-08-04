package model

import "fmt"

// Question 表示一道选择题及其正确答案。
type Question struct {
	// ID 字段用于存储问题的唯一标识符
	ID int `json:"id"`
	// Category 字段用于存储问题的分类信息
	Category string `json:"category"`
	// Question 字段用于存储问题的文本内容
	Question string `json:"question"`
	// Choices 字段用于存储问题的选项列表
	Choices []string `json:"choices"`
	// Answer 字段用于存储问题的正确答案的索引（从0开始）
	Answer int `json:"answer"`
	// Explanation 字段用于存储问题的解析信息
	Explanation string `json:"explanation"`
}

// CheckAnswer 报告 choice 是否为 q 的正确答案。
func (q *Question) CheckAnswer(choice int) bool {
	return choice == q.Answer
}

// String 返回 Question 的可读字符串表示。
func (q *Question) String() string {
	return fmt.Sprintf("编号: %d\n 分类: %s\n 问题: %s\n 选项: %v\n 答案: %d\n 解析: %s\n",
		q.ID, q.Category, q.Question, q.Choices, q.Answer, q.Explanation)
}
