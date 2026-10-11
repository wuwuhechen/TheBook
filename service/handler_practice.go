package service

import (
	"TheBook/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HandlerPostPracticeInit 创建一套随机练习并跳转到第一题。
func (s *Server) HandlerPostPracticeInit(c *gin.Context) {
	var req model.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	practice := (&model.Practice{}).NewPractice().GenerateExam(s.DB, req.PracticeSize)
	practice.UserID = userID
	if err := s.PM.Create(practice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	practice.Reset()
	s.businessLog().Info("套题已创建", zap.Uint("user_id", practice.UserID), zap.Int("practice_id", practice.ID), zap.Int("question_count", practice.TotalQuestions))
	c.Redirect(http.StatusFound, "/practice/"+strconv.Itoa(practice.ID))
}

// HandlerGetPracticePage 渲染当前练习题目并处理题目导航。
func (s *Server) HandlerGetPracticePage(c *gin.Context) {
	practiceID, err := strconv.Atoi(c.Param("practice_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid practice ID"})
		return
	}
	practice, err := s.getOwnedPractice(c, practiceID)
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
		if practice.CurrentIndex > 0 {
			practice.CurrentIndex--
		}
	case "next":
		if practice.CurrentIndex < len(practice.Questions)-1 {
			practice.CurrentIndex++
		}
	}
	question, err := s.DB.GetQuestion(practice.GetCurrentQuestionID())
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

// HandlerPostSubmitAnswer 保存当前套题中一题的答案，重复提交覆盖旧值。
func (s *Server) HandlerPostSubmitAnswer(c *gin.Context) {
	var req model.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	practice, err := s.getOwnedPractice(c, req.PracticeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	practice.Answers[req.QuestionID] = req.Choice
	s.businessLog().Info("套题答案已保存", zap.Int("practice_id", req.PracticeID), zap.Int("question_id", req.QuestionID))
	c.JSON(http.StatusOK, gin.H{"message": "saved"})
}

// HandlerSubmitPractice 批改练习，持久化结果并记录错题。
func (s *Server) HandlerSubmitPractice(c *gin.Context) {
	practiceID, err := strconv.Atoi(c.Param("practice_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid practice ID"})
		return
	}
	practice, err := s.getOwnedPractice(c, practiceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if practice.Completed {
		c.JSON(http.StatusConflict, gin.H{"error": "Practice already submitted"})
		return
	}
	results := practice.CheckPractice(s.DB)
	record, err := s.convertToPracticeRecord(practice, results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to build practice record"})
		return
	}
	if err := s.PM.Persist(record); err != nil {
		s.appLog().Error("保存套题记录失败", zap.Uint("user_id", practice.UserID), zap.Int("practice_id", practice.ID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save practice record"})
		return
	}
	if err := s.recordWrongAnswers(record); err != nil {
		s.appLog().Error("保存错题本失败", zap.Uint("user_id", practice.UserID), zap.Int("practice_id", practice.ID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save wrong questions"})
		return
	}
	practice.Completed = true
	s.businessLog().Info("套题已提交", zap.Uint("user_id", practice.UserID), zap.Int("practice_id", practice.ID), zap.Int("correct_count", results.CorrectCount), zap.Int("wrong_count", results.WrongCount))
	c.JSON(http.StatusOK, results)
}

// HandlerGetPracticeResultPage 优先从已持久化记录渲染练习结果。
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
	if record, err := s.PM.GetRecordByPracticeID(userID, practiceID); err == nil {
		c.HTML(http.StatusOK, "practice_result_page.html", s.practiceResultFromRecord(record))
		return
	}
	practice, err := s.getOwnedPractice(c, practiceID)
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
	for index, questionID := range practice.Questions {
		question, err := s.DB.GetQuestion(questionID)
		if err != nil {
			continue
		}
		answer, answered := practice.Answers[questionID]
		correct := answered && answer == question.Answer
		if correct {
			correctCount++
		} else {
			wrongCount++
		}
		items = append(items, practiceResultItem(index+1, question, answer, answered, correct))
	}
	c.HTML(http.StatusOK, "practice_result_page.html", model.PracticeResultPageData{PracticeID: practiceID, Total: len(items), CorrectCount: correctCount, WrongCount: wrongCount, Items: items})
}

// HandlerGetPracticeHistoryPage 分页渲染当前用户的历史套题记录。
func (s *Server) HandlerGetPracticeHistoryPage(c *gin.Context) {
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
		items = append(items, model.PracticeHistoryItem{PracticeID: record.PracticeID, TotalQuestions: record.TotalQuestions, CorrectCount: record.CorrectCount, WrongCount: record.WrongCount, CorrectRate: rate, SubmitTime: record.SubmitTime})
	}
	c.HTML(http.StatusOK, "practice_history_page.html", model.PracticeHistoryPageData{Items: items, Page: page, HasPrev: page > 1, HasNext: hasNext, PrevPage: page - 1, NextPage: page + 1})
}

// getOwnedPractice 查询并校验当前用户拥有的套题。
func (s *Server) getOwnedPractice(c *gin.Context, practiceID int) (*model.Practice, error) {
	practice, err := s.PM.FindByID(practiceID)
	if err != nil {
		return nil, err
	}
	if practice == nil {
		return nil, fmt.Errorf("practice not found")
	}
	userID, ok := currentUserID(c)
	if !ok || practice.UserID != userID {
		return nil, fmt.Errorf("no permission to access this practice")
	}
	return practice, nil
}

// convertToPracticeRecord 将完成的套题及批改结果转换为持久化记录。
func (s *Server) convertToPracticeRecord(practice *model.Practice, results *model.PracticeResponse) (*model.PracticeRecord, error) {
	answers := make([]model.AnswerRecord, 0, len(practice.Questions))
	for _, id := range practice.Questions {
		question, err := s.DB.GetQuestion(id)
		if err != nil {
			return nil, err
		}
		answer, answered := practice.Answers[id]
		answers = append(answers, model.AnswerRecord{QuestionID: uint(id), Answer: answer, Answered: answered, Correct: answered && answer == question.Answer})
	}
	return &model.PracticeRecord{UserID: practice.UserID, PracticeID: practice.ID, TotalQuestions: results.Total, CorrectCount: results.CorrectCount, WrongCount: results.WrongCount, StartTime: practice.StartTime, SubmitTime: time.Now(), Answers: answers}, nil
}

// recordWrongAnswers 将练习记录中的错误或未作答题目写入错题本。
func (s *Server) recordWrongAnswers(record *model.PracticeRecord) error {
	for _, answer := range record.Answers {
		if !answer.Correct {
			if s.WQ == nil {
				return fmt.Errorf("wrong question manager is not configured")
			}
			if err := s.WQ.RecordWrongQuestion(record.UserID, answer.QuestionID, record.SubmitTime); err != nil {
				return err
			}
		}
	}
	return nil
}

// practiceResultFromRecord 将已持久化练习记录转换为结果页面数据。
func (s *Server) practiceResultFromRecord(record *model.PracticeRecord) model.PracticeResultPageData {
	items := make([]model.PracticeResultItem, 0, len(record.Answers))
	for index, answer := range record.Answers {
		question, err := s.DB.GetQuestion(int(answer.QuestionID))
		if err == nil {
			items = append(items, practiceResultItem(index+1, question, answer.Answer, answer.Answered, answer.Correct))
		}
	}
	return model.PracticeResultPageData{PracticeID: record.PracticeID, Total: record.TotalQuestions, CorrectCount: record.CorrectCount, WrongCount: record.WrongCount, Items: items}
}

// practiceResultItem 组合一道题及其作答状态为结果页面条目。
func practiceResultItem(number int, question *model.Question, answer int, answered, correct bool) model.PracticeResultItem {
	return model.PracticeResultItem{Number: number, RealID: question.ID, Category: question.Category, Question: question.Question, Choices: question.Choices, UserAnswer: answer, Answered: answered, UserAnswerText: choiceText(question.Choices, answer), CorrectAnswer: question.Answer, CorrectAnswerText: choiceText(question.Choices, question.Answer), Correct: correct, Explanation: question.Explanation}
}

// requestedPage 解析查询参数中的正整数页码。
func requestedPage(c *gin.Context) (int, bool) {
	value := c.Query("page")
	if value == "" {
		return 1, true
	}
	page, err := strconv.Atoi(value)
	return page, err == nil && page >= 1
}
