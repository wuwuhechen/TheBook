package model

import (
	"math/rand"
	"time"
)

// 套题结构体
// Practice 保存一次练习的题目、答案和计时状态。
type Practice struct {
	// 套题ID
	ID int `json:"id"`

	// 题库总题数
	TotalQuestions int `json:"total_questions"`

	// 套题大小
	PracticeSize int `json:"practice_size"`

	// 现在题目
	CurrentIndex int `json:"current_index"`

	// 套题包含的问题列表
	Questions []int `json:"questions"`

	// 用户答案，键为题目ID，值为用户的答案
	Answers map[int]int `json:"answers"`

	// 开始时间
	StartTime time.Time `json:"start_time"`

	// 持续时间
	Duration time.Duration `json:"duration"`

	// 是否已完成
	Completed bool `json:"completed"`
}

// NewPractice 创建具有默认时限且题目和答案为空的练习。
func (p *Practice) NewPractice() *Practice {
	return &Practice{
		ID:             0,
		TotalQuestions: 0,
		CurrentIndex:   0,
		PracticeSize:   0,
		Questions:      []int{},
		Answers:        map[int]int{},
		StartTime:      time.Time{},
		Duration:       60 * time.Minute,
		Completed:      false,
	}
}

// GenerateExam 从 qm 中随机选取最多 size 道题并创建练习。
func (p *Practice) GenerateExam(qm QuestionManager, size int) *Practice {
	n := qm.GetTotalCount()

	size = min(size, n)

	ids := qm.GetALLQuestionIDs()

	r := rand.New(
		rand.NewSource(
			time.Now().UnixNano(),
		),
	)

	r.Shuffle(
		n,
		func(i, j int) {
			ids[i], ids[j] = ids[j], ids[i]
		},
	)

	return &Practice{
		ID:             int(r.Int63n(9000000000000000)),
		TotalQuestions: size,
		CurrentIndex:   0,
		PracticeSize:   size,
		Questions:      ids[:size],
		Answers:        map[int]int{},
		StartTime:      time.Now(),
		Duration:       60 * time.Minute,
		Completed:      false,
	}
}

// GetCurrentQuestionID 返回当前练习题目的真实 ID。
func (p *Practice) GetCurrentQuestionID() int {
	if p.CurrentIndex < len(p.Questions) {
		return p.Questions[p.CurrentIndex]
	}
	return p.TotalQuestions // 返回题库总数表示没有更多问题
}

// NextQuestionID 在存在下一题时前进并返回其 ID。
func (p *Practice) NextQuestionID() int {
	if p.CurrentIndex < len(p.Questions)-1 {
		return p.Questions[p.CurrentIndex+1]
	}
	return p.TotalQuestions
}

// LastQuestionID 在存在上一题时后退并返回其 ID。
func (p *Practice) LastQuestionID() int {
	if p.CurrentIndex > 0 {
		return p.Questions[p.CurrentIndex-1]
	}
	return 0
}

// Reset 清空答案并重置练习计时和导航状态。
func (p *Practice) Reset() {
	p.CurrentIndex = 0
	p.Answers = make(map[int]int)
	p.StartTime = time.Now()
	p.Completed = false
}

// CheckPractice 使用已保存的答案批改练习中的全部题目。
func (p *Practice) CheckPractice(qm QuestionManager) *PracticeResponse {
	Total := len(p.Questions)
	CorrectCount := 0
	WrongCount := 0
	Details := make([]QuestionResponse, 0, Total)

	for _, id := range p.Questions {
		question, err := qm.GetQuestion(id)
		if err != nil {
			continue
		}

		answer, ok := p.Answers[id]

		correct := ok && answer == question.Answer

		if correct {
			CorrectCount++
		} else {
			WrongCount++
		}

		Details = append(Details, QuestionResponse{
			Correct:     correct,
			Explanation: question.Explanation,
		})
	}

	return &PracticeResponse{
		PracticeID:   p.ID,
		Total:        Total,
		CorrectCount: CorrectCount,
		WrongCount:   WrongCount,
		Details:      Details,
	}
}

// GetDuration 返回练习配置的时限。
func (p *Practice) GetDuration() time.Duration {
	return p.Duration
}
