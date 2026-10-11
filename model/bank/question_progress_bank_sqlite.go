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

// NewQuestionProgressBankSQLite 创建 SQLite 版用户答题进度仓库。
func NewQuestionProgressBankSQLite(db *gorm.DB) *QuestionProgressBankSQLite {
	return &QuestionProgressBankSQLite{
		questionProgressBank: db,
	}
}

// FindByUserID 从 SQLite 查询用户进度，不存在时返回初始进度。
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

// Upsert 在 SQLite 中创建或更新用户答题进度。
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

// Persist 保持接口兼容；SQLite 写入已在 Upsert 中完成。
func (b *QuestionProgressBankSQLite) Persist() error {
	// SQLite automatically persists data, so no action is needed here.
	return nil
}
