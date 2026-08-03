package service

import (
	"TheBook/config"
	"TheBook/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

/*
InitSystem 初始化系统，包括加载配置文件、初始化数据和Gin引擎

参数

	无

返回值

	*gin.Engine: Gin引擎实例
	*QuestionServer: 问题服务器实例
*/
func InitSystem() (*gin.Engine, *QuestionServer, error) {
	rootPath, err := utils.FindProjectRoot()

	cfg, err := config.LoadConfig(fmt.Sprintf("%s/config/config.yaml", rootPath))

	if err != nil {
		return nil, nil, fmt.Errorf("failed to find project root: %v", err)
	}

	qsPath := fmt.Sprintf("%s/%s", rootPath, cfg.Database.DatabasePath)
	questionServer, err := DataInit(qsPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize data: %v", err)
	}

	frontEndPath := fmt.Sprintf("%s/%s", rootPath, cfg.FrontEnd.TemplatePath)
	r, err := GinInit(frontEndPath, questionServer)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize Gin: %v", err)
	}

	return r, questionServer, nil
}

/*
DataInit 初始化数据，包括加载问题数据

参数

	path string: 问题数据文件的路径

返回值

	*QuestionServer: 问题服务器实例
	error: 错误信息
*/
func DataInit(path string) (*QuestionServer, error) {
	db, err := LoadQuestions(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load questions: %v", err)
	}

	pm := make(map[int]*Practice)

	var questionServer *QuestionServer
	questionServer = &QuestionServer{
		DB: db,
		PM: pm,
	}

	return questionServer, nil
}

/*
GinInit 初始化Gin引擎，包括设置路由和加载HTML模板

参数

	path string: HTML模板文件的路径
	questionServer *QuestionServer: 问题服务器实例

返回值

	*gin.Engine: Gin引擎实例
	error: 错误信息
*/
func GinInit(path string, questionServer *QuestionServer) (*gin.Engine, error) {
	r := gin.Default()

	r.LoadHTMLGlob(path)

	r.POST("/random", questionServer.HandlerPostRandomQuestion)

	r.POST("/request", questionServer.HandlerPostQuestion)

	r.POST("/check_answer", questionServer.HandlerPostCheckAnswer)

	r.POST("/practice_init", questionServer.HandlerPostPracticeInit)

	r.POST("/submit_answer", questionServer.HandlerPostSubmitAnswer)

	r.POST("/submit_practice/:practice_id", questionServer.HandlerSubmitPractice)

	r.GET("/question", questionServer.HandlerGetQuestionPage)

	r.GET("/practice/:practice_id", questionServer.HandlerGetPracticePage)

	r.GET("/", questionServer.HandlerGetHomePage)

	return r, nil
}

/*
LoadQuestions 加载问题数据

参数

	path string: 问题数据文件的路径

返回值

	*QuestionServer: 问题服务器实例
	error: 错误信息
*/
func LoadQuestions(path string) (QuestionManager, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open questions file: %v", err)
	}

	defer file.Close()

	var questionBank []Question
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&questionBank)
	if err != nil {
		return nil, fmt.Errorf("failed to decode questions: %v", err)
	}

	bank := &QuestionBank{
		Questions: make(map[int]Question),
	}

	for _, question := range questionBank {
		bank.AddQuestion(question)
	}

	return QuestionManager(bank), nil
}

