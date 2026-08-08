package service

import (
	"TheBook/model"
	"sync"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router       *gin.Engine
	DB           model.QuestionManager
	PM           map[int]*model.Practice
	RS           map[int]*model.RandomSession
	US           UserServer
	QS           map[uint]*model.QuestionProgress
	UserPath     string
	UserHashPath string
	RecordPath   string
	RecordMu     sync.Mutex
	RootPath     string
}

type UserServer interface {
	RegisterUser(req model.RegisterRequest) (*model.User, error)
	LoginUser(req model.LoginRequest) (*model.User, error)
}

var _ UserServer = (*model.UserBank)(nil)
