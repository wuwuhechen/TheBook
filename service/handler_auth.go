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

// HandlerGetRegisterPage 渲染用户注册页面。
func (s *Server) HandlerGetRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register_page.html", nil)
}

// HandlerGetLoginPage 渲染用户登录页面。
func (s *Server) HandlerGetLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login_page.html", nil)
}

// HandlerPostLogin 校验账户并建立登录 Cookie。
func (s *Server) HandlerPostLogin(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := s.UM.LoginUser(req)
	if err != nil {
		s.businessLog().Warn("用户登录失败", zap.String("username", req.Username))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	token, err := auth.GenerateToken(user.UserID, user.Username, user.Role, user.Nickname, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.SetCookie("access_token", token, 24*3600, "/", "", false, true)
	c.SetCookie("userID", strconv.Itoa(int(user.UserID)), 24*3600, "/", "", false, true)
	s.businessLog().Info("用户登录成功", zap.Uint("user_id", user.UserID), zap.String("username", user.Username))
	c.Redirect(http.StatusFound, "/")
}

// HandlerPostRegister 创建账户并在成功后写入登录 Cookie。
func (s *Server) HandlerPostRegister(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := s.UM.RegisterUser(req)
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
	token, err := auth.GenerateToken(user.UserID, user.Username, user.Role, user.Nickname, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.Set("userID", user.UserID)
	c.SetCookie("access_token", token, 24*3600, "/", "", false, true)
	s.businessLog().Info("用户注册成功", zap.Uint("user_id", user.UserID), zap.String("username", user.Username))
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful", "token": token})
}

// HandlerPostLogout 清除认证 Cookie。
func (s *Server) HandlerPostLogout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	s.businessLog().Info("用户已退出登录")
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// persistUsers 将内存用户仓库及其密码哈希持久化到 JSON 文件。
func (s *Server) persistUsers() error {
	if _, ok := s.UM.(*model.UserBankSQLite); ok {
		return nil
	}
	bank, ok := s.UM.(*model.UserBank)
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
	sort.Slice(users, func(i, j int) bool { return users[i].UserID < users[j].UserID })
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

// persistUserPasswordHashes 将用户密码哈希写入独立文件。
func persistUserPasswordHashes(path string, bank *model.UserBank) error {
	hashUsers := make([]*model.User, 0, len(bank.Users))
	for _, user := range bank.Users {
		passwordHash, exists := bank.PasswordHashes[user.Username]
		if !exists {
			return fmt.Errorf("password hash not found for user %s", user.Username)
		}
		hashUsers = append(hashUsers, &model.User{UserID: user.UserID, Username: user.Username, Password: passwordHash, Nickname: user.Nickname, Role: user.Role})
	}
	sort.Slice(hashUsers, func(i, j int) bool { return hashUsers[i].UserID < hashUsers[j].UserID })
	data, err := json.MarshalIndent(hashUsers, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
