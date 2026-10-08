package bank

import (
	"TheBook/model/manager"
	"TheBook/model/structs"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// WrongQuestionRow 是 wrong_questions 表的 SQLite 映射。
type WrongQuestionRow struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;uniqueIndex:idx_wrong_question" json:"user_id"`
	QuestionID   uint      `gorm:"not null;uniqueIndex:idx_wrong_question" json:"question_id"`
	WrongCount   int       `gorm:"not null" json:"wrong_count"`
	FirstWrongAt time.Time `gorm:"not null" json:"first_wrong_at"`
	LastWrongAt  time.Time `gorm:"not null" json:"last_wrong_at"`
}

func (WrongQuestionRow) TableName() string {
	return "wrong_questions"
}

// WrongQuestionBank 负责错题本持久化，不依赖练习记录实现。
type WrongQuestionBank struct {
	db *gorm.DB
	mu sync.Mutex
}

var _ manager.WrongQuestionManager = (*WrongQuestionBank)(nil)

func NewWrongQuestionBank(db *gorm.DB) *WrongQuestionBank {
	return &WrongQuestionBank{db: db}
}

// RecordWrongQuestion 新增或累计一次错误记录。
func (b *WrongQuestionBank) RecordWrongQuestion(userID, questionID uint, wrongAt time.Time) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	row := WrongQuestionRow{
		UserID:       userID,
		QuestionID:   questionID,
		WrongCount:   1,
		FirstWrongAt: wrongAt,
		LastWrongAt:  wrongAt,
	}
	if err := b.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "question_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"wrong_count":   gorm.Expr("wrong_count + ?", 1),
			"last_wrong_at": wrongAt,
		}),
	}).Create(&row).Error; err != nil {
		return fmt.Errorf("failed to record wrong question: %w", err)
	}
	return nil
}

// ListWrongQuestionsByUser 按最近错误时间倒序读取指定用户的错题本。
func (b *WrongQuestionBank) ListWrongQuestionsByUser(userID uint, limit, offset int) ([]*structs.WrongQuestion, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	var rows []WrongQuestionRow
	if err := b.db.Where("user_id = ?", userID).
		Order("last_wrong_at DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to list wrong questions: %w", err)
	}

	questions := make([]*structs.WrongQuestion, 0, len(rows))
	for _, row := range rows {
		questions = append(questions, &structs.WrongQuestion{
			UserID: row.UserID, QuestionID: row.QuestionID, WrongCount: row.WrongCount,
			FirstWrongAt: row.FirstWrongAt, LastWrongAt: row.LastWrongAt,
		})
	}
	return questions, nil
}

// ListWrongQuestionsByIDs 返回指定用户错题本内的指定题目。
func (b *WrongQuestionBank) ListWrongQuestionsByIDs(userID uint, questionIDs []uint) ([]*structs.WrongQuestion, error) {
	if len(questionIDs) == 0 {
		return []*structs.WrongQuestion{}, nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	var rows []WrongQuestionRow
	if err := b.db.Where("user_id = ? AND question_id IN ?", userID, questionIDs).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to list selected wrong questions: %w", err)
	}

	questions := make([]*structs.WrongQuestion, 0, len(rows))
	for _, row := range rows {
		questions = append(questions, &structs.WrongQuestion{
			UserID:       row.UserID,
			QuestionID:   row.QuestionID,
			WrongCount:   row.WrongCount,
			FirstWrongAt: row.FirstWrongAt,
			LastWrongAt:  row.LastWrongAt,
		})
	}
	return questions, nil
}
