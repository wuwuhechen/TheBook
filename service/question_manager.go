package service

// QuestionManager 定义了获取问题的接口
type QuestionManager interface {
	GetQuestion(id int) (*Question, error)

	GetAllQuestionsSorted() []Question

	GetAllQuestions() []Question

	GetALLQuestionIDs() []int

	GetRandomQuestionID() (*Question, error)

	GetTotalCount() int
}

var _ QuestionManager = (*QuestionBank)(nil)

type PracticeManager interface {
	GetCurrentQuestionID() int

	NextQuestionID() int

	LastQuestionID() int

	GenerateExam(qm QuestionManager, size int) *Practice

	CheckPractice(qs *QuestionServer) *PracticeResponse

	Reset()
}

var _ PracticeManager = (*Practice)(nil)

// QuestionServer 定义了问题服务器结构体
type QuestionServer struct {
	DB QuestionManager

	PM map[int]*Practice
}
