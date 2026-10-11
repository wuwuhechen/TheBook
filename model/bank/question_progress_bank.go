package bank

import (
	"TheBook/model/manager"
	"TheBook/utils"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type QuestionProgressBank struct {
	questionProgresses map[uint]*QuestionProgress

	mu sync.RWMutex
}

var _ manager.QuestionProgressManager = (*QuestionProgressBank)(nil)

// NewQuestionProgressBank 创建内存版用户答题进度仓库。
func NewQuestionProgressBank() *QuestionProgressBank {
	return &QuestionProgressBank{
		questionProgresses: make(map[uint]*QuestionProgress),
	}
}

// FindByUserID 查询用户进度，不存在时返回初始进度。
func (b *QuestionProgressBank) FindByUserID(userID uint) (*QuestionProgress, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	progress, ok := b.questionProgresses[userID]
	if !ok {
		progress = NewQuestionProgress(userID)
	}
	return progress, nil
}

// Upsert 创建或更新用户答题进度。
func (b *QuestionProgressBank) Upsert(progress *QuestionProgress) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.questionProgresses[progress.UserID] = progress
	return nil
}

// Persist 将内存中的用户答题进度写入 JSON 文件。
func (b *QuestionProgressBank) Persist() error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	progresses := make([]*QuestionProgress, 0, len(b.questionProgresses))
	for _, progress := range b.questionProgresses {
		progresses = append(progresses, progress)
	}

	data, err := json.MarshalIndent(progresses, "", "  ")
	if err != nil {
		return err
	}

	path, err := utils.FindProjectRoot()
	if err != nil {
		return err
	}
	path = path + "/database/question_progress.json"
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
