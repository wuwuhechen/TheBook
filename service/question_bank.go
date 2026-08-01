package service

import (
	"fmt"
	"sort"
)

// QuestionBank 用于存储和管理问题的结构体
type QuestionBank struct {
	Questions map[int]Question

	IDs []int
}

/*
NewQuestionBank 创建一个新的问题库实例

返回值

	*QuestionBank: 新的问题库实例
*/
func NewQuestionBank() *QuestionBank {
	return &QuestionBank{
		Questions: make(map[int]Question),
		IDs:       []int{},
	}
}

/*
GetQuestion 根据问题ID获取问题实例

参数

	id int: 问题ID

返回值

	*Question: 问题实例
	error: 错误信息
*/
func (qb *QuestionBank) GetQuestion(id int) (*Question, error) {
	question, exists := qb.Questions[id]
	if !exists {
		return nil, fmt.Errorf("question with ID %d not found", id)
	}
	return &question, nil
}

/*
AddQuestion 向问题库中添加一个问题

参数

	question Question: 要添加的问题

返回值

	无
*/
func (qb *QuestionBank) AddQuestion(question Question) {
	qb.Questions[question.ID] = question
	qb.IDs = append(qb.IDs, question.ID)
}

/*
GetAllQuestions 返回问题库中所有的问题

参数

	无

返回值

	[]Question: 问题列表
*/
func (qb *QuestionBank) GetAllQuestions() []Question {
	questions := make([]Question, 0, len(qb.Questions))
	for _, question := range qb.Questions {
		questions = append(questions, question)
	}
	return questions
}

/*
GetAllQuestionsSorted 返回按ID排序的问题列表

参数

	无

返回值

	[]Question: 按ID排序的问题列表
*/
func (qb *QuestionBank) GetAllQuestionsSorted() []Question {
	sort.Ints(qb.IDs)

	questions := make([]Question, 0, len(qb.Questions))
	for _, questionID := range qb.IDs {
		question := qb.Questions[questionID]
		questions = append(questions, question)
	}

	return questions
}
