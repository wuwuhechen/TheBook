package service

import (
	"TheBook/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HandlerGetWrongQuestionPage 分页渲染当前用户的错题本。
func (s *Server) HandlerGetWrongQuestionPage(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated user"})
		return
	}
	page, ok := requestedPage(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
		return
	}
	const pageSize = 5
	wrongQuestions, err := s.WQ.ListWrongQuestionsByUser(userID, pageSize+1, (page-1)*pageSize)
	if err != nil {
		s.appLog().Error("读取错题本失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load wrong questions"})
		return
	}
	hasNext := len(wrongQuestions) > pageSize
	if hasNext {
		wrongQuestions = wrongQuestions[:pageSize]
	}
	items := make([]model.WrongQuestionItem, 0, len(wrongQuestions))
	for index, wrongQuestion := range wrongQuestions {
		question, err := s.DB.GetQuestion(int(wrongQuestion.QuestionID))
		if err != nil {
			continue
		}
		items = append(items, model.WrongQuestionItem{QuestionID: int(wrongQuestion.QuestionID), Number: index + 1, Category: question.Category, Question: question.Question, Choices: question.Choices, CorrectAnswerText: choiceText(question.Choices, question.Answer), Explanation: question.Explanation, WrongCount: wrongQuestion.WrongCount, LastWrongAt: wrongQuestion.LastWrongAt})
	}
	c.HTML(http.StatusOK, "wrong_question_page.html", model.WrongQuestionPageData{UserID: userID, Items: items, Page: page, HasPrev: page > 1, HasNext: hasNext, PrevPage: page - 1, NextPage: page + 1})
}

// HandlerPostWrongQuestionPractice 根据当前用户选中的错题创建新的练习。
func (s *Server) HandlerPostWrongQuestionPractice(c *gin.Context) {
	var req model.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.QuestionIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one wrong question ID is required"})
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	selectedIDs := make([]uint, 0, len(req.QuestionIDs))
	selectedSet := make(map[uint]struct{}, len(req.QuestionIDs))
	for _, questionID := range req.QuestionIDs {
		if questionID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Question IDs must be positive"})
			return
		}
		id := uint(questionID)
		if _, duplicate := selectedSet[id]; duplicate {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Question IDs must not contain duplicates"})
			return
		}
		selectedSet[id] = struct{}{}
		selectedIDs = append(selectedIDs, id)
	}
	wrongQuestions, err := s.WQ.ListWrongQuestionsByIDs(userID, selectedIDs)
	if err != nil {
		s.appLog().Error("读取选中错题失败", zap.Uint("user_id", userID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load wrong questions"})
		return
	}
	if len(wrongQuestions) != len(selectedIDs) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Selected questions must belong to your wrong-question list"})
		return
	}
	for _, questionID := range req.QuestionIDs {
		if _, err := s.DB.GetQuestion(questionID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}
	}
	practice := (&model.Practice{}).GenerateExamByQuestionIDs(req.QuestionIDs)
	practice.UserID = userID
	if err := s.PM.Create(practice); err != nil {
		s.appLog().Error("创建错题练习失败", zap.Uint("user_id", userID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wrong-question practice"})
		return
	}
	s.businessLog().Info("错题练习已创建", zap.Uint("user_id", userID), zap.Int("practice_id", practice.ID), zap.Int("question_count", practice.TotalQuestions))
	c.Redirect(http.StatusFound, "/practice/"+strconv.Itoa(practice.ID))
}
