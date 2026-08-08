package model

import (
	"TheBook/auth"
	"fmt"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserBank struct {
	Users          map[string]*User  `json:"users"`
	PasswordHashes map[string]string `json:"-"`
	// NextID 是下一个可分配的用户编号。
	// 它在加载已有用户时同步推进，避免每次注册都扫描整个用户表。
	NextID uint `json:"next_id"`
}

func NewUserBank() *UserBank {
	return &UserBank{
		Users:          make(map[string]*User),
		PasswordHashes: make(map[string]string),
		NextID:         0,
	}
}

// SetPasswordHash 保存用户对应的 bcrypt 密码哈希。
func (ub *UserBank) SetPasswordHash(username, passwordHash string) {
	ub.PasswordHashes[username] = passwordHash
}

func (ub *UserBank) AddUser(user *User) {
	if user == nil {
		return
	}
	ub.Users[user.Username] = user
	if user.UserID >= ub.NextID {
		ub.NextID = user.UserID + 1
	}
}

func (ub *UserBank) GetUser(username string) (*User, bool) {
	user, exists := ub.Users[username]
	return user, exists
}

func (ub *UserBank) RegisterUser(req RegisterRequest) (*User, error) {
	if _, exists := ub.Users[req.Username]; exists {
		return nil, fmt.Errorf("user already exists")
	}

	user := &User{
		UserID:   ub.NextID,
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
	}
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}
	ub.AddUser(user)
	ub.SetPasswordHash(user.Username, passwordHash)
	return user, nil
}

func (ub *UserBank) LoginUser(req LoginRequest) (*User, error) {
	user, exists := ub.GetUser(req.Username)
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	passwordHash, exists := ub.PasswordHashes[req.Username]
	if !exists || !auth.VerifyPassword(passwordHash, req.Password) {
		return nil, fmt.Errorf("incorrect password")
	}

	return user, nil
}
