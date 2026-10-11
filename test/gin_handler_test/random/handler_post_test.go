package random_test

import (
	"TheBook/test/gin_handler_test/shared"
	"net/http"
	"strings"
	"testing"
)

// TestRandomPostCreateSession 验证随机练习会话创建 Handler。
func TestRandomPostCreateSession(t *testing.T) {
	f := shared.NewFixture(t)
	w := shared.Serve(f, f.Request(http.MethodPost, "/question/random", nil))
	if w.Code != http.StatusFound || !strings.HasPrefix(w.Header().Get("Location"), "/question/random/") {
		t.Fatalf("unexpected random redirect: %d %q", w.Code, w.Header().Get("Location"))
	}
}
