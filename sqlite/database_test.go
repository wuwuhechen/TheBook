package sqlite

import (
	"path/filepath"
	"testing"

	"gorm.io/gorm"
)

// closeTestDatabase 关闭测试数据库连接。
func closeTestDatabase(t *testing.T, db *gorm.DB) {
	t.Helper()
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get underlying database connection: %v", err)
	}
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Errorf("close test database: %v", err)
		}
	})
}

// TestOpenDatabaseAndCreateDatabase 验证数据库打开和模型迁移。
func TestOpenDatabaseAndCreateDatabase(t *testing.T) {
	db, err := OpenDatabase(filepath.Join(t.TempDir(), "thebook.db"))
	if err != nil {
		t.Fatalf("OpenDatabase() error = %v", err)
	}
	closeTestDatabase(t, db)

	models := []any{
		&Question{},
		&User{},
		&PracticeRecord{},
		&PracticeAnswer{},
		&QuestionProgress{},
		&WrongQuestion{},
	}
	if err := CreateDatabase(db, models); err != nil {
		t.Fatalf("CreateDatabase() error = %v", err)
	}

	for _, model := range models {
		if !db.Migrator().HasTable(model) {
			t.Errorf("expected table for %T to exist", model)
		}
	}
}
