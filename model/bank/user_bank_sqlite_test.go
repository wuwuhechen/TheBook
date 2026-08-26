package bank

import (
	"TheBook/auth"
	"testing"

	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserBankSQLiteLoginUpgradesMigratedPlaintextPassword(t *testing.T) {
	db, err := gorm.Open(sqliteDriver.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("Failed to migrate users table: %v", err)
	}

	user := &User{
		Username: "migrated_user",
		Password: "Correct123!",
		Nickname: "Migrated User",
		Role:     "user",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to insert migrated user: %v", err)
	}

	userBank := NewUserBankSQLite(db)
	loggedIn, err := userBank.LoginUser(LoginRequest{
		Username: user.Username,
		Password: "Correct123!",
	})
	if err != nil {
		t.Fatalf("Expected migrated user login to succeed: %v", err)
	}
	if !auth.VerifyPassword(loggedIn.Password, "Correct123!") {
		t.Fatal("Expected plaintext password to be upgraded to bcrypt")
	}

	var persisted User
	if err := db.Where("username = ?", user.Username).First(&persisted).Error; err != nil {
		t.Fatalf("Failed to reload migrated user: %v", err)
	}
	if !auth.VerifyPassword(persisted.Password, "Correct123!") {
		t.Fatal("Expected upgraded bcrypt password to be persisted")
	}
}
