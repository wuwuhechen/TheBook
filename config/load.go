package config

import (
	"fmt"

	"github.com/spf13/viper"
)

/*
LoadConfig 加载配置文件

参数

	path string: 配置文件的路径

返回值

	*Config: 配置实例
*/
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

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
