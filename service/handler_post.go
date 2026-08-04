package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HandlerPostQuestion 校验请求的题目并重定向到题目页面。
func (qs *QuestionServer) HandlerPostQuestion(c *gin.Context) {
	var userReq Request
	if err := c.ShouldBind(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := qs.DB.GetQuestion(userReq.QuestionID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/question?question_id="+strconv.Itoa(userReq.QuestionID))
}

// HandlerPostRandomQuestion 重定向到随机选择的题目。
func (qs *QuestionServer) HandlerPostRandomQuestion(c *gin.Context) {
	session := NewRandomSession(qs.DB)
	if len(session.Questions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No questions available"})
		return
	}
	qs.RS[session.ID] = session
	c.Redirect(http.StatusFound, "/question/random/"+strconv.Itoa(session.ID))
}

// HandlerPostCheckAnswer 检查单题答案并返回 JSON 结果。
func (qs *QuestionServer) HandlerPostCheckAnswer(c *gin.Context) {
	var userReq Request
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	question, err := qs.DB.GetQuestion(userReq.QuestionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, NewResponse(userReq.Choice == question.Answer, question.Explanation))
}

// HandlerPostPracticeInit 创建练习并重定向到第一题。
func (qs *QuestionServer) HandlerPostPracticeInit(c *gin.Context) {
	var userReq Request
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	practice := (&Practice{}).NewPractice().GenerateExam(qs.DB, userReq.PracticeSize)
	qs.PM[practice.ID] = practice
	practice.Reset()
	c.Redirect(http.StatusFound, "/practice/"+strconv.Itoa(practice.ID))
}

// HandlerPostSubmitAnswer 保存练习中一道题的答案。
func (qs *QuestionServer) HandlerPostSubmitAnswer(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	practice := qs.PM[req.PracticeID]
	if practice == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Practice not found"})
		return
	}
	practice.Answers[req.QuestionID] = req.Choice
	c.JSON(http.StatusOK, gin.H{"message": "saved"})
}

// HandlerSubmitPractice 批改练习、标记完成状态并返回汇总结果。
func (qs *QuestionServer) HandlerSubmitPractice(c *gin.Context) {
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
	results := practice.CheckPractice(qs)
	practice.Completed = true
	c.JSON(http.StatusOK, results)
}
