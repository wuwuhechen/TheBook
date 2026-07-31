package service

type QuestionManager interface {
	GetQuestion(id int) (*Question, error)
}

var _ QuestionManager = (*QuestionBank)(nil)

type QuestionServer struct {
	DB *QuestionBank
}
