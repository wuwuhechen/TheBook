package model

// QuestionManager 定义查询题库所需的操作。
type QuestionManager interface {
	GetQuestion(id int) (*Question, error)

	GetAllQuestionsSorted() []Question

	GetAllQuestions() []Question

	GetALLQuestionIDs() []int

	GetRandomQuestionID() (*Question, error)

	GetTotalCount() int
}

var _ QuestionManager = (*QuestionBank)(nil)

// PracticeManager 定义练习会话支持的操作。
type PracticeManager interface {
	GetCurrentQuestionID() int

	NextQuestionID() int

	LastQuestionID() int

	GenerateExam(qm QuestionManager, size int) *Practice

	CheckPractice(qm QuestionManager) *PracticeResponse

	Reset()
}

var _ PracticeManager = (*Practice)(nil)
