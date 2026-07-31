package service

import "fmt"

type QuestionPageData struct {
	ID             int
	Category       string
	Question       string
	OptionA        string
	OptionB        string
	OptionC        string
	OptionD        string
	Explanation    string
	TotalQuestions int
}

func NewQuestionData(question *Question, totalQuestions int) *QuestionPageData {
	return &QuestionPageData{
		ID:             question.ID,
		Category:       question.Category,
		Question:       question.Question,
		OptionA:        question.Choices[0],
		OptionB:        question.Choices[1],
		OptionC:        question.Choices[2],
		OptionD:        question.Choices[3],
		Explanation:    question.Explanation,
		TotalQuestions: totalQuestions,
	}
}

func (q *QuestionPageData) String() string {
	return fmt.Sprintf("编号：%d\n 分类：%s\n 问题：%s\n 选项A：%s\n 选项B：%s\n 选项C：%s\n 选项D：%s\n 解析：%s\n 总题数：%d",
		q.ID, q.Category, q.Question, q.OptionA, q.OptionB, q.OptionC, q.OptionD, q.Explanation, q.TotalQuestions)
}
