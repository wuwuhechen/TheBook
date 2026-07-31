package service

import (
	"TheBook/config"
	"TheBook/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

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

func DataInit(path string) (*QuestionServer, error) {
	questionServer, err := LoadQuestions(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load questions: %v", err)
	}

	return questionServer, nil
}

func GinInit(path string, questionServer *QuestionServer) (*gin.Engine, error) {
	r := gin.Default()

	r.POST("/request", questionServer.HandlerPostSubmitAnswer)
	r.LoadHTMLFiles(path)

	return r, nil
}

func LoadQuestions(path string) (*QuestionServer, error) {
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

	return &QuestionServer{DB: bank}, nil
}

func (qs *QuestionServer) HandlerPostSubmitAnswer(c *gin.Context) {
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

	log.Printf("Question: %s\n", question.String())

	pageData := NewQuestionData(
		question,
		len(qs.DB.Questions),
	)

	// log.Printf("PageData: %s\n", pageData.String())

	c.HTML(200, "question_page.html", pageData)
}
