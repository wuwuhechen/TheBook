package service

import "github.com/gin-gonic/gin"

type Server struct {
	Router *gin.Engine
	QS     *QuestionServer
	US     UserServer
}
