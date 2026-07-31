package service

import (
	"fmt"
	"sort"
)

type QuestionBank struct {
	Questions map[int]Question

	IDs []int
}

func (qb *QuestionBank) GetQuestion(id int) (*Question, error) {
	question, exists := qb.Questions[id]
	if !exists {
		return nil, fmt.Errorf("question with ID %d not found", id)
	}
	return &question, nil
}

func (qb *QuestionBank) AddQuestion(question Question) {
	qb.Questions[question.ID] = question
	qb.IDs = append(qb.IDs, question.ID)
}

func (qb *QuestionBank) GetAllQuestions() []Question {
	questions := make([]Question, 0, len(qb.Questions))
	for _, question := range qb.Questions {
		questions = append(questions, question)
	}
	return questions
}

func (qb *QuestionBank) GetAllQuestionsSorted() []Question {
	sort.Ints(qb.IDs)

	questions := make([]Question, 0, len(qb.Questions))
	for _, questionID := range qb.IDs {
		question := qb.Questions[questionID]
		questions = append(questions, question)
	}

	return questions
}
