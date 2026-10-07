package manager

import (
	"TheBook/model/structs"
)

type PracticeManager interface {
	Create(practice *structs.Practice) error

	FindByID(id int) (*structs.Practice, error)

	Save(practice *structs.Practice) error

	Delete(practice *structs.Practice) error

	Persist(record *structs.PracticeRecord) error

	ListRecordsByUser(userID uint, limit, offset int) ([]*structs.PracticeRecord, error)

	GetRecordByPracticeID(userID uint, practiceID int) (*structs.PracticeRecord, error)

	ListAnswers(recordID int) ([]PracticeAnswer, error)
}

type PracticeAnswer struct {
	ID               uint `gorm:"primaryKey" json:"id"`
	PracticeRecordID uint `gorm:"not null;uniqueIndex:idx_practice_answer" json:"practice_record_id"`
	QuestionID       uint `gorm:"not null;uniqueIndex:idx_practice_answer" json:"question_id"`
	Answer           int  `gorm:"not null" json:"answer"`
	Answered         bool `gorm:"not null" json:"answered"`
	Correct          bool `gorm:"not null" json:"correct"`
}
