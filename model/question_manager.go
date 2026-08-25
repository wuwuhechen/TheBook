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
