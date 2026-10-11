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

// TableName 返回 SQLite 题目表名称。
func (sqliteQuestionRow) TableName() string {
	return "questions"
}

type SQLiteQuestionBank struct {
	db *gorm.DB
}

var _ manager.QuestionManager = (*SQLiteQuestionBank)(nil)

// NewSQLiteQuestionBank 创建以 SQLite 为数据源的题库。
func NewSQLiteQuestionBank(db *gorm.DB) *SQLiteQuestionBank {
	return &SQLiteQuestionBank{db: db}
}

// GetQuestion 按题目 ID 从 SQLite 读取题目。
func (sqb *SQLiteQuestionBank) GetQuestion(id int) (*Question, error) {
	var question sqliteQuestionRow
	result := sqb.db.First(&question, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return convertSQLiteQuestion(&question), nil
}

// convertSQLiteQuestion 将 SQLite 行记录转换为领域题目。
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

// GetAllQuestionsSorted 按题目 ID 升序读取全部题目。
func (sqb *SQLiteQuestionBank) GetAllQuestionsSorted() []Question {
	var rows []sqliteQuestionRow
	sqb.db.Order("id").Find(&rows)
	return convertSQLiteQuestions(rows)
}

// GetAllQuestions 读取 SQLite 中的全部题目。
func (sqb *SQLiteQuestionBank) GetAllQuestions() []Question {
	var rows []sqliteQuestionRow
	sqb.db.Find(&rows)
	return convertSQLiteQuestions(rows)
}

// GetALLQuestionIDs 读取 SQLite 中的全部题目 ID。
func (sqb *SQLiteQuestionBank) GetALLQuestionIDs() []int {
	var ids []int
	sqb.db.Model(&sqliteQuestionRow{}).Pluck("id", &ids)
	return ids
}

// GetRandomQuestionID 随机读取 SQLite 中的一道题目。
func (sqb *SQLiteQuestionBank) GetRandomQuestionID() (*Question, error) {
	var question sqliteQuestionRow
	result := sqb.db.Order("RANDOM()").First(&question)
	if result.Error != nil {
		return nil, result.Error
	}
	return convertSQLiteQuestion(&question), nil
}

// GetTotalCount 返回 SQLite 题库中的题目数量。
func (sqb *SQLiteQuestionBank) GetTotalCount() int {
	var count int64
	sqb.db.Model(&sqliteQuestionRow{}).Count(&count)
	return int(count)
}

// convertSQLiteQuestions 将多条 SQLite 题目记录转换为领域题目。
func convertSQLiteQuestions(rows []sqliteQuestionRow) []Question {
	questions := make([]Question, 0, len(rows))
	for i := range rows {
		questions = append(questions, *convertSQLiteQuestion(&rows[i]))
	}
	return questions
}
