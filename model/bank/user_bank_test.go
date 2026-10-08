package bank

import "testing"

func TestUserBankRegisterLoginAndDelete(t *testing.T) {
	bank := NewUserBank()
	if err := bank.AddUser(nil); err == nil {
		t.Fatal("expected nil user error")
	}
	if err := bank.AddUser(&User{UserID: 4, Username: "existing"}); err != nil {
		t.Fatalf("add existing user: %v", err)
	}
	user, err := bank.RegisterUser(RegisterRequest{Username: "new", Password: "Correct123!", Nickname: "New User"})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if user.UserID != 5 || user.Password != "Correct123!" {
		t.Fatalf("unexpected registered user: %+v", user)
	}
	if _, ok := bank.PasswordHashes[user.Username]; !ok {
		t.Fatal("expected password hash to be stored")
	}
	if _, err := bank.RegisterUser(RegisterRequest{Username: "new", Password: "again"}); err == nil {
		t.Fatal("expected duplicate user error")
	}
	if _, err := bank.LoginUser(LoginRequest{Username: "new", Password: "wrong"}); err == nil {
		t.Fatal("expected incorrect-password error")
	}
	loggedIn, err := bank.LoginUser(LoginRequest{Username: "new", Password: "Correct123!"})
	if err != nil || loggedIn.Username != "new" {
		t.Fatalf("unexpected login result: %+v, %v", loggedIn, err)
	}
	if err := bank.DeleteUser("new"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, ok := bank.GetUser("new"); ok {
		t.Fatal("user should be removed")
	}
	if err := bank.DeleteUser("new"); err == nil {
		t.Fatal("expected missing-user delete error")
	}
}
