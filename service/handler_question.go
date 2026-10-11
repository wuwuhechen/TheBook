package service

import (
	"TheBook/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HandlerGetQuestionPage 渲染独立答题页面。
func (s *Server) HandlerGetQuestionPage(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	questionID, err := strconv.Atoi(c.Query("question_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}
	question, err := s.DB.GetQuestion(questionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	pageData := model.NewQuestionData(question, s.DB.GetTotalCount())
	pageData.SetUserID(userID)
	c.HTML(http.StatusOK, "question_page.html", pageData)
}

// HandlerPostQuestion 更新顺序答题进度并跳转到指定题目。
func (s *Server) HandlerPostQuestion(c *gin.Context) {
	var req model.Request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := s.DB.GetQuestion(req.QuestionID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	progress, err := s.QS.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve question progress"})
		return
	}
	progress.CurrentQuestionID, progress.UpdatedAt = req.QuestionID, time.Now()
	if err := s.QS.Upsert(progress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question progress"})
		return
	}
	if err := s.QS.Persist(); err != nil {
		s.appLog().Error("保存顺序答题进度失败", zap.Uint("user_id", userID), zap.Int("question_id", req.QuestionID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save question progress"})
		return
	}
	s.businessLog().Info("顺序答题进度已更新", zap.Uint("user_id", userID), zap.Int("question_id", req.QuestionID))
	c.Redirect(http.StatusFound, "/question?question_id="+strconv.Itoa(req.QuestionID))
}

// HandlerPostCurrentQuestion 恢复当前用户上次答题位置；失效位置会回退到第一题。
func (s *Server) HandlerPostCurrentQuestion(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	progress, err := s.QS.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve question progress"})
		return
	}
	questionID := progress.CurrentQuestionID
	if questionID == 0 {
		var exists bool
		questionID, exists = s.firstQuestionID()
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "No questions available"})
			return
		}
	} else if _, err := s.DB.GetQuestion(questionID); err != nil {
		var exists bool
		questionID, exists = s.firstQuestionID()
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "No questions available"})
			return
		}
	}
	progress.CurrentQuestionID, progress.UpdatedAt = questionID, time.Now()
	if err := s.QS.Upsert(progress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question progress"})
		return
	}
	if err := s.QS.Persist(); err != nil {
		s.appLog().Error("恢复顺序答题进度后保存失败", zap.Uint("user_id", userID), zap.Int("question_id", questionID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save question progress"})
		return
	}
	s.businessLog().Info("顺序答题进度已恢复", zap.Uint("user_id", userID), zap.Int("question_id", questionID))
	c.Redirect(http.StatusFound, "/question?question_id="+strconv.Itoa(questionID))
}

// firstQuestionID 返回题库中编号最小的有效题目 ID。
func (s *Server) firstQuestionID() (int, bool) {
	questions := s.DB.GetAllQuestionsSorted()
	if len(questions) == 0 {
		return 0, false
	}
	return questions[0].ID, true
}

// HandlerPostCheckAnswer 批改单题答案并记录错误题。
func (s *Server) HandlerPostCheckAnswer(c *gin.Context) {
	var req model.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	question, err := s.DB.GetQuestion(req.QuestionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	correct := req.Choice == question.Answer
	if !correct {
		_ = s.WQ.RecordWrongQuestion(uint(req.UserID), uint(req.QuestionID), time.Now())
	}
	s.businessLog().Info("单题判题完成", zap.Int("question_id", req.QuestionID), zap.Bool("correct", correct))
	c.JSON(http.StatusOK, model.NewResponse(correct, question.Explanation))
}
