package bank

import (
	"testing"
	"time"

	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestWrongQuestionBankRecordsAndListsQuestions(t *testing.T) {
	db, err := gorm.Open(sqliteDriver.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	if err := db.AutoMigrate(&WrongQuestionRow{}); err != nil {
		t.Fatalf("Failed to migrate wrong question table: %v", err)
	}

	bank := NewWrongQuestionBank(db)
	firstWrongAt := time.Date(2026, 10, 7, 10, 0, 0, 0, time.UTC)
	lastWrongAt := firstWrongAt.Add(time.Minute)
	if err := bank.RecordWrongQuestion(3, 11, firstWrongAt); err != nil {
		t.Fatalf("Failed to record first wrong answer: %v", err)
	}
	if err := bank.RecordWrongQuestion(3, 11, lastWrongAt); err != nil {
		t.Fatalf("Failed to record repeated wrong answer: %v", err)
	}
	if err := bank.RecordWrongQuestion(4, 12, lastWrongAt); err != nil {
		t.Fatalf("Failed to record another user's wrong answer: %v", err)
	}

	questions, err := bank.ListWrongQuestionsByUser(3, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list wrong questions: %v", err)
	}
	if len(questions) != 1 {
		t.Fatalf("Expected one wrong question, got %d", len(questions))
	}
	if questions[0].QuestionID != 11 || questions[0].WrongCount != 2 {
		t.Fatalf("Unexpected wrong question: %+v", questions[0])
	}
	if !questions[0].FirstWrongAt.Equal(firstWrongAt) || !questions[0].LastWrongAt.Equal(lastWrongAt) {
		t.Fatalf("Wrong timestamps: %+v", questions[0])
	}
}
