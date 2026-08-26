package bank

import (
	"TheBook/auth"
	"TheBook/model/manager"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

type UserBankSQLite struct {
	mu sync.Mutex
	db *gorm.DB
}

var _ manager.UserManager = (*UserBankSQLite)(nil)

func NewUserBankSQLite(db *gorm.DB) *UserBankSQLite {
	return &UserBankSQLite{
		db: db,
	}
}

func (ub *UserBankSQLite) AddUser(user *User) error {
	ub.mu.Lock()
	defer ub.mu.Unlock()
	return ub.db.Create(user).Error
}

func (ub *UserBankSQLite) SetPasswordHash(username, passwordHash string) error {
	ub.mu.Lock()
	defer ub.mu.Unlock()
	return ub.db.Model(&User{}).Where("username = ?", username).Update("password", passwordHash).Error
}

func (ub *UserBankSQLite) RegisterUser(req RegisterRequest) (*User, error) {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	exists := ub.db.Where("username = ?", req.Username).First(&User{}).Error
	if exists == nil {
		return nil, fmt.Errorf("user already exists")
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	user := &User{
		Username: req.Username,
		Password: passwordHash,
		Nickname: req.Nickname,
	}

	err = ub.AddUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to add user: %v", err)
	}

	return user, nil
}

func (ub *UserBankSQLite) GetUser(username string) (*User, bool) {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	var user User
	err := ub.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, false
	}

	return &user, true
}

func (ub *UserBankSQLite) LoginUser(req LoginRequest) (*User, error) {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	user, exist := ub.GetUser(req.Username)
	if !exist {
		return nil, fmt.Errorf("user not found")
	}

	passwordHash := user.Password
	if !auth.VerifyPassword(passwordHash, req.Password) {
		return nil, fmt.Errorf("incorrect password")
	}

	return user, nil
}

func (ub *UserBankSQLite) DeleteUser(username string) error {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	err := ub.db.Where("username = ?", username).Delete(&User{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}
