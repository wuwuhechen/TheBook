package main

import (
	"TheBook/model"
	"encoding/json"
	"fmt"
	"os"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func addUser(db *gorm.DB, user *User) error {
	return db.Clauses(clause.OnConflict{UpdateAll: true}).Create(user).Error
}

func addUsers(db *gorm.DB, users []User) error {
	for i := range users {
		if err := addUser(db, &users[i]); err != nil {
			return err
		}
	}
	return nil
}

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

func addQuestion(db *gorm.DB, question *Question) error {
	return db.Clauses(clause.OnConflict{UpdateAll: true}).Create(question).Error
}

func addQuestions(db *gorm.DB, questions []Question) error {
	for i := range questions {
		if err := addQuestion(db, &questions[i]); err != nil {
			return err
		}
	}
	return nil
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

// addPracticeRecord 将 JSON 中的一次套题记录及其逐题答案写入 SQLite。
// 父记录和答案记录在同一事务中创建，避免只写入其中一部分数据。
func addPracticeRecord(db *gorm.DB, source *model.PracticeRecord) error {
	return db.Transaction(func(tx *gorm.DB) error {
		record := PracticeRecord{
			ID:             uint(source.ID),
			UserID:         source.UserID,
			PracticeID:     uint(source.PracticeID),
			TotalQuestions: source.TotalQuestions,
			CorrectCount:   source.CorrectCount,
			WrongCount:     source.WrongCount,
			StartTime:      source.StartTime,
			SubmitTime:     source.SubmitTime,
		}

		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&record).Error; err != nil {
			return err
		}

		answers := make([]PracticeAnswer, 0, len(source.Answers))
		for _, sourceAnswer := range source.Answers {
			answers = append(answers, PracticeAnswer{
				PracticeRecordID: record.ID,
				QuestionID:       uint(sourceAnswer.QuestionID),
				Answer:           sourceAnswer.Answer,
				Answered:         sourceAnswer.Answered,
				Correct:          sourceAnswer.Correct,
			})
		}

		if len(answers) == 0 {
			return nil
		}
		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "practice_record_id"},
				{Name: "question_id"},
			},
			UpdateAll: true,
		}).Create(&answers).Error
	})
}

// addPracticeRecords 批量迁移 JSON 中的套题记录。
func addPracticeRecords(db *gorm.DB, records []model.PracticeRecord) error {
	for i := range records {
		if err := addPracticeRecord(db, &records[i]); err != nil {
			return err
		}
	}
	return nil
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

func addQuestionProgress(db *gorm.DB, progress *QuestionProgress) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(progress).Error
}

func addQuestionProgresses(db *gorm.DB, progresses []QuestionProgress) error {
	for i := range progresses {
		if err := addQuestionProgress(db, &progresses[i]); err != nil {
			return err
		}
	}
	return nil
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
