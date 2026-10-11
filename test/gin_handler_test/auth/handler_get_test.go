package auth_test

import (
	"TheBook/test/gin_handler_test/shared"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAuthenticationPages 验证登录和注册页面 Handler。
func TestAuthenticationPages(t *testing.T) {
	f := shared.NewFixture(t)
	for _, target := range []string{"/auth/login", "/auth/register"} {
		w := shared.Serve(f, httptest.NewRequest(http.MethodGet, target, nil))
		if w.Code != http.StatusOK {
			t.Errorf("GET %s: expected 200, got %d", target, w.Code)
		}
	}
}
