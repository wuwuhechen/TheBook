package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig用于加载配置文件并将其解析为Config结构体
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

	viper.SetEnvPrefix("THEBOOK")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.BindEnv("sqlite.database_path")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &cfg, nil
}
