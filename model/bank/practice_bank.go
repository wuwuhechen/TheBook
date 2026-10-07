package bank

import (
	"TheBook/model/manager"
	"TheBook/model/structs"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

type PracticeBank struct {
	practices map[int]*Practice

	records *gorm.DB

	mu sync.Mutex
}

type PracticeAnswer struct {
	ID               uint `gorm:"primaryKey" json:"id"`
	PracticeRecordID uint `gorm:"not null;uniqueIndex:idx_practice_answer" json:"practice_record_id"`
	QuestionID       uint `gorm:"not null;uniqueIndex:idx_practice_answer" json:"question_id"`
	Answer           int  `gorm:"not null" json:"answer"`
	Answered         bool `gorm:"not null" json:"answered"`
	Correct          bool `gorm:"not null" json:"correct"`
}

type PracticeRecordRow struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null" json:"user_id"`

	PracticeID     uint `gorm:"not null" json:"practice_id"`
	TotalQuestions int  `gorm:"not null" json:"total_questions"`

	CorrectCount int `gorm:"not null" json:"correct_count"`
	WrongCount   int `gorm:"not null" json:"wrong_count"`

	StartTime  time.Time `gorm:"not null" json:"start_time"`
	SubmitTime time.Time `gorm:"not null" json:"submit_time"`
}

func NewPracticeBank(db *gorm.DB) *PracticeBank {
	return &PracticeBank{
		practices: make(map[int]*Practice),
		records:   db,
	}
}

func (pb *PracticeBank) Create(practice *Practice) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	if _, exists := pb.practices[practice.ID]; exists {
		return fmt.Errorf("practice with ID %d already exists", practice.ID)
	}
	pb.practices[practice.ID] = practice
	return nil
}

func (pb *PracticeBank) FindByID(id int) (*Practice, error) {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	practice, exists := pb.practices[id]
	if !exists {
		return nil, fmt.Errorf("practice with ID %d not found", id)
	}
	return practice, nil
}

func (pb *PracticeBank) Save(practice *Practice) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	if _, exists := pb.practices[practice.ID]; !exists {
		return fmt.Errorf("practice with ID %d does not exist", practice.ID)
	}
	pb.practices[practice.ID] = practice
	return nil
}

func (pb *PracticeBank) Delete(practice *Practice) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	if _, exists := pb.practices[practice.ID]; !exists {
		return fmt.Errorf("practice with ID %d does not exist", practice.ID)
	}
	delete(pb.practices, practice.ID)
	return nil
}

func (PracticeRecordRow) TableName() string {
	return "practice_records"
}

func (PracticeAnswer) TableName() string {
	return "practice_answers"
}

func convertToSQLiteStructs(record *PracticeRecord) (*PracticeRecordRow, []PracticeAnswer) {
	practiceRecordRow := &PracticeRecordRow{
		UserID:         record.UserID,
		PracticeID:     (uint)(record.PracticeID),
		TotalQuestions: record.TotalQuestions,
		CorrectCount:   record.CorrectCount,
		WrongCount:     record.WrongCount,
		StartTime:      record.StartTime,
		SubmitTime:     record.SubmitTime,
	}

	var practiceAnswers []PracticeAnswer
	for _, answer := range record.Answers {
		practiceAnswers = append(practiceAnswers, PracticeAnswer{
			QuestionID: (uint)(answer.QuestionID),
			Answer:     answer.Answer,
			Answered:   answer.Answered,
			Correct:    answer.Correct,
		})
	}

	return practiceRecordRow, practiceAnswers
}

func (pb *PracticeBank) Persist(record *PracticeRecord) error {
	practiceRecordRow, practiceAnswers := convertToSQLiteStructs(record)

	return pb.records.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(practiceRecordRow).Error; err != nil {
			return fmt.Errorf("failed to persist practice record: %w", err)
		}

		for i := range practiceAnswers {
			practiceAnswers[i].PracticeRecordID = practiceRecordRow.ID
		}
		if len(practiceAnswers) > 0 {
			if err := tx.Create(&practiceAnswers).Error; err != nil {
				return fmt.Errorf("failed to persist practice answers: %w", err)
			}
		}
		return nil
	})
}

func (pb *PracticeBank) ListRecordsByUser(userID uint, limit, offset int) ([]*PracticeRecord, error) {
	var records []PracticeRecordRow
	if err := pb.records.Where("user_id = ?", userID).Order("submit_time DESC").Limit(limit).Offset(offset).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to list practice records: %w", err)
	}

	var practiceRecords []*PracticeRecord
	for _, record := range records {
		practiceRecords = append(practiceRecords, &PracticeRecord{
			ID:             record.ID,
			UserID:         record.UserID,
			PracticeID:     (int)(record.PracticeID),
			TotalQuestions: record.TotalQuestions,
			CorrectCount:   record.CorrectCount,
			WrongCount:     record.WrongCount,
			StartTime:      record.StartTime,
			SubmitTime:     record.SubmitTime,
			Answers:        nil,
		})
	}
	return practiceRecords, nil
}

func convertToPracticeRecord(answers []manager.PracticeAnswer) []structs.AnswerRecord {
	var answerRecords []structs.AnswerRecord
	for _, answer := range answers {
		answerRecords = append(answerRecords, structs.AnswerRecord{
			QuestionID: answer.QuestionID,
			Answer:     answer.Answer,
			Answered:   answer.Answered,
			Correct:    answer.Correct,
		})
	}
	return answerRecords
}

func (pb *PracticeBank) GetRecordByPracticeID(userID uint, practiceID int) (*PracticeRecord, error) {
	var record PracticeRecordRow
	if err := pb.records.Where("user_id = ? AND practice_id = ?", userID, practiceID).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("practice record not found for user_id %d and practice_id %d", userID, practiceID)
		}
		return nil, fmt.Errorf("failed to get practice record: %w", err)
	}
	var answers []manager.PracticeAnswer
	if err := pb.records.Where("practice_record_id = ?", record.ID).Order("id ASC").Find(&answers).Error; err != nil {
		return nil, fmt.Errorf("failed to get practice answers: %w", err)
	}
	answersRecords := convertToPracticeRecord(answers)
	practicRecord := &PracticeRecord{
		ID:             record.ID,
		UserID:         record.UserID,
		PracticeID:     (int)(record.PracticeID),
		TotalQuestions: record.TotalQuestions,
		CorrectCount:   record.CorrectCount,
		WrongCount:     record.WrongCount,
		StartTime:      record.StartTime,
		SubmitTime:     record.SubmitTime,
		Answers:        answersRecords,
	}
	return practicRecord, nil
}

func (pb *PracticeBank) ListAnswers(recordID int) ([]manager.PracticeAnswer, error) {
	var answers []manager.PracticeAnswer
	if err := pb.records.Where("practice_record_id = ?", recordID).Order("id ASC").Find(&answers).Error; err != nil {
		return nil, fmt.Errorf("failed to list practice answers: %w", err)
	}
	return answers, nil
}
