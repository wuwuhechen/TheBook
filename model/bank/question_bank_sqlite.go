package bank

import (
	"TheBook/model/manager"
	"gorm.io/gorm"
)

type sqliteQuestionRow struct {
	ID          uint `gorm:"primaryKey"`
	Category    string
	Question    string
	OptionA     string
	OptionB     string
	OptionC     string
	OptionD     string
	Answer      int
	Explanation string
}

func (sqliteQuestionRow) TableName() string {
	return "questions"
}

type SQLiteQuestionBank struct {
	db *gorm.DB
}

var _ manager.QuestionManager = (*SQLiteQuestionBank)(nil)

func NewSQLiteQuestionBank(db *gorm.DB) *SQLiteQuestionBank {
	return &SQLiteQuestionBank{db: db}
}

func (sqb *SQLiteQuestionBank) GetQuestion(id int) (*Question, error) {
	var question sqliteQuestionRow
	result := sqb.db.First(&question, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return convertSQLiteQuestion(&question), nil
}

func convertSQLiteQuestion(question *sqliteQuestionRow) *Question {
	return &Question{
		ID:          int(question.ID),
		Category:    question.Category,
		Question:    question.Question,
		Choices:     []string{question.OptionA, question.OptionB, question.OptionC, question.OptionD},
		Answer:      question.Answer,
		Explanation: question.Explanation,
	}
}

func (sqb *SQLiteQuestionBank) GetAllQuestionsSorted() []Question {
	var rows []sqliteQuestionRow
	sqb.db.Order("id").Find(&rows)
	return convertSQLiteQuestions(rows)
}

func (sqb *SQLiteQuestionBank) GetAllQuestions() []Question {
	var rows []sqliteQuestionRow
	sqb.db.Find(&rows)
	return convertSQLiteQuestions(rows)
}

func (sqb *SQLiteQuestionBank) GetALLQuestionIDs() []int {
	var ids []int
	sqb.db.Model(&sqliteQuestionRow{}).Pluck("id", &ids)
	return ids
}

func (sqb *SQLiteQuestionBank) GetRandomQuestionID() (*Question, error) {
	var question sqliteQuestionRow
	result := sqb.db.Order("RANDOM()").First(&question)
	if result.Error != nil {
		return nil, result.Error
	}
	return convertSQLiteQuestion(&question), nil
}

func (sqb *SQLiteQuestionBank) GetTotalCount() int {
	var count int64
	sqb.db.Model(&sqliteQuestionRow{}).Count(&count)
	return int(count)
}

func convertSQLiteQuestions(rows []sqliteQuestionRow) []Question {
	questions := make([]Question, 0, len(rows))
	for i := range rows {
		questions = append(questions, *convertSQLiteQuestion(&rows[i]))
	}
	return questions
}
