package bank

import (
	"testing"

	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newQuestionSQLiteTestBank 创建带测试题目的 SQLite 题库。
func newQuestionSQLiteTestBank(t *testing.T) *SQLiteQuestionBank {
	t.Helper()
	db, err := gorm.Open(sqliteDriver.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&sqliteQuestionRow{}); err != nil {
		t.Fatalf("migrate questions: %v", err)
	}
	rows := []sqliteQuestionRow{
		{ID: 2, Category: "Go", Question: "second", OptionA: "A", OptionB: "B", OptionC: "C", OptionD: "D", Answer: 1, Explanation: "two"},
		{ID: 1, Category: "Go", Question: "first", OptionA: "A1", OptionB: "B1", OptionC: "C1", OptionD: "D1", Answer: 0, Explanation: "one"},
	}
	if err := db.Create(&rows).Error; err != nil {
		t.Fatalf("seed questions: %v", err)
	}
	return NewSQLiteQuestionBank(db)
}

// TestSQLiteQuestionBankReadsAndSorts 验证 SQLite 题库查询与排序。
func TestSQLiteQuestionBankReadsAndSorts(t *testing.T) {
	bank := newQuestionSQLiteTestBank(t)
	question, err := bank.GetQuestion(2)
	if err != nil || question.Question != "second" || len(question.Choices) != 4 || question.Choices[1] != "B" {
		t.Fatalf("unexpected question: %+v, %v", question, err)
	}
	if _, err := bank.GetQuestion(99); err == nil {
		t.Fatal("expected missing question error")
	}
	sorted := bank.GetAllQuestionsSorted()
	if len(sorted) != 2 || sorted[0].ID != 1 || sorted[1].ID != 2 {
		t.Fatalf("unexpected sorted questions: %+v", sorted)
	}
	if bank.GetTotalCount() != 2 {
		t.Fatalf("unexpected total: %d", bank.GetTotalCount())
	}
	ids := bank.GetALLQuestionIDs()
	if len(ids) != 2 {
		t.Fatalf("unexpected IDs: %v", ids)
	}
}

// TestSQLiteQuestionBankRandomQuestion 验证 SQLite 题库随机取题。
func TestSQLiteQuestionBankRandomQuestion(t *testing.T) {
	bank := newQuestionSQLiteTestBank(t)
	question, err := bank.GetRandomQuestionID()
	if err != nil || (question.ID != 1 && question.ID != 2) {
		t.Fatalf("unexpected random question: %+v, %v", question, err)
	}
}
