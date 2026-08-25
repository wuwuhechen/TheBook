package model

import "gorm.io/gorm"

type SQLiteQuestionBank struct {
	db *gorm.DB
}

var _ QuestionManager = (*SQLiteQuestionBank)(nil)

func (sqb *SQLiteQuestionBank) GetQuestion(id int) (*Question, error) {
	var question Question
	result := sqb.db.First(&question, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &question, nil
}

func (sqb *SQLiteQuestionBank) GetAllQuestionsSorted() []Question {
	var questions []Question
	sqb.db.Order("id").Find(&questions)
	return questions
}

func (sqb *SQLiteQuestionBank) GetAllQuestions() []Question {
	var questions []Question
	sqb.db.Find(&questions)
	return questions
}

func (sqb *SQLiteQuestionBank) GetALLQuestionIDs() []int {
	var ids []int
	sqb.db.Model(&Question{}).Pluck("id", &ids)
	return ids
}

func (sqb *SQLiteQuestionBank) GetRandomQuestionID() (*Question, error) {
	var question Question
	result := sqb.db.Order("RANDOM()").First(&question)
	if result.Error != nil {
		return nil, result.Error
	}
	return &question, nil
}

func (sqb *SQLiteQuestionBank) GetTotalCount() int {
	var count int64
	sqb.db.Model(&Question{}).Count(&count)
	return int(count)
}
