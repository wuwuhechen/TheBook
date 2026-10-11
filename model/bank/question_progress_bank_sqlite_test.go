package bank

import (
	"testing"
	"time"

	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// progressTestRow mirrors the SQLite persistence schema. The domain
// QuestionProgress deliberately has no GORM tags, so the test must create the
// same unique user_id constraint that production migration creates.
type progressTestRow struct {
	ID                uint      `gorm:"primaryKey"`
	UserID            uint      `gorm:"not null;uniqueIndex"`
	CurrentQuestionID int       `gorm:"not null"`
	UpdatedAt         time.Time `gorm:"not null"`
}

// TableName 返回测试进度记录对应的数据表名称。
func (progressTestRow) TableName() string { return "question_progresses" }

// TestQuestionProgressBankSQLiteFindAndUpsert 验证 SQLite 进度查询和更新。
func TestQuestionProgressBankSQLiteFindAndUpsert(t *testing.T) {
	db, err := gorm.Open(sqliteDriver.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&progressTestRow{}); err != nil {
		t.Fatalf("migrate progress: %v", err)
	}
	bank := NewQuestionProgressBankSQLite(db)

	missing, err := bank.FindByUserID(5)
	if err != nil || missing.UserID != 5 || missing.CurrentQuestionID != 0 {
		t.Fatalf("unexpected missing progress: %+v, %v", missing, err)
	}
	if err := bank.Upsert(&QuestionProgress{UserID: 5, CurrentQuestionID: 7}); err != nil {
		t.Fatalf("first upsert: %v", err)
	}
	if err := bank.Upsert(&QuestionProgress{UserID: 5, CurrentQuestionID: 9}); err != nil {
		t.Fatalf("second upsert: %v", err)
	}
	stored, err := bank.FindByUserID(5)
	if err != nil || stored.CurrentQuestionID != 9 {
		t.Fatalf("unexpected persisted progress: %+v, %v", stored, err)
	}
	if err := bank.Persist(); err != nil {
		t.Fatalf("SQLite persist should be a no-op: %v", err)
	}
}
