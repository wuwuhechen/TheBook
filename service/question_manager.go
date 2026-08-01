package service

import "github.com/gin-gonic/gin"

// QuestionManager 定义了获取问题的接口
type QuestionManager interface {
	GetQuestion(id int) (*Question, error)
}

var _ QuestionManager = (*QuestionBank)(nil)

// QuestionServer 定义了问题服务器结构体
type QuestionServer struct {
	DB *QuestionBank
}

// APP 定义了应用程序结构体，包含Gin引擎和问题服务器
type APP struct {
	// Gin引擎实例
	r *gin.Engine

	// QuestionServer 实例
	qs *QuestionServer
}
