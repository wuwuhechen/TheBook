package sqlite

import (
	"TheBook/utils"
	"fmt"
)

func Migrate() error {
	root, err := utils.FindProjectRoot()
	if err != nil {
		return err
	}

	db, err := OpenDatabase(fmt.Sprintf("%s/database/thebook.db", root))
	if err != nil {
		return err
	}

	err = CreateDatabase(db, []any{
		&Question{},
		&User{},
		&PracticeRecord{},
		&PracticeAnswer{},
		&QuestionProgress{},
	})
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	err = MigrateDatabase(db, fmt.Sprintf("%s/database/", root))
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	fmt.Println("Data migration completed successfully.")

	return nil
}
