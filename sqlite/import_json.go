package sqlite

import (
	"TheBook/auth"
	"TheBook/model"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"gorm.io/gorm"
)

func ReadUsers(path string) ([]User, error) {
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

func ReadQuestions(path string) ([]Question, error) {
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

// ReadPracticeRecords 读取旧 practice_records.json 的套题记录。
func ReadPracticeRecords(path string) ([]model.PracticeRecord, error) {
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

func ReadQuestionProgress(path string) ([]QuestionProgress, error) {
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

func ReadWrongQuestions(path string) ([]WrongQuestion, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var wrongQuestions []WrongQuestion
	if err := json.NewDecoder(file).Decode(&wrongQuestions); err != nil {
		return nil, err
	}
	return wrongQuestions, nil
}

func MigrateDatabase(db *gorm.DB, dataPath string) error {
	fmt.Println("Migrating data from JSON files to SQLite database...")

	questions, err := ReadQuestions(dataPath + "/data.json")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("No questions data found, skipping questions migration.")
	} else if err != nil {
		return err
	} else {
		fmt.Println("Read questions:", len(questions))
		if err := AddQuestions(db, questions); err != nil {
			return err
		}
	}

	users, err := ReadUsers(dataPath + "/users.json")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("No users data found, skipping users migration.")
	} else if err != nil {
		return err
	} else {
		for i := range users {
			if isBcryptPassword(users[i].Password) {
				continue
			}
			passwordHash, err := auth.HashPassword(users[i].Password)
			if err != nil {
				return fmt.Errorf("hash password for %s: %w", users[i].Username, err)
			}
			users[i].Password = passwordHash
		}
		if err := AddUsers(db, users); err != nil {
			return err
		}
	}

	practiceRecords, err := ReadPracticeRecords(dataPath + "/practice_records.json")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("No practice records data found, skipping practice records migration.")
	} else if err != nil {
		return err
	} else {
		fmt.Println("Read practice records:", len(practiceRecords))
		if err := AddPracticeRecords(db, practiceRecords); err != nil {
			return err
		}
	}

	progresses, err := ReadQuestionProgress(dataPath + "/question_progress.json")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("No question progress data found, skipping question progress migration.")
	} else if err != nil {
		return err
	} else {
		fmt.Println("Read question progresses:", len(progresses))
		if err := AddQuestionProgresses(db, progresses); err != nil {
			return err
		}
	}

	wrongQuestions, err := ReadWrongQuestions(dataPath + "/wrong_questions.json")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("No wrong questions data found, skipping wrong questions migration.")
	} else if err != nil {
		return err
	} else {
		fmt.Println("Read wrong questions:", len(wrongQuestions))
		if err := AddWrongQuestions(db, wrongQuestions); err != nil {
			return err
		}

	}

	return nil
}

func isBcryptPassword(password string) bool {
	return len(password) >= 4 &&
		(password[:4] == "$2a$" || password[:4] == "$2b$" || password[:4] == "$2y$")
}
