package service

import (
	"TheBook/auth"
	"TheBook/model"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// HandlerPostQuestion 校验请求的题目并重定向到题目页面。
func (s *Server) HandlerPostQuestion(c *gin.Context) {
	var userReq model.Request
	if err := c.ShouldBind(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := s.DB.GetQuestion(userReq.QuestionID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/question?question_id="+strconv.Itoa(userReq.QuestionID))
}

// HandlerPostRandomQuestion 重定向到随机选择的题目。
func (s *Server) HandlerPostRandomQuestion(c *gin.Context) {
	session := model.NewRandomSession(s.DB)
	if len(session.Questions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No questions available"})
		return
	}
	s.RS[session.ID] = session
	c.Redirect(http.StatusFound, "/question/random/"+strconv.Itoa(session.ID))
}

// HandlerPostCheckAnswer 检查单题答案并返回 JSON 结果。
func (s *Server) HandlerPostCheckAnswer(c *gin.Context) {
	var userReq model.Request
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	question, err := s.DB.GetQuestion(userReq.QuestionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.NewResponse(userReq.Choice == question.Answer, question.Explanation))
}

// HandlerPostPracticeInit 创建练习并重定向到第一题。
func (s *Server) HandlerPostPracticeInit(c *gin.Context) {
	var userReq model.Request
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	practice := (&model.Practice{}).NewPractice().GenerateExam(s.DB, userReq.PracticeSize)
	s.PM[practice.ID] = practice
	practice.Reset()
	c.Redirect(http.StatusFound, "/practice/"+strconv.Itoa(practice.ID))
}

// HandlerPostSubmitAnswer 保存练习中一道题的答案。
func (s *Server) HandlerPostSubmitAnswer(c *gin.Context) {
	var req model.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	practice := s.PM[req.PracticeID]
	if practice == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Practice not found"})
		return
	}
	practice.Answers[req.QuestionID] = req.Choice
	c.JSON(http.StatusOK, gin.H{"message": "saved"})
}

// HandlerSubmitPractice 批改练习、标记完成状态并返回汇总结果。
func (s *Server) HandlerSubmitPractice(c *gin.Context) {
	practiceID, err := strconv.Atoi(c.Param("practice_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid practice ID"})
		return
	}
	practice := s.PM[practiceID]
	if practice == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Practice not found"})
		return
	}
	results := practice.CheckPractice(s.DB)
	practice.Completed = true
	c.JSON(http.StatusOK, results)
}

func (s *Server) HandlerPostLogin(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.US.LoginUser(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.GenerateToken(user.UserID, user.Username, user.Role, user.Nickname, time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("access_token", token, 3600*24, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

// TODO：重定向到登录界面
func (s *Server) HandlerPostRegister(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.US.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if err := s.persistUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	token, err := auth.GenerateToken(user.UserID, user.Username, user.Role, user.Nickname, time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.Set("userID", user.UserID)
	c.SetCookie("access_token", token, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful", "token": token})
}

func (s *Server) persistUsers() error {
	bank, ok := s.US.(*model.UserBank)
	if !ok {
		return fmt.Errorf("unsupported user service type")
	}
	users := make([]*model.User, 0, len(bank.Users))
	for _, user := range bank.Users {
		users = append(users, user)
	}
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.UserPath, data, 0644)
}

func (s *Server) HandlerPostLogout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
