package bank

import (
	"TheBook/model/manager"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QuestionProgressBankSQLite struct {
	questionProgressBank *gorm.DB

	mu sync.RWMutex
}

var _ manager.QuestionProgressManager = (*QuestionProgressBankSQLite)(nil)

func NewQuestionProgressBankSQLite(db *gorm.DB) *QuestionProgressBankSQLite {
	return &QuestionProgressBankSQLite{
		questionProgressBank: db,
	}
}

func (b *QuestionProgressBankSQLite) FindByUserID(userID uint) (*QuestionProgress, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var progress QuestionProgress
	err := b.questionProgressBank.Where("user_id = ?", userID).First(&progress).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			progress = *NewQuestionProgress(userID)
			return &progress, nil
		}
		return nil, err
	}

	return &progress, nil
}

func (b *QuestionProgressBankSQLite) Upsert(progress *QuestionProgress) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.questionProgressBank.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_id"},
			},
			UpdateAll: true,
		}).
		Create(progress).Error
}

func (b *QuestionProgressBankSQLite) Persist() error {
	// SQLite automatically persists data, so no action is needed here.
	return nil
}
