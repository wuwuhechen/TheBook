package service

import "fmt"

// Question 定义了问题的结构体
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

/*
CheckAnswer 检查用户选择的答案是否正确

参数

	choice int: 用户选择的答案的索引

返回值

	bool: 如果答案正确则返回true，否则返回false
*/
func (q *Question) CheckAnswer(choice int) bool {
	return choice == q.Answer
}

func (q *Question) String() string {
	return fmt.Sprintf("编号: %d\n 分类: %s\n 问题: %s\n 选项: %v\n 答案: %d\n 解析: %s\n",
		q.ID, q.Category, q.Question, q.Choices, q.Answer, q.Explanation)
}
