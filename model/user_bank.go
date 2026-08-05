package model

import "fmt"

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
	Users map[string]*User `json:"users"`
	// NextID 是下一个可分配的用户编号。
	// 它在加载已有用户时同步推进，避免每次注册都扫描整个用户表。
	NextID uint `json:"next_id"`
}

func NewUserBank() *UserBank {
	return &UserBank{
		Users:  make(map[string]*User),
		NextID: 0,
	}
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
	ub.AddUser(user)
	return user, nil
}

func (ub *UserBank) LoginUser(req LoginRequest) (*User, error) {
	user, exists := ub.GetUser(req.Username)
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	if user.Password != req.Password {
		return nil, fmt.Errorf("incorrect password")
	}
	return user, nil
}
