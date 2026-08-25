package structs

// QuestionManager is the minimal question source required by domain behavior.
// The public manager package aliases this contract to keep dependencies acyclic.
type QuestionManager interface {
	GetQuestion(id int) (*Question, error)
	GetAllQuestionsSorted() []Question
	GetAllQuestions() []Question
	GetALLQuestionIDs() []int
	GetRandomQuestionID() (*Question, error)
	GetTotalCount() int
}
