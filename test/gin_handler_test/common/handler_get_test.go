package common_test

import (
	"TheBook/test/gin_handler_test/shared"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHomePage 验证首页 Handler 可正常渲染。
func TestHomePage(t *testing.T) {
	f := shared.NewFixture(t)
	w := shared.Serve(f, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected home page 200, got %d", w.Code)
	}
}
