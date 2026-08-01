package config

import "fmt"

type Config struct {
	// Database 字段用于存储数据库相关的配置
	Database DatabaseConfig `mapstructure:"database"`

	// FrontEnd 字段用于存储前端相关的配置
	FrontEnd FrontEndConfig `mapstructure:"front_end"`
}

type DatabaseConfig struct {
	// DatabasePath 字段用于存储数据库文件的路径
	DatabasePath string `mapstructure:"database_path"`
}

type FrontEndConfig struct {
	// TemplatePath 字段用于存储前端模板文件的路径
	TemplatePath string `mapstructure:"template_path"`
}

func (c *Config) String() string {
	return fmt.Sprintf("Database Path: %s\nFront End Template Path: %s\n", c.Database.DatabasePath, c.FrontEnd.TemplatePath)
}
