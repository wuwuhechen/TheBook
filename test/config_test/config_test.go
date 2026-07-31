package config_test

import (
	"TheBook/config"
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configPath := "../../config/config.yaml"

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Loaded config: %s", cfg)
}
