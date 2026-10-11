package service

import (
	"TheBook/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HandlerPostRandomQuestion 创建固定顺序的随机答题会话。
func (s *Server) HandlerPostRandomQuestion(c *gin.Context) {
	session := model.NewRandomSession(s.DB)
	if len(session.Questions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No questions available"})
		return
	}
	_ = s.RS.Create(session)
	s.businessLog().Info("随机练习已创建", zap.Int("session_id", session.ID), zap.Int("question_count", len(session.Questions)))
	c.Redirect(http.StatusFound, "/question/random/"+strconv.Itoa(session.ID))
}

// HandlerGetRandomQuestionPage 渲染随机会话当前题并处理前后切换。
func (s *Server) HandlerGetRandomQuestionPage(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	sessionID, err := strconv.Atoi(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid random session ID"})
		return
	}
	session, err := s.RS.FindByID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve random session"})
		return
	}
	switch c.Query("direction") {
	case "last":
		if session.CurrentIndex > 0 {
			session.CurrentIndex--
		}
	case "next":
		if session.CurrentIndex < len(session.Questions)-1 {
			session.CurrentIndex++
		}
	}
	question, err := s.DB.GetQuestion(session.CurrentQuestionID())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	pageData := model.NewQuestionData(question, len(session.Questions))
	pageData.SetID(session.CurrentIndex + 1)
	pageData.SetUserID(userID)
	pageData.SetRandomSessionID(sessionID)
	pageData.SetHasLastID(session.CurrentIndex > 0)
	pageData.SetHasNextID(session.CurrentIndex < len(session.Questions)-1)
	c.HTML(http.StatusOK, "question_page.html", pageData)
}
