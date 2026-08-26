package sqlite

import (
	"TheBook/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AddUser(db *gorm.DB, user *User) error {
	// UserID 0 is a valid legacy user ID. GORM normally treats a zero-valued
	// primary key as auto-generated and omits it from INSERT, so select the
	// columns explicitly to preserve IDs during JSON-to-SQLite migration.
	return db.Model(&User{}).Clauses(clause.OnConflict{UpdateAll: true}).Create(map[string]any{
		"user_id":  user.UserID,
		"username": user.Username,
		"password": user.Password,
		"nickname": user.Nickname,
		"role":     user.Role,
	}).Error
}

func AddUsers(db *gorm.DB, users []User) error {
	for i := range users {
		if err := AddUser(db, &users[i]); err != nil {
			return err
		}
	}
	return nil
}

func AddQuestion(db *gorm.DB, question *Question) error {
	return db.Clauses(clause.OnConflict{UpdateAll: true}).Create(question).Error
}

func AddQuestions(db *gorm.DB, questions []Question) error {
	for i := range questions {
		if err := AddQuestion(db, &questions[i]); err != nil {
			return err
		}
	}
	return nil
}

// AddPracticeRecord 将 JSON 中的一次套题记录及其逐题答案写入 SQLite。
// 父记录和答案记录在同一事务中创建，避免只写入其中一部分数据。
func AddPracticeRecord(db *gorm.DB, source *model.PracticeRecord) error {
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

// AddPracticeRecords 批量迁移 JSON 中的套题记录。
func AddPracticeRecords(db *gorm.DB, records []model.PracticeRecord) error {
	for i := range records {
		if err := AddPracticeRecord(db, &records[i]); err != nil {
			return err
		}
	}
	return nil
}

func AddQuestionProgress(db *gorm.DB, progress *QuestionProgress) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(progress).Error
}

func AddQuestionProgresses(db *gorm.DB, progresses []QuestionProgress) error {
	for i := range progresses {
		if err := AddQuestionProgress(db, &progresses[i]); err != nil {
			return err
		}
	}
	return nil
}
