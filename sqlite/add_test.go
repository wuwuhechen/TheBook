package sqlite

import (
	"testing"

	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAddUserPreservesZeroUserID(t *testing.T) {
	db, err := gorm.Open(sqliteDriver.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("Failed to migrate users table: %v", err)
	}

	user := &User{
		UserID:   0,
		Username: "admin",
		Password: "password-hash",
		Nickname: "Administrator",
		Role:     "admin",
	}
	if err := AddUser(db, user); err != nil {
		t.Fatalf("Failed to add user: %v", err)
	}

	var stored User
	if err := db.Where("user_id = ?", 0).First(&stored).Error; err != nil {
		t.Fatalf("Expected user ID 0 to be stored: %v", err)
	}
	if stored.Username != "admin" {
		t.Fatalf("Expected admin user, got %q", stored.Username)
	}
}
