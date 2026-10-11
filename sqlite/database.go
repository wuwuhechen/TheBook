package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// OpenDatabase 打开指定路径的 SQLite 数据库连接。
func OpenDatabase(path string) (*gorm.DB, error) {
	db, err := gorm.Open(
		sqlite.Open(path),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CreateDatabase 对给定模型执行自动迁移。
func CreateDatabase(db *gorm.DB, models []any) error {
	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}
	return nil
}
