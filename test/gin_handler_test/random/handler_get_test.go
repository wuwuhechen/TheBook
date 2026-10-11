package random_test

import (
	"TheBook/model"
	"TheBook/test/gin_handler_test/shared"
	"net/http"
	"strings"
	"testing"
)

// TestRandomGetNavigation 验证随机练习页面和导航 Handler。
func TestRandomGetNavigation(t *testing.T) {
	f := shared.NewFixture(t)
	session := &model.RandomSession{ID: 88, Questions: []int{1, 2}}
	if err := f.Server.RS.Create(session); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Server.RS.Delete(session.ID) })
	for _, target := range []string{"/question/random/not-a-number", "/question/random/404"} {
		w := shared.Serve(f, f.Request(http.MethodGet, target, nil))
		if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
			t.Fatalf("GET %s: got %d", target, w.Code)
		}
	}
	w := shared.Serve(f, f.Request(http.MethodGet, "/question/random/88?direction=next", nil))
	if w.Code != http.StatusOK || session.CurrentIndex != 1 || !strings.Contains(w.Body.String(), "第二题") {
		t.Fatalf("expected second question, got %d index=%d", w.Code, session.CurrentIndex)
	}
	w = shared.Serve(f, f.Request(http.MethodGet, "/question/random/88?direction=next", nil))
	if session.CurrentIndex != 1 {
		t.Fatalf("navigation exceeded last question: %d", session.CurrentIndex)
	}
}
