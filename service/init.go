package service

import (
	"TheBook/auth"
	"TheBook/config"
	"TheBook/logger"
	"TheBook/model"
	"TheBook/utils"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// InitSystem 加载配置并初始化题目服务与 Gin 引擎。
func InitSystem(log *logger.Logger) (*Server, error) {
	rootPath, err := utils.FindProjectRoot()

	cfg, err := config.LoadConfig(fmt.Sprintf("%s/config/config.yaml", rootPath))
	if err != nil {
		return nil, fmt.Errorf("failed to find project root: %v", err)
	}

	dbPath := fmt.Sprintf("%s/%s", rootPath, cfg.Database.DatabasePath)
	questionServer, err := DataInit(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize data: %v", err)
	}

	usPath := fmt.Sprintf("%s/%s", rootPath, cfg.User.UserBankPath)
	userHashPath := fmt.Sprintf("%s/%s", rootPath, cfg.User.UserHashPath)
	userServer, err := UserInit(usPath, userHashPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user data: %v", err)
	}

	qsPath := fmt.Sprintf("%s/%s", rootPath, cfg.User.QuestionProgressPath)
	questionProgresses, err := QuestionProgressInit(qsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize question progresses: %v", err)
	}

	server := &Server{
		DB:           questionServer.DB,
		PM:           questionServer.PM,
		RS:           questionServer.RS,
		US:           userServer,
		QS:           questionProgresses,
		UserPath:     usPath,
		UserHashPath: userHashPath,
		RecordPath:   fmt.Sprintf("%s/database/practice_records.json", rootPath),
		RootPath:     rootPath,
	}

	frontEndPath := fmt.Sprintf("%s/%s", rootPath, cfg.FrontEnd.TemplatePath)
	r, err := GinInit(frontEndPath, server, log)
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

func UserInit(path, hashPath string) (*model.UserBank, error) {
	userBank, err := LoadUsers(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load user bank: %v", err)
	}

	if err := LoadUserPasswordHashes(hashPath, userBank); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load user password hashes: %v", err)
		}
		for username, user := range userBank.Users {
			passwordHash, err := auth.HashPassword(user.Password)
			if err != nil {
				return nil, fmt.Errorf("failed to hash user password: %v", err)
			}
			userBank.SetPasswordHash(username, passwordHash)
		}
		if err := persistUserPasswordHashes(hashPath, userBank); err != nil {
			return nil, fmt.Errorf("failed to create user password hash file: %v", err)
		}
	}

	return userBank, nil
}

func QuestionProgressInit(path string) (map[uint]*model.QuestionProgress, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open question progress file: %v", err)
	}
	defer file.Close()

	var progresses []*model.QuestionProgress
	if err := json.NewDecoder(file).Decode(&progresses); err != nil {
		return nil, fmt.Errorf("failed to decode question progresses: %v", err)
	}

	progressMap := make(map[uint]*model.QuestionProgress)
	for _, progress := range progresses {
		progressMap[progress.UserID] = progress
	}

	return progressMap, nil
}

// GinInit 创建 Gin 引擎、加载 path 中的模板并注册路由。
func GinInit(path string, Server *Server, log *logger.Logger) (*gin.Engine, error) {
	r := gin.New()

	r.Use(
		logger.RequestLogger(log),
		gin.Recovery(),
	)

	r.LoadHTMLGlob(path)

	// homeMode 管理首页入口。
	homeMode := r.Group("/")
	homeMode.GET("", Server.HandlerGetHomePage)

	// authMode 管理用户注册、登录和登出。
	authMode := r.Group("/auth")
	authMode.GET("/register", Server.HandlerGetRegisterPage)
	authMode.GET("/login", Server.HandlerGetLoginPage)
	authMode.POST("/register", Server.HandlerPostRegister)
	authMode.POST("/login", Server.HandlerPostLogin)
	authMode.POST("/logout", Server.HandlerPostLogout)

	// questionMode 管理单题浏览、随机出题和即时判题。
	questionMode := r.Group("/question", auth.AuthMiddleware())
	questionMode.POST("/current", Server.HandlerPostCurrentQuestion)
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

// LoadUserPasswordHashes 将哈希账户文件载入到 UserBank 的密码哈希索引中。
func LoadUserPasswordHashes(path string, userBank *model.UserBank) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var users []*model.User
	if err := json.NewDecoder(file).Decode(&users); err != nil {
		return fmt.Errorf("failed to decode user password hashes: %v", err)
	}
	for _, user := range users {
		if user == nil {
			continue
		}
		if _, exists := userBank.GetUser(user.Username); exists {
			userBank.SetPasswordHash(user.Username, user.Password)
		}
	}
	return nil
}
