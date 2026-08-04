package service

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

	CheckPractice(qs *QuestionServer) *PracticeResponse

	Reset()
}

var _ PracticeManager = (*Practice)(nil)

// QuestionServer 协调题库存储、练习状态和 HTTP 处理器。
type QuestionServer struct {
	DB QuestionManager

	PM map[int]*Practice

	RS map[int]*RandomSession
}
