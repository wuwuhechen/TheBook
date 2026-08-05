package config

import "fmt"

type Config struct {
	// Database 字段用于存储数据库相关的配置
	Database DatabaseConfig `mapstructure:"database"`

	// FrontEnd 字段用于存储前端相关的配置
	FrontEnd FrontEndConfig `mapstructure:"front_end"`

	// User 字段用于存储用户相关的配置
	User UserConfig `mapstructure:"user"`

	// Auth 字段用于存储认证相关的配置
	Auth AuthConfig `mapstructure:"auth"`
}

type DatabaseConfig struct {
	// DatabasePath 字段用于存储数据库文件的路径
	DatabasePath string `mapstructure:"database_path"`
}

type FrontEndConfig struct {
	// TemplatePath 字段用于存储前端模板文件的路径
	TemplatePath string `mapstructure:"template_path"`
}

type UserConfig struct {
	// UserBankPath 字段用于存储用户数据文件的路径
	UserBankPath string `mapstructure:"user_bank_path"`
}

type AuthConfig struct {
	// JWTSecret 字段用于存储JWT的密钥
	JWTSecret string `mapstructure:"jwt_secret"`

	// ExpirationTime 字段用于存储JWT的过期时间（以秒为单位）
	ExpirationTime int64 `mapstructure:"expiration_time"`
}

func (c *Config) String() string {
	return fmt.Sprintf("Database Path: %s\nFront End Template Path: %s\nUser Bank Path: %s\nAuth Config: %+v\n",
		c.Database.DatabasePath, c.FrontEnd.TemplatePath, c.User.UserBankPath, c.Auth)
}
