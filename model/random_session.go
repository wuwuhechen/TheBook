package model

import (
	"math/rand"
	"time"
)

// RandomSession 保存一次随机答题的固定题目顺序和当前位置。
type RandomSession struct {
	ID           int
	Questions    []int
	CurrentIndex int
}

// NewRandomSession 创建包含全部题目的随机答题会话。
func NewRandomSession(qm QuestionManager) *RandomSession {
	questions := qm.GetALLQuestionIDs()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	return &RandomSession{
		ID:        int(random.Int63n(9000000000000000)),
		Questions: questions,
	}
}

// CurrentQuestionID 返回当前随机题目的真实 ID。
func (s *RandomSession) CurrentQuestionID() int {
	if s.CurrentIndex < 0 || s.CurrentIndex >= len(s.Questions) {
		return 0
	}
	return s.Questions[s.CurrentIndex]
}
