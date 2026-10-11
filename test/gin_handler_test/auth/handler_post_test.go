package auth_test

import (
	"TheBook/model"
	"TheBook/test/gin_handler_test/shared"
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestAuthenticationPosts 验证注册、登录和退出 Handler。
func TestAuthenticationPosts(t *testing.T) {
	f := shared.NewFixture(t)
	w := shared.Serve(f, httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte(`{`))))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid registration: got %d", w.Code)
	}
	register := model.RegisterRequest{Username: "new_user", Password: "password", Nickname: "New User"}
	w = shared.Serve(f, httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(shared.JSONBody(t, register))))
	if w.Code != http.StatusOK || !strings.Contains(w.Header().Get("Set-Cookie"), "access_token=") {
		t.Fatalf("registration failed: %d: %s", w.Code, w.Body.String())
	}
	w = shared.Serve(f, httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(shared.JSONBody(t, register))))
	if w.Code != http.StatusConflict {
		t.Fatalf("duplicate registration: got %d", w.Code)
	}
	w = shared.Serve(f, httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(shared.JSONBody(t, model.LoginRequest{Username: "new_user", Password: "bad"}))))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("invalid login: got %d", w.Code)
	}
	w = shared.Serve(f, httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(shared.JSONBody(t, model.LoginRequest{Username: "new_user", Password: "password"}))))
	if w.Code != http.StatusFound || w.Header().Get("Location") != "/" {
		t.Fatalf("valid login: got %d %q", w.Code, w.Header().Get("Location"))
	}
	w = shared.Serve(f, httptest.NewRequest(http.MethodPost, "/auth/logout", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Header().Get("Set-Cookie"), "Max-Age=0") {
		t.Fatalf("logout failed: %d %q", w.Code, w.Header().Get("Set-Cookie"))
	}
}
