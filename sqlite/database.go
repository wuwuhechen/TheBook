package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

func CreateDatabase(db *gorm.DB, models []any) error {
	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}
	return nil
}
