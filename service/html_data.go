package service

import (
	"fmt"
	"time"
)

// QuestionPageData 用于在HTML模板中渲染问题页面的数据结构
type QuestionPageData struct {
	// ID 字段用于存储问题的唯一标识符
	ID int

	// RealID 字段用于存储问题的实际唯一标识符
	RealID int

	// LastID 字段用于存储上一个问题的唯一标识符
	LastID int

	// HasLastID 字段用于指示是否存在上一个问题
	HasLastID bool

	// HasNextID 字段用于指示是否存在下一个问题
	HasNextID bool

	// NextID 字段用于存储下一个问题的唯一标识符
	NextID int

	// PracticeID 字段用于存储练习的唯一标识符
	PracticeID int

	// Category 字段用于存储问题的分类信息
	Category string

	// Question 字段用于存储问题的文本内容
	Question string

	// OptionA 字段用于存储问题的选项A
	OptionA string

	// OptionB 字段用于存储问题的选项B
	OptionB string

	// OptionC 字段用于存储问题的选项C
	OptionC string

	// OptionD 字段用于存储问题的选项D
	OptionD string

	// Explanation 字段用于存储问题的解析信息
	Explanation string

	// TotalQuestions 字段用于存储总题数
	TotalQuestions int

	// Duration 字段用于存储练习的持续时间（以秒为单位）
	Duration int
}

/*
NewQuestionData 创建新的问题页面数据

参数

	question *Question: 问题实例
	totalQuestions int: 总题数

返回值

	*QuestionPageData: 问题页面数据实例
*/
func NewQuestionData(question *Question, totalQuestions int) *QuestionPageData {
	return &QuestionPageData{
		ID:             question.ID,
		RealID:         question.ID,
		LastID:         question.ID - 1,
		HasLastID:      question.ID > 1,
		NextID:         question.ID + 1,
		HasNextID:      question.ID < totalQuestions,
		Category:       question.Category,
		Question:       question.Question,
		OptionA:        question.Choices[0],
		OptionB:        question.Choices[1],
		OptionC:        question.Choices[2],
		OptionD:        question.Choices[3],
		Explanation:    question.Explanation,
		TotalQuestions: totalQuestions,
		PracticeID:     0,
		Duration:       0,
	}
}

func (q *QuestionPageData) SetLastID(lastID int) {
	q.LastID = lastID
}

func (q *QuestionPageData) SetNextID(nextID int, totalQuestions int) {
	q.NextID = nextID
	q.TotalQuestions = totalQuestions
}

func (q *QuestionPageData) SetHasLastID(hasLastID bool) {
	q.HasLastID = hasLastID
}

func (q *QuestionPageData) SetHasNextID(hasNextID bool) {
	q.HasNextID = hasNextID
}

func (q *QuestionPageData) SetDuration(duration time.Duration) {
	q.Duration = int(duration.Seconds())
}

func (q *QuestionPageData) SetPracticeID(practiceID int) {
	q.PracticeID = practiceID
}

func (q *QuestionPageData) SetID(ID int) {
	q.ID = ID
}

/*
String 返回问题页面数据的字符串表示

参数

	无

返回值

	string: 问题页面数据的字符串表示
*/
func (q *QuestionPageData) String() string {
	return fmt.Sprintf("编号：%d\n 分类：%s\n 问题：%s\n 选项A：%s\n 选项B：%s\n 选项C：%s\n 选项D：%s\n 解析：%s\n 总题数：%d",
		q.ID, q.Category, q.Question, q.OptionA, q.OptionB, q.OptionC, q.OptionD, q.Explanation, q.TotalQuestions)
}

type PracticeResultItem struct {
	Number            int      `json:"number"`
	RealID            int      `json:"real_id"`
	Category          string   `json:"category"`
	Question          string   `json:"question"`
	Choices           []string `json:"choices"`
	UserAnswer        int      `json:"user_answer"`
	Answered          bool     `json:"answered"`
	UserAnswerText    string   `json:"user_answer_text"`
	CorrectAnswer     int      `json:"correct_answer"`
	CorrectAnswerText string   `json:"correct_answer_text"`
	Correct           bool     `json:"correct"`
	Explanation       string   `json:"explanation"`
}

type PracticeResultPageData struct {
	PracticeID   int                  `json:"practice_id"`
	Total        int                  `json:"total"`
	CorrectCount int                  `json:"correct_count"`
	WrongCount   int                  `json:"wrong_count"`
	Items        []PracticeResultItem `json:"Items"`
}
