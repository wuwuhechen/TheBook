package config

import "fmt"

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`

	FrontEnd FrontEndConfig `mapstructure:"front_end"`
}

type DatabaseConfig struct {
	DatabasePath string `mapstructure:"database_path"`
}

type FrontEndConfig struct {
	TemplatePath string `mapstructure:"template_path"`
}

func (c *Config) String() string {
	return fmt.Sprintf("Database Path: %s\nFront End Template Path: %s\n", c.Database.DatabasePath, c.FrontEnd.TemplatePath)
}
