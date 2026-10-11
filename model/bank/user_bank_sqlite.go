package bank

import (
	"TheBook/auth"
	"TheBook/model/manager"
	"crypto/subtle"
	"fmt"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type UserBankSQLite struct {
	mu sync.Mutex
	db *gorm.DB
}

var _ manager.UserManager = (*UserBankSQLite)(nil)

// NewUserBankSQLite 创建 SQLite 用户仓库。
func NewUserBankSQLite(db *gorm.DB) *UserBankSQLite {
	return &UserBankSQLite{
		db: db,
	}
}

// AddUser 将用户写入 SQLite。
func (ub *UserBankSQLite) AddUser(user *User) error {
	ub.mu.Lock()
	defer ub.mu.Unlock()
	return ub.db.Create(user).Error
}

// SetPasswordHash 更新指定用户的密码哈希。
func (ub *UserBankSQLite) SetPasswordHash(username, passwordHash string) error {
	ub.mu.Lock()
	defer ub.mu.Unlock()
	return ub.db.Model(&User{}).Where("username = ?", username).Update("password", passwordHash).Error
}

// RegisterUser 创建并持久化 SQLite 用户账户。
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
		Role:     "user",
	}

	err = ub.db.Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to add user: %v", err)
	}

	return user, nil
}

// GetUser 按用户名从 SQLite 查询用户。
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

// LoginUser 验证 SQLite 用户密码并升级旧明文密码。
func (ub *UserBankSQLite) LoginUser(req LoginRequest) (*User, error) {
	user, exist := ub.GetUser(req.Username)
	if !exist {
		return nil, fmt.Errorf("user not found")
	}

	if isBcryptHash(user.Password) {
		if !auth.VerifyPassword(user.Password, req.Password) {
			return nil, fmt.Errorf("incorrect password")
		}
		return user, nil
	}

	// 处理非 bcrypt 密码哈希的情况
	if subtle.ConstantTimeCompare([]byte(user.Password), []byte(req.Password)) != 1 {
		return nil, fmt.Errorf("incorrect password")
	}
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade password: %v", err)
	}
	if err := ub.SetPasswordHash(user.Username, passwordHash); err != nil {
		return nil, fmt.Errorf("failed to upgrade password: %v", err)
	}
	user.Password = passwordHash

	return user, nil
}

// isBcryptHash 判断字符串是否具有 bcrypt 哈希前缀。
func isBcryptHash(password string) bool {
	return strings.HasPrefix(password, "$2a$") ||
		strings.HasPrefix(password, "$2b$") ||
		strings.HasPrefix(password, "$2y$")
}

// DeleteUser 从 SQLite 删除指定用户。
func (ub *UserBankSQLite) DeleteUser(username string) error {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	err := ub.db.Where("username = ?", username).Delete(&User{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}
