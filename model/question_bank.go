package model

import (
	"fmt"
	"math/rand"
	"sort"
)

// QuestionBank 按真实题目 ID 存储题目。
type QuestionBank struct {
	Questions map[int]Question

	IDs []int
}

// NewQuestionBank 创建空的 QuestionBank。
func NewQuestionBank() *QuestionBank {
	return &QuestionBank{
		Questions: make(map[int]Question),
		IDs:       []int{},
	}
}

// GetQuestion 返回 id 对应的题目。
func (qb *QuestionBank) GetQuestion(id int) (*Question, error) {
	question, exists := qb.Questions[id]
	if !exists {
		return nil, fmt.Errorf("question with ID %d not found", id)
	}
	return &question, nil
}

// AddQuestion 按题目 ID 添加或替换题目。
func (qb *QuestionBank) AddQuestion(question Question) {
	qb.Questions[question.ID] = question
	qb.IDs = append(qb.IDs, question.ID)
}

// GetAllQuestions 返回所有题目，顺序不保证。
func (qb *QuestionBank) GetAllQuestions() []Question {
	questions := make([]Question, 0, len(qb.Questions))
	for _, question := range qb.Questions {
		questions = append(questions, question)
	}
	return questions
}

// GetAllQuestionsSorted 按题目 ID 顺序返回所有题目。
func (qb *QuestionBank) GetAllQuestionsSorted() []Question {
	sort.Ints(qb.IDs)

	questions := make([]Question, 0, len(qb.Questions))
	for _, questionID := range qb.IDs {
		question := qb.Questions[questionID]
		questions = append(questions, question)
	}

	return questions
}

// GetRandomQuestionID 返回随机选择的一道题目。
func (qb *QuestionBank) GetRandomQuestionID() (*Question, error) {
	if len(qb.IDs) == 0 {
		return nil, fmt.Errorf("no questions available")
	}

	randomIndex := rand.Intn(len(qb.IDs))
	randomQuestionID := qb.IDs[randomIndex]
	question, exists := qb.Questions[randomQuestionID]
	if !exists {
		return nil, fmt.Errorf("question with ID %d not found", randomQuestionID)
	}

	return &question, nil
}

// GetTotalCount 返回题库中的题目数量。
func (qb *QuestionBank) GetTotalCount() int {
	return len(qb.Questions)
}

// GetALLQuestionIDs 返回题库中所有题目的 ID。
func (qb *QuestionBank) GetALLQuestionIDs() []int {
	ids := make([]int, len(qb.IDs))
	copy(ids, qb.IDs)
	return ids
}
