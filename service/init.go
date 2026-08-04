package service

import (
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

	frontEndPath := fmt.Sprintf("%s/%s", rootPath, cfg.FrontEnd.TemplatePath)
	r, err := GinInit(frontEndPath, questionServer)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gin: %v", err)
	}

	usPath := fmt.Sprintf("%s/%s", rootPath, cfg.User.UserBankPath)
	userServer, err := UserInit(usPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user data: %v", err)
	}

	return &Server{
		Router: r,
		QS:     questionServer,
		US:     userServer,
	}, nil
}

// DataInit 从 path 加载题目，并创建空的练习管理器。
func DataInit(path string) (*QuestionServer, error) {
	db, err := LoadQuestions(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load questions: %v", err)
	}

	return &QuestionServer{
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
func GinInit(path string, questionServer *QuestionServer) (*gin.Engine, error) {
	r := gin.Default()
	r.LoadHTMLGlob(path)

	// homeMode 管理首页入口。
	homeMode := r.Group("/")
	homeMode.GET("", questionServer.HandlerGetHomePage)

	// questionMode 管理单题浏览、随机出题和即时判题。
	questionMode := r.Group("/question")
	questionMode.POST("/random", questionServer.HandlerPostRandomQuestion)
	questionMode.POST("/request", questionServer.HandlerPostQuestion)
	questionMode.POST("/check_answer", questionServer.HandlerPostCheckAnswer)
	questionMode.GET("/random/:session_id", questionServer.HandlerGetRandomQuestionPage)
	questionMode.GET("", questionServer.HandlerGetQuestionPage)

	// practiceMode 管理套题初始化、答题、提交和结果查看。
	practiceMode := r.Group("/practice")
	practiceMode.POST("/init", questionServer.HandlerPostPracticeInit)
	practiceMode.POST("/answer", questionServer.HandlerPostSubmitAnswer)
	practiceMode.POST("/:practice_id/submit", questionServer.HandlerSubmitPractice)
	practiceMode.GET("/:practice_id", questionServer.HandlerGetPracticePage)
	practiceMode.GET("/:practice_id/result", questionServer.HandlerGetPracticeResultPage)

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
