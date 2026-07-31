package service

import "fmt"

type Question struct {
	ID          int      `json:"id"`
	Category    string   `json:"category"`
	Question    string   `json:"question"`
	Choices     []string `json:"choices"`
	Answer      int      `json:"answer"`
	Explanation string   `json:"explanation"`
}

func (q *Question) CheckAnswer(choice int) bool {
	return choice == q.Answer
}

func (q *Question) String() string {
	return fmt.Sprintf("编号: %d\n 分类: %s\n 问题: %s\n 选项: %v\n 答案: %d\n 解析: %s\n",
		q.ID, q.Category, q.Question, q.Choices, q.Answer, q.Explanation)
}
