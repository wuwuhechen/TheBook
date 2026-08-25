package model

import (
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

var _ QuestionProgressManager = (*QuestionProgressBank)(nil)

func NewQuestionProgressBank() *QuestionProgressBank {
	return &QuestionProgressBank{
		questionProgresses: make(map[uint]*QuestionProgress),
	}
}

func (b *QuestionProgressBank) FindByUserID(userID uint) (*QuestionProgress, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	progress, ok := b.questionProgresses[userID]
	if !ok {
		progress = NewQuestionProgress(userID)
	}
	return progress, nil
}

func (b *QuestionProgressBank) Upsert(progress *QuestionProgress) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.questionProgresses[progress.UserID] = progress
	return nil
}

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
