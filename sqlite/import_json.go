package sqlite

import (
	"TheBook/model"
	"encoding/json"
	"fmt"
	"os"

	"gorm.io/gorm"
)

func readUsers(path string) ([]User, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []User
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func readQuestions(path string) ([]Question, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var que []model.Question
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&que); err != nil {
		return nil, err
	}

	var questions []Question
	for _, q := range que {
		questions = append(questions, Question{
			ID:          uint(q.ID),
			Category:    q.Category,
			Question:    q.Question,
			OptionA:     q.Choices[0],
			OptionB:     q.Choices[1],
			OptionC:     q.Choices[2],
			OptionD:     q.Choices[3],
			Answer:      q.Answer,
			Explanation: q.Explanation,
		})
	}

	return questions, nil
}

// readPracticeRecords 读取旧 practice_records.json 的套题记录。
func readPracticeRecords(path string) ([]model.PracticeRecord, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []model.PracticeRecord
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&records); err != nil {
		return nil, err
	}

	return records, nil
}

func readQuestionProgress(path string) ([]QuestionProgress, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var progresses []QuestionProgress
	if err := json.NewDecoder(file).Decode(&progresses); err != nil {
		return nil, err
	}

	return progresses, nil
}

func migrateDatabase(db *gorm.DB, dataPath string) error {
	fmt.Println("Migrating data from JSON files to SQLite database...")

	questions, err := readQuestions(dataPath + "/data.json")
	fmt.Println("Read questions:", len(questions))
	if err != nil {
		return err
	}
	if err := addQuestions(db, questions); err != nil {
		return err
	}

	users, err := readUsers(dataPath + "/users.json")
	fmt.Println("Read users:", len(users))
	if err != nil {
		return err
	}
	if err := addUsers(db, users); err != nil {
		return err
	}

	practiceRecords, err := readPracticeRecords(dataPath + "/practice_records.json")
	fmt.Println("Read practice records:", len(practiceRecords))
	if err != nil {
		return err
	}
	if err := addPracticeRecords(db, practiceRecords); err != nil {
		return err
	}

	progresses, err := readQuestionProgress(dataPath + "/question_progress.json")
	fmt.Println("Read question progresses:", len(progresses))
	if err != nil {
		return err
	}
	if err := addQuestionProgresses(db, progresses); err != nil {
		return err
	}

	return nil
}
