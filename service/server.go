package service

import (
	"TheBook/logger"
	"TheBook/model"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	Router *gin.Engine

	DB model.QuestionManager
	PM model.PracticeManager
	RS model.RandomSessionManager
	UM model.UserManager
	QS model.QuestionProgressManager

	UserPath     string
	UserHashPath string
	RecordPath   string
	RecordMu     sync.Mutex
	RootPath     string
	Log          *logger.Logger
}

// businessLog 返回业务日志器；未注入日志器时返回空日志器，便于测试环境使用。
func (s *Server) businessLog() *zap.Logger {
	if s.Log == nil || s.Log.Business == nil {
		return zap.NewNop()
	}
	return s.Log.Business
}

// appLog 返回应用日志器；未注入日志器时返回空日志器，便于测试环境使用。
func (s *Server) appLog() *zap.Logger {
	if s.Log == nil || s.Log.App == nil {
		return zap.NewNop()
	}
	return s.Log.App
}
