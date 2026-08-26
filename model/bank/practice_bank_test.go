package bank

import (
	"TheBook/model/structs"
	"testing"

	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPracticeBankPersistUsesMigratedTables(t *testing.T) {
	db, err := gorm.Open(sqliteDriver.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	if err := db.AutoMigrate(&PracticeRecordRow{}, &PracticeAnswer{}); err != nil {
		t.Fatalf("Failed to migrate practice record tables: %v", err)
	}

	bank := NewPracticeBank(db)
	record := &structs.PracticeRecord{
		UserID:         2,
		PracticeID:     123,
		TotalQuestions: 1,
		CorrectCount:   1,
		Answers: []structs.AnswerRecord{{
			QuestionID: 7,
			Answer:     1,
			Answered:   true,
			Correct:    true,
		}},
	}
	if err := bank.Persist(record); err != nil {
		t.Fatalf("Failed to persist practice record: %v", err)
	}

	var storedRecord PracticeRecordRow
	if err := db.First(&storedRecord).Error; err != nil {
		t.Fatalf("Failed to load practice record: %v", err)
	}
	var storedAnswer PracticeAnswer
	if err := db.First(&storedAnswer).Error; err != nil {
		t.Fatalf("Failed to load practice answer: %v", err)
	}
	if storedAnswer.PracticeRecordID != storedRecord.ID {
		t.Fatalf("Expected answer record ID %d, got %d", storedRecord.ID, storedAnswer.PracticeRecordID)
	}
}
