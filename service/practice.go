package service

import (
	"math/rand"
	"time"
)

// 套题结构体
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

// NewPractice 创建一个新的套题实例
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

/*
GenerateExam 生成一个新的套题

参数

	qm QuestionManager: 问题管理器接口，用于获取问题数据
	size int: 套题的大小，即包含的问题数量

返回值

	*Practice: 生成的套题实例
*/
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

// GetCurrentQuestionID 返回套题的当前题目ID
func (p *Practice) GetCurrentQuestionID() int {
	if p.CurrentIndex < len(p.Questions) {
		return p.Questions[p.CurrentIndex]
	}
	return p.TotalQuestions // 返回题库总数表示没有更多问题
}

// NextQuestionID 返回套题的下一个题目ID
func (p *Practice) NextQuestionID() int {
	if p.CurrentIndex < len(p.Questions)-1 {
		return p.Questions[p.CurrentIndex+1]
	}
	return p.TotalQuestions
}

// LastQuestionID 返回套题的上一个题目ID
func (p *Practice) LastQuestionID() int {
	if p.CurrentIndex > 0 {
		return p.Questions[p.CurrentIndex-1]
	}
	return 0
}

func (p *Practice) Reset() {
	p.CurrentIndex = 0
	p.Answers = make(map[int]int)
	p.StartTime = time.Now()
	p.Completed = false
}

func (p *Practice) CheckPractice(qs *QuestionServer) *PracticeResponse {
	Total := len(p.Questions)
	CorrectCount := 0
	WrongCount := 0
	Details := make([]QuestionResponse, 0, Total)

	for _, id := range p.Questions {
		question, err := qs.DB.GetQuestion(id)
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

func (p *Practice) GetDuration() time.Duration {
	return p.Duration
}
