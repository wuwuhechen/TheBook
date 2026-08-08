package service

import (
	"TheBook/auth"
	"TheBook/model"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDValue.(uint)

	userProgress, exists := s.QS[userID]
	if !exists {
		s.QS[userID] = &model.QuestionProgress{
			UserID:            userID,
			CurrentQuestionID: userReq.QuestionID,
			UpdatedAt:         time.Now(),
		}
	} else {
		userProgress.CurrentQuestionID = userReq.QuestionID
		userProgress.UpdatedAt = time.Now()
	}

	err := s.persistQuestionProgress()
	if err != nil {
		s.appLog().Error("保存顺序答题进度失败", zap.Uint("user_id", userID), zap.Int("question_id", userReq.QuestionID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save question progress"})
		return
	}
	s.businessLog().Info("顺序答题进度已更新", zap.Uint("user_id", userID), zap.Int("question_id", userReq.QuestionID))

	c.Redirect(http.StatusFound, "/question?question_id="+strconv.Itoa(userReq.QuestionID))
}

func (s *Server) HandlerPostCurrentQuestion(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDValue.(uint)

	userProgress, exists := s.QS[userID]
	if !exists {
		userProgress = &model.QuestionProgress{UserID: userID}
		s.QS[userID] = userProgress
	}

	currentQuestionID := userProgress.CurrentQuestionID
	if currentQuestionID == 0 {
		var ok bool
		currentQuestionID, ok = s.firstQuestionID()
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "No questions available"})
			return
		}
	} else if _, err := s.DB.GetQuestion(currentQuestionID); err != nil {
		var ok bool
		currentQuestionID, ok = s.firstQuestionID()
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "No questions available"})
			return
		}
	}
	userProgress.CurrentQuestionID = currentQuestionID
	userProgress.UpdatedAt = time.Now()

	err := s.persistQuestionProgress()
	if err != nil {
		s.appLog().Error("恢复顺序答题进度后保存失败", zap.Uint("user_id", userID), zap.Int("question_id", currentQuestionID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save question progress"})
		return
	}
	s.businessLog().Info("顺序答题进度已恢复", zap.Uint("user_id", userID), zap.Int("question_id", currentQuestionID))

	c.Redirect(http.StatusFound, "/question?question_id="+strconv.Itoa(currentQuestionID))
}

// firstQuestionID 返回题库中编号最小的有效题目 ID。
func (s *Server) firstQuestionID() (int, bool) {
	questions := s.DB.GetAllQuestionsSorted()
	if len(questions) == 0 {
		return 0, false
	}
	return questions[0].ID, true
}

func (s *Server) persistQuestionProgress() error {
	s.RecordMu.Lock()
	defer s.RecordMu.Unlock()

	progresses := make([]*model.QuestionProgress, 0, len(s.QS))
	for _, progress := range s.QS {
		progresses = append(progresses, progress)
	}

	data, err := json.MarshalIndent(progresses, "", "  ")
	if err != nil {
		return err
	}

	path := s.RootPath + "/database/question_progress.json"
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// HandlerPostRandomQuestion 重定向到随机选择的题目。
func (s *Server) HandlerPostRandomQuestion(c *gin.Context) {
	session := model.NewRandomSession(s.DB)
	if len(session.Questions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No questions available"})
		return
	}
	s.RS[session.ID] = session
	s.businessLog().Info("随机练习已创建", zap.Int("session_id", session.ID), zap.Int("question_count", len(session.Questions)))
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
	correct := userReq.Choice == question.Answer
	s.businessLog().Info("单题判题完成", zap.Int("question_id", userReq.QuestionID), zap.Bool("correct", correct))
	c.JSON(http.StatusOK, model.NewResponse(correct, question.Explanation))
}

// HandlerPostPracticeInit 创建练习并重定向到第一题。
func (s *Server) HandlerPostPracticeInit(c *gin.Context) {
	var userReq model.Request
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	practice := (&model.Practice{}).NewPractice().GenerateExam(s.DB, userReq.PracticeSize)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	practice.UserID = userID.(uint)
	s.PM[practice.ID] = practice
	practice.Reset()
	s.businessLog().Info("套题已创建", zap.Uint("user_id", practice.UserID), zap.Int("practice_id", practice.ID), zap.Int("question_count", practice.TotalQuestions))
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
	s.businessLog().Info("套题答案已保存", zap.Int("practice_id", req.PracticeID), zap.Int("question_id", req.QuestionID))
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
	userIDValue, exists := c.Get("userID")
	if !exists || practice.UserID != userIDValue.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No permission to submit this practice"})
		return
	}
	if practice.Completed {
		c.JSON(http.StatusConflict, gin.H{"error": "Practice already submitted"})
		return
	}
	results := practice.CheckPractice(s.DB)
	if err := s.persistPracticeRecord(practice, results); err != nil {
		s.appLog().Error("保存套题记录失败", zap.Uint("user_id", practice.UserID), zap.Int("practice_id", practice.ID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save practice record"})
		return
	}
	practice.Completed = true
	s.businessLog().Info("套题已提交", zap.Uint("user_id", practice.UserID), zap.Int("practice_id", practice.ID), zap.Int("correct_count", results.CorrectCount), zap.Int("wrong_count", results.WrongCount))
	c.JSON(http.StatusOK, results)
}

func (s *Server) persistPracticeRecord(practice *model.Practice, results *model.PracticeResponse) error {
	s.RecordMu.Lock()
	defer s.RecordMu.Unlock()
	path := s.RecordPath
	if path == "" {
		path = "database/practice_records.json"
	}
	var records []model.PracticeRecord
	if data, err := os.ReadFile(path); err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &records); err != nil {
			return err
		}
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	answers := make([]model.AnswerRecord, 0, len(practice.Questions))
	correct := make(map[int]bool)
	for _, id := range practice.Questions {
		q, err := s.DB.GetQuestion(id)
		if err != nil {
			continue
		}
		answer, answered := practice.Answers[id]
		correct[id] = answered && answer == q.Answer
		answers = append(answers, model.AnswerRecord{QuestionID: id, Answer: answer, Answered: answered, Correct: correct[id]})
	}
	records = append(records, model.PracticeRecord{ID: len(records) + 1, UserID: practice.UserID, PracticeID: practice.ID, TotalQuestions: results.Total, CorrectCount: results.CorrectCount, WrongCount: results.WrongCount, StartTime: practice.StartTime, SubmitTime: time.Now(), Answers: answers})
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *Server) HandlerPostLogin(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.US.LoginUser(req)
	if err != nil {
		s.businessLog().Warn("用户登录失败", zap.String("username", req.Username))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.GenerateToken(user.UserID, user.Username, user.Role, user.Nickname, time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("access_token", token, 3600*24, "/", "", false, true)
	s.businessLog().Info("用户登录成功", zap.Uint("user_id", user.UserID), zap.String("username", user.Username))
	c.SetCookie("userID", strconv.Itoa(int(user.UserID)), 3600*24, "/", "", false, true)
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
		s.businessLog().Warn("用户注册失败", zap.String("username", req.Username))
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if err := s.persistUsers(); err != nil {
		s.appLog().Error("保存用户账户失败", zap.Uint("user_id", user.UserID), zap.String("username", user.Username), zap.Error(err))
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
	s.businessLog().Info("用户注册成功", zap.Uint("user_id", user.UserID), zap.String("username", user.Username))
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful", "token": token})
}

func (s *Server) persistUsers() error {
	bank, ok := s.US.(*model.UserBank)
	if !ok {
		return fmt.Errorf("unsupported user service type")
	}

	usersByName := make(map[string]*model.User)
	if data, err := os.ReadFile(s.UserPath); err == nil && len(data) > 0 {
		var existingUsers []*model.User
		if err := json.Unmarshal(data, &existingUsers); err != nil {
			return err
		}
		for _, user := range existingUsers {
			if user != nil {
				usersByName[user.Username] = user
			}
		}
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}

	for username, user := range bank.Users {
		usersByName[username] = user
	}
	users := make([]*model.User, 0, len(usersByName))
	for _, user := range usersByName {
		users = append(users, user)
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].UserID < users[j].UserID
	})
	for _, user := range usersByName {
		bank.AddUser(user)
	}
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(s.UserPath, data, 0644); err != nil {
		return err
	}
	return persistUserPasswordHashes(s.UserHashPath, bank)
}

// persistUserPasswordHashes 将用户资料和对应密码哈希写入独立的哈希账户文件。
func persistUserPasswordHashes(path string, bank *model.UserBank) error {
	hashUsers := make([]*model.User, 0, len(bank.Users))
	for _, user := range bank.Users {
		passwordHash, exists := bank.PasswordHashes[user.Username]
		if !exists {
			return fmt.Errorf("password hash not found for user %s", user.Username)
		}
		hashUsers = append(hashUsers, &model.User{
			UserID:   user.UserID,
			Username: user.Username,
			Password: passwordHash,
			Nickname: user.Nickname,
			Role:     user.Role,
		})
	}
	sort.Slice(hashUsers, func(i, j int) bool {
		return hashUsers[i].UserID < hashUsers[j].UserID
	})
	data, err := json.MarshalIndent(hashUsers, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *Server) HandlerPostLogout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	s.businessLog().Info("用户已退出登录")
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
