package bank

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
// NewQuestionBank 创建空的内存题库。
func NewQuestionBank() *QuestionBank {
	return &QuestionBank{
		Questions: make(map[int]Question),
		IDs:       []int{},
	}
}

// GetQuestion 返回 id 对应的题目。
// GetQuestion 按题目 ID 获取内存题库中的题目。
func (qb *QuestionBank) GetQuestion(id int) (*Question, error) {
	question, exists := qb.Questions[id]
	if !exists {
		return nil, fmt.Errorf("question with ID %d not found", id)
	}
	return &question, nil
}

// AddQuestion 按题目 ID 添加或替换题目。
// AddQuestion 添加或替换内存题库中的题目。
func (qb *QuestionBank) AddQuestion(question Question) {
	if _, exists := qb.Questions[question.ID]; !exists {
		qb.IDs = append(qb.IDs, question.ID)
	}
	qb.Questions[question.ID] = question
}

// GetAllQuestions 返回所有题目，顺序不保证。
// GetAllQuestions 返回题库中的全部题目副本。
func (qb *QuestionBank) GetAllQuestions() []Question {
	questions := make([]Question, 0, len(qb.Questions))
	for _, question := range qb.Questions {
		questions = append(questions, question)
	}
	return questions
}

// GetAllQuestionsSorted 按题目 ID 顺序返回所有题目。
// GetAllQuestionsSorted 按题目 ID 升序返回全部题目。
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
// GetRandomQuestionID 随机返回题库中的一道题目。
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
// GetTotalCount 返回题库中的题目数量。
func (qb *QuestionBank) GetTotalCount() int {
	return len(qb.Questions)
}

// GetALLQuestionIDs 返回题库中所有题目的 ID。
// GetALLQuestionIDs 返回题库中全部真实题目 ID。
func (qb *QuestionBank) GetALLQuestionIDs() []int {
	ids := make([]int, len(qb.IDs))
	copy(ids, qb.IDs)
	return ids
}
