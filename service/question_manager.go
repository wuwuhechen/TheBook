package service

// QuestionManager 定义了获取问题的接口
type QuestionManager interface {
	GetQuestion(id int) (*Question, error)
}

var _ QuestionManager = (*QuestionBank)(nil)

// QuestionServer 定义了问题服务器结构体
type QuestionServer struct {
	DB *QuestionBank
}
