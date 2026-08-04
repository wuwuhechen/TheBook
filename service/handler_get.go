package service

import (
	"TheBook/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// HandlerGetQuestionPage 渲染独立答题页面。
func (qs *QuestionServer) HandlerGetQuestionPage(c *gin.Context) {
	questionID, err := strconv.Atoi(c.Query("question_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}
	question, err := qs.DB.GetQuestion(questionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	c.HTML(http.StatusOK, "question_page.html", model.NewQuestionData(question, qs.DB.GetTotalCount()))
}

// HandlerGetRandomQuestionPage 渲染随机答题会话中的当前题目并处理前后切换。
func (qs *QuestionServer) HandlerGetRandomQuestionPage(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid random session ID"})
		return
	}

	session := qs.RS[sessionID]
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Random session not found"})
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

	question, err := qs.DB.GetQuestion(session.CurrentQuestionID())
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
func (qs *QuestionServer) HandlerGetHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "home_page.html", nil)
}

// HandlerGetPracticePage 渲染当前练习题目并处理题目导航。
func (qs *QuestionServer) HandlerGetPracticePage(c *gin.Context) {
	practiceID, err := strconv.Atoi(c.Param("practice_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid practice ID"})
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
	question, err := qs.DB.GetQuestion(questionID)
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
func (qs *QuestionServer) HandlerGetPracticeResultPage(c *gin.Context) {
	practiceID, err := strconv.Atoi(c.Param("practice_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid practice ID"})
		return
	}
	practice := qs.PM[practiceID]
	if practice == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Practice not found"})
		return
	}
	if !practice.Completed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Practice not completed"})
		return
	}

	items := make([]model.PracticeResultItem, 0, len(practice.Questions))
	correctCount, wrongCount := 0, 0
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

func choiceText(choices []string, index int) string {
	if index < 0 || index >= len(choices) {
		return ""
	}
	return string(rune('A'+index)) + ". " + choices[index]
}
