package main

import (
	"fmt"
	"os"
)

func main() {
	// 打印当前路径
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	fmt.Println("Current directory:", path)

	db, err := openDatabase(fmt.Sprintf("%s/../test.db", path))
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

	err = migrateDatabase(db, fmt.Sprintf("%s/../", path))
	if err != nil {
		panic("failed to migrate database" + err.Error())
	}
	fmt.Println("Data migration completed successfully.")
}
