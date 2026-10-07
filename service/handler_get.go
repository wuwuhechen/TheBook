package service

import (
	"TheBook/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// HandlerGetRegisterPage 渲染用户注册页面。
func (s *Server) HandlerGetRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register_page.html", nil)
}

// HandlerGetLoginPage 渲染用户登录页面。
func (s *Server) HandlerGetLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login_page.html", nil)
}

// HandlerGetQuestionPage 渲染独立答题页面。
func (s *Server) HandlerGetQuestionPage(c *gin.Context) {
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
	c.HTML(http.StatusOK, "question_page.html", model.NewQuestionData(question, s.DB.GetTotalCount()))
}

// HandlerGetRandomQuestionPage 渲染随机答题会话中的当前题目并处理前后切换。
func (s *Server) HandlerGetRandomQuestionPage(c *gin.Context) {
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
	pageData.SetRandomSessionID(sessionID)
	pageData.SetHasLastID(session.CurrentIndex > 0)
	pageData.SetHasNextID(session.CurrentIndex < len(session.Questions)-1)
	c.HTML(http.StatusOK, "question_page.html", pageData)
}

// HandlerGetHomePage 渲染应用首页。
func (s *Server) HandlerGetHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "home_page.html", nil)
}

// HandlerGetPracticePage 渲染当前练习题目并处理题目导航。
func (s *Server) HandlerGetPracticePage(c *gin.Context) {
	practiceID, err := strconv.Atoi(c.Param("practice_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid practice ID"})
		return
	}
	practice, err := s.getOwnedPracitce(c, practiceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
	question, err := s.DB.GetQuestion(questionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	pageData := model.NewQuestionData(question, practice.TotalQuestions)
	pageData.SetID(practice.CurrentIndex + 1)
	pageData.SetPracticeID(practiceID)
	pageData.SetHasLastID(practice.CurrentIndex > 0)
	pageData.SetHasNextID(practice.CurrentIndex < len(practice.Questions)-1)
	pageData.SetDuration(practice.GetDuration() - time.Since(practice.StartTime))
	c.HTML(http.StatusOK, "practice_page.html", pageData)
}

// HandlerGetPracticeResultPage 渲染已完成练习的答案与解析页面。
func (s *Server) HandlerGetPracticeResultPage(c *gin.Context) {
	practiceID, err := strconv.Atoi(c.Param("practice_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid practice ID"})
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated user"})
		return
	}

	if record, recordErr := s.PM.GetRecordByPracticeID(userID, practiceID); recordErr == nil {
		c.HTML(http.StatusOK, "practice_result_page.html", s.practiceResultFromRecord(record))
		return
	}

	practice, err := s.getOwnedPracitce(c, practiceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Practice result not found"})
		return
	}
	if !practice.Completed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Practice not completed"})
		return
	}

	items := make([]model.PracticeResultItem, 0, len(practice.Questions))
	correctCount, wrongCount := 0, 0
	for idx, questionID := range practice.Questions {
		question, err := s.DB.GetQuestion(questionID)
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
		items = append(items, model.PracticeResultItem{
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
		})
	}

	c.HTML(http.StatusOK, "practice_result_page.html", model.PracticeResultPageData{
		PracticeID: practiceID, Total: len(items), CorrectCount: correctCount,
		WrongCount: wrongCount, Items: items,
	})
}

// HandlerGetPracticeHistoryPage 渲染当前登录用户的历史套题记录。
func (s *Server) HandlerGetPracticeHistoryPage(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated user"})
		return
	}

	page := 1
	if value := c.Query("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
			return
		}
		page = parsed
	}

	const pageSize = 5
	records, err := s.PM.ListRecordsByUser(userID, pageSize+1, (page-1)*pageSize)
	if err != nil {
		s.appLog().Error("读取练习历史失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load practice history"})
		return
	}

	hasNext := len(records) > pageSize
	if hasNext {
		records = records[:pageSize]
	}
	items := make([]model.PracticeHistoryItem, 0, len(records))
	for _, record := range records {
		rate := 0
		if record.TotalQuestions > 0 {
			rate = record.CorrectCount * 100 / record.TotalQuestions
		}
		items = append(items, model.PracticeHistoryItem{
			PracticeID:     record.PracticeID,
			TotalQuestions: record.TotalQuestions,
			CorrectCount:   record.CorrectCount,
			WrongCount:     record.WrongCount,
			CorrectRate:    rate,
			SubmitTime:     record.SubmitTime,
		})
	}
	c.HTML(http.StatusOK, "practice_history_page.html", model.PracticeHistoryPageData{
		Items:    items,
		Page:     page,
		HasPrev:  page > 1,
		HasNext:  hasNext,
		PrevPage: page - 1,
		NextPage: page + 1,
	})
}

func (s *Server) practiceResultFromRecord(record *model.PracticeRecord) model.PracticeResultPageData {
	items := make([]model.PracticeResultItem, 0, len(record.Answers))
	for index, answer := range record.Answers {
		question, err := s.DB.GetQuestion(int(answer.QuestionID))
		if err != nil {
			continue
		}
		items = append(items, model.PracticeResultItem{
			Number:            index + 1,
			RealID:            question.ID,
			Category:          question.Category,
			Question:          question.Question,
			Choices:           question.Choices,
			UserAnswer:        answer.Answer,
			Answered:          answer.Answered,
			UserAnswerText:    choiceText(question.Choices, answer.Answer),
			CorrectAnswer:     question.Answer,
			CorrectAnswerText: choiceText(question.Choices, question.Answer),
			Correct:           answer.Correct,
			Explanation:       question.Explanation,
		})
	}
	return model.PracticeResultPageData{
		PracticeID:   record.PracticeID,
		Total:        record.TotalQuestions,
		CorrectCount: record.CorrectCount,
		WrongCount:   record.WrongCount,
		Items:        items,
	}
}

func currentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

func choiceText(choices []string, index int) string {
	if index < 0 || index >= len(choices) {
		return ""
	}
	return string(rune('A'+index)) + ". " + choices[index]
}
