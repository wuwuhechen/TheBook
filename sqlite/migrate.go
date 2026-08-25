package sqlite

import (
	"TheBook/utils"
	"fmt"
)

func Migrate() {
	root, err := utils.FindProjectRoot()
	if err != nil {
		panic(err)
	}

	db, err := openDatabase(fmt.Sprintf("%s/database/test.db", root))
	if err != nil {
		panic(err)
	}

	err = createDatabase(db, []any{
		&Question{},
		&User{},
		&PracticeRecord{},
		&PracticeAnswer{},
		&QuestionProgress{},
	})
	if err != nil {
		panic("failed to create database")
	}

	err = migrateDatabase(db, fmt.Sprintf("%s/database/", root))
	if err != nil {
		panic("failed to migrate database" + err.Error())
	}
	fmt.Println("Data migration completed successfully.")
}
