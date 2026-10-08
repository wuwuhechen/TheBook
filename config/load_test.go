package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const validConfigYAML = `database:
  database_path: database/questions.json
front_end:
  template_path: front_end/*.html
user:
  user_bank_path: database/users.json
  user_hash_path: database/users_hash.json
  question_progress_path: database/question_progress.json
auth:
  jwt_secret: test-secret
  expiration_time: 3600
sqlite:
  database_path: database/test.db
`

func writeTestConfig(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write test config: %v", err)
	}
	return path
}

func TestLoadConfig(t *testing.T) {
	t.Setenv("THEBOOK_SQLITE_DATABASE_PATH", "")

	cfg, err := LoadConfig(writeTestConfig(t, validConfigYAML))
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Database.DatabasePath != "database/questions.json" {
		t.Errorf("Database.DatabasePath = %q, want %q", cfg.Database.DatabasePath, "database/questions.json")
	}
	if cfg.FrontEnd.TemplatePath != "front_end/*.html" {
		t.Errorf("FrontEnd.TemplatePath = %q, want %q", cfg.FrontEnd.TemplatePath, "front_end/*.html")
	}
	if cfg.User.UserBankPath != "database/users.json" || cfg.User.UserHashPath != "database/users_hash.json" || cfg.User.QuestionProgressPath != "database/question_progress.json" {
		t.Errorf("User config = %+v, want all configured paths", cfg.User)
	}
	if cfg.Auth.JWTSecret != "test-secret" || cfg.Auth.ExpirationTime != 3600 {
		t.Errorf("Auth config = %+v, want secret test-secret and expiration 3600", cfg.Auth)
	}
	if cfg.SQLite.DatabasePath != "database/test.db" {
		t.Errorf("SQLite.DatabasePath = %q, want %q", cfg.SQLite.DatabasePath, "database/test.db")
	}
}

func TestLoadConfigUsesSQLiteEnvironmentOverride(t *testing.T) {
	t.Setenv("THEBOOK_SQLITE_DATABASE_PATH", "database/overridden.db")

	cfg, err := LoadConfig(writeTestConfig(t, validConfigYAML))
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg.SQLite.DatabasePath != "database/overridden.db" {
		t.Errorf("SQLite.DatabasePath = %q, want environment override %q", cfg.SQLite.DatabasePath, "database/overridden.db")
	}
}

func TestLoadConfigReturnsErrorForMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.yaml")

	if _, err := LoadConfig(path); err == nil {
		t.Fatal("LoadConfig() error = nil, want error for missing file")
	}
}

func TestLoadConfigReturnsErrorForInvalidYAML(t *testing.T) {
	path := writeTestConfig(t, "sqlite:\n  database_path: [unterminated\n")

	if _, err := LoadConfig(path); err == nil {
		t.Fatal("LoadConfig() error = nil, want error for invalid YAML")
	}
}

func TestConfigString(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{DatabasePath: "database/questions.json"},
		FrontEnd: FrontEndConfig{TemplatePath: "front_end/*.html"},
		User:     UserConfig{UserBankPath: "database/users.json"},
		Auth:     AuthConfig{JWTSecret: "test-secret", ExpirationTime: 3600},
	}

	got := cfg.String()
	for _, expected := range []string{
		"Database Path: database/questions.json",
		"Front End Template Path: front_end/*.html",
		"User Bank Path: database/users.json",
		"JWTSecret:test-secret",
		"ExpirationTime:3600",
	} {
		if !strings.Contains(got, expected) {
			t.Errorf("Config.String() = %q, want substring %q", got, expected)
		}
	}
}
