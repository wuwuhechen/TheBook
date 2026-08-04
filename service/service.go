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

// InitSystem 加载配置并初始化题目服务与 Gin 引擎。
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

// DataInit 从 path 加载题目，并创建空的练习管理器。
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

// GinInit 创建 Gin 引擎、加载 path 中的模板并注册路由。
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

	r.GET("/practice/:practice_id/result", questionServer.HandlerGetPracticeResultPage)

	r.GET("/practice/:practice_id", questionServer.HandlerGetPracticePage)

	r.GET("/", questionServer.HandlerGetHomePage)

	return r, nil
}

// LoadQuestions 解析 path 指向的 JSON 题库文件。
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

// HandlerPostQuestion 校验请求的题目并重定向到题目页面。
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

// HandlerPostRandomQuestion 重定向到随机选择的题目。
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

// HandlerPostCheckAnswer 检查单题答案并返回 JSON 结果。
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

// HandlerGetQuestionPage 渲染独立答题页面。
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

// HandlerGetHomePage 渲染应用首页。
func (qs *QuestionServer) HandlerGetHomePage(c *gin.Context) {
	c.HTML(200, "home_page.html", nil)
}

// HandlerPostPracticeInit 创建练习并重定向到第一题。
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

// HandlerPostSubmitAnswer 保存练习中一道题的答案。
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

// HandlerSubmitPractice 批改练习、标记完成状态并返回汇总结果。
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

// HandlerGetPracticePage 渲染当前练习题目并处理题目导航。
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

// HandlerGetPracticeResultPage 渲染已完成练习的答案与解析页面。
func (qs *QuestionServer) HandlerGetPracticeResultPage(c *gin.Context) {
	practiceIDStr := c.Param("practice_id")
	practiceID, err := strconv.Atoi(practiceIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid practice ID"})
		return
	}

	practice := qs.PM[practiceID]
	if practice == nil {
		c.JSON(404, gin.H{"error": "Practice not found"})
		return
	}

	if !practice.Completed {
		c.JSON(400, gin.H{"error": "Practice not completed"})
		return
	}

	items := make([]PracticeResultItem, 0, len(practice.Questions))

	correctCount := 0
	wrongCount := 0

	for idx, questionID := range practice.Questions {
		question, err := qs.DB.GetQuestion(questionID)
		if err != nil {
			continue
		}

		userAnswer, answered := practice.Answers[questionID]
		correct := answered && userAnswer == question.Answer

		if correct {
			correctCount++
		} else {
			wrongCount++
		}

		item := PracticeResultItem{
			Number:            idx + 1,
			RealID:            question.ID,
			Category:          question.Category,
			Question:          question.Question,
			Choices:           question.Choices,
			UserAnswer:        userAnswer,
			Answered:          answered,
			UserAnswerText:    choiceText(question.Choices, userAnswer),
			CorrectAnswer:     question.Answer,
			CorrectAnswerText: choiceText(question.Choices, question.Answer),
			Correct:           correct,
			Explanation:       question.Explanation,
		}

		items = append(items, item)
	}

	c.HTML(http.StatusOK, "practice_result_page.html", PracticeResultPageData{
		PracticeID:   practiceID,
		Total:        len(items),
		CorrectCount: correctCount,
		WrongCount:   wrongCount,
		Items:        items,
	})
}

func choiceText(choices []string, index int) string {
	if index < 0 || index >= len(choices) {
		return ""
	}

	return string(rune('A'+index)) + ". " + choices[index]
}
