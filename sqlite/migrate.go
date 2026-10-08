package sqlite

import (
	"TheBook/utils"
	"fmt"
	"os"
)

func Migrate(path string) error {
	root, err := utils.FindProjectRoot()
	if err != nil {
		return err
	}

	db_path := fmt.Sprintf("%s/database/%s", root, path)

	// 如果数据库文件已经存在，报错
	_, err = os.Open(db_path)
	if err == nil {
		return fmt.Errorf("database file already exists: %s", path)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check database file: %w", err)
	}

	// 创建数据库文件
	_, err = os.Create(db_path)
	if err != nil {
		return fmt.Errorf("failed to create database file: %w", err)
	}

	db, err := OpenDatabase(db_path)
	if err != nil {
		return err
	}

	err = CreateDatabase(db, []any{
		&Question{},
		&User{},
		&PracticeRecord{},
		&PracticeAnswer{},
		&QuestionProgress{},
		&WrongQuestion{},
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
