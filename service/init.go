package service

import (
	"TheBook/auth"
	"TheBook/config"
	"TheBook/model"
	"TheBook/utils"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// InitSystem 加载配置并初始化题目服务与 Gin 引擎。
func InitSystem() (*Server, error) {
	rootPath, err := utils.FindProjectRoot()

	cfg, err := config.LoadConfig(fmt.Sprintf("%s/config/config.yaml", rootPath))
	if err != nil {
		return nil, fmt.Errorf("failed to find project root: %v", err)
	}

	qsPath := fmt.Sprintf("%s/%s", rootPath, cfg.Database.DatabasePath)
	questionServer, err := DataInit(qsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize data: %v", err)
	}

	usPath := fmt.Sprintf("%s/%s", rootPath, cfg.User.UserBankPath)
	userServer, err := UserInit(usPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user data: %v", err)
	}

	server := &Server{
		DB: questionServer.DB,
		PM: questionServer.PM,
		RS: questionServer.RS,
		US: userServer,
	}

	frontEndPath := fmt.Sprintf("%s/%s", rootPath, cfg.FrontEnd.TemplatePath)
	r, err := GinInit(frontEndPath, server)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gin: %v", err)
	}

	auth.InitJWT(cfg.Auth.JWTSecret)

	server.Router = r

	return server, nil
}

// DataInit 从 path 加载题目，并创建空的练习管理器。
func DataInit(path string) (*Server, error) {
	db, err := LoadQuestions(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load questions: %v", err)
	}

	return &Server{
		DB: db,
		PM: make(map[int]*model.Practice),
		RS: make(map[int]*model.RandomSession),
	}, nil
}

func UserInit(path string) (*model.UserBank, error) {
	userBank, err := LoadUsers(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load user bank: %v", err)
	}

	return userBank, nil
}

// GinInit 创建 Gin 引擎、加载 path 中的模板并注册路由。
func GinInit(path string, Server *Server) (*gin.Engine, error) {
	r := gin.Default()
	r.LoadHTMLGlob(path)

	// homeMode 管理首页入口。
	homeMode := r.Group("/")
	homeMode.GET("", Server.HandlerGetHomePage)

	// authMode 管理用户注册、登录和登出。
	authMode := r.Group("/auth")
	authMode.POST("/register", Server.HandlerPostRegister)
	authMode.POST("/login", Server.HandlerPostLogin)
	authMode.POST("/logout", Server.HandlerPostLogout)

	// questionMode 管理单题浏览、随机出题和即时判题。
	questionMode := r.Group("/question", auth.AuthMiddleware())
	questionMode.POST("/random", Server.HandlerPostRandomQuestion)
	questionMode.POST("/request", Server.HandlerPostQuestion)
	questionMode.POST("/check_answer", Server.HandlerPostCheckAnswer)
	questionMode.GET("/random/:session_id", Server.HandlerGetRandomQuestionPage)
	questionMode.GET("", Server.HandlerGetQuestionPage)

	// practiceMode 管理套题初始化、答题、提交和结果查看。
	practiceMode := r.Group("/practice", auth.AuthMiddleware())
	practiceMode.POST("/init", Server.HandlerPostPracticeInit)
	practiceMode.POST("/answer", Server.HandlerPostSubmitAnswer)
	practiceMode.POST("/:practice_id/submit", Server.HandlerSubmitPractice)
	practiceMode.GET("/:practice_id", Server.HandlerGetPracticePage)
	practiceMode.GET("/:practice_id/result", Server.HandlerGetPracticeResultPage)

	return r, nil
}

// LoadQuestions 解析 path 指向的 JSON 题库文件。
func LoadQuestions(path string) (model.QuestionManager, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open questions file: %v", err)
	}
	defer file.Close()

	var questionBank []model.Question
	if err := json.NewDecoder(file).Decode(&questionBank); err != nil {
		return nil, fmt.Errorf("failed to decode questions: %v", err)
	}

	bank := &model.QuestionBank{Questions: make(map[int]model.Question)}
	for _, question := range questionBank {
		bank.AddQuestion(question)
	}

	return bank, nil
}

func LoadUsers(path string) (*model.UserBank, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open user bank file: %v", err)
	}
	defer file.Close()

	var users []*model.User
	if err := json.NewDecoder(file).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %v", err)
	}

	userBank := model.NewUserBank()
	for _, user := range users {
		userBank.AddUser(user)
	}

	return userBank, nil
}