/*
HandlerPostQuestion 处理用户提交答案的请求

参数

	c *gin.Context: Gin上下文对象
*/
func (qs *QuestionServer) HandlerPostQuestion(c *gin.Context) {
	var userReq Request
	if err := c.ShouldBind(&userReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := qs.DB.GetQuestion(userReq.QuestionID)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(
		http.StatusFound,
		"/question?question_id="+strconv.Itoa(userReq.QuestionID),
	)
}

/*
HandlerPostRandomQuestion 处理用户请求随机获取题目的请求

参数

	c *gin.Context: Gin上下文对象
*/
func (qs *QuestionServer) HandlerPostRandomQuestion(c *gin.Context) {
	question, err := qs.DB.GetRandomQuestionID()
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(
		http.StatusFound,
		"/question?question_id="+strconv.Itoa(question.ID),
	)
}

/*
HandlerPostCheckAnswer 处理用户提交答案的请求

参数

	c *gin.Context: Gin上下文对象
*/
func (qs *QuestionServer) HandlerPostCheckAnswer(c *gin.Context) {
	var userReq Request

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	question, err := qs.DB.GetQuestion(userReq.QuestionID)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	correct := userReq.Choice == question.Answer

	response := NewResponse(correct, question.Explanation)
	c.JSON(200, response)
}

/*
HandlerGetQuestionPage 处理获取问题页面的请求

参数

	c *gin.Context: Gin上下文对象
*/
func (qs *QuestionServer) HandlerGetQuestionPage(c *gin.Context) {
	questionID, err := strconv.Atoi(c.Query("question_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid question ID"})
		return
	}

	question, err := qs.DB.GetQuestion(questionID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Question not found"})
		return
	}

	pageData := NewQuestionData(
		question,
		qs.DB.GetTotalCount(),
	)

	c.HTML(http.StatusOK, "question_page.html", pageData)

}

/*
HandlerGetHomePage 处理获取首页的请求

参数

	c *gin.Context: Gin上下文对象
*/
func (qs *QuestionServer) HandlerGetHomePage(c *gin.Context) {
	c.HTML(200, "home_page.html", nil)
}

func (qs *QuestionServer) HandlerPostPracticeInit(c *gin.Context) {
	var userReq Request
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	practice := (&Practice{}).NewPractice()
	practice = practice.GenerateExam(qs.DB, userReq.PracticeSize)

	qs.PM[practice.ID] = practice
	practice.Reset()

	c.Redirect(
		http.StatusFound,
		"/practice/"+strconv.Itoa(practice.ID),
	)
}

func (qs *QuestionServer) HandlerPostSubmitAnswer(c *gin.Context) {
	var req Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	practice := qs.PM[req.PracticeID]
	if practice == nil {
		c.JSON(404, gin.H{"error": "Practice not found"})
		return
	}

	practice.Answers[req.QuestionID] = req.Choice

	c.JSON(200, gin.H{
		"message": "saved",
	})
}

func (qs *QuestionServer) HandlerSubmitPractice(c *gin.Context) {
	requestIDStr := c.Param("practice_id")
	practiceID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid practice ID"})
		return
	}

	practice := qs.PM[practiceID]
	if practice == nil {
		c.JSON(404, gin.H{"error": "Practice not found"})
		return
	}

	results := practice.CheckPractice(qs)

	practice.Completed = true

	c.JSON(200, results)
}

func (qs *QuestionServer) HandlerGetPracticePage(c *gin.Context) {
	requestIDStr := c.Param("practice_id")
	practiceID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid practice ID"})
		return
	}

	practice := qs.PM[practiceID]
	if practice == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Practice not found"})
		return
	}
	if practice.Completed {
		c.JSON(http.StatusGone, gin.H{"error": "Practice already submitted"})
		return
	}

	switch c.Query("direction") {
	case "last":
		if practice.CurrentIndex >= 0 {
			practice.CurrentIndex--
		}
	case "next":
		if practice.CurrentIndex < len(practice.Questions)-1 {
			practice.CurrentIndex++
		}
	}

	questionID := practice.GetCurrentQuestionID()
	lastID := 0
	if practice.CurrentIndex > 0 {
		lastID = practice.CurrentIndex - 1
	}
	nextID := 0
	if practice.CurrentIndex < len(practice.Questions)-1 {
		nextID = practice.CurrentIndex + 1
	}

	question, err := qs.DB.GetQuestion(questionID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Question not found"})
		return
	}

	remainingTime := practice.GetDuration() - time.Since(practice.StartTime)

	pageData := NewQuestionData(question, practice.TotalQuestions)
	pageData.SetID(practice.CurrentIndex + 1)
	pageData.SetPracticeID(practiceID)
	pageData.SetLastID(lastID)
	pageData.SetNextID(nextID, practice.TotalQuestions)
	pageData.SetHasLastID(practice.CurrentIndex > 0)
	pageData.SetHasNextID(practice.CurrentIndex < len(practice.Questions)-1)
	pageData.SetDuration(remainingTime)

	c.HTML(200, "practice_page.html", pageData)
}
