package gin_handler_test

import (
	"TheBook/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(M *testing.M) {

}

func TestGet(t *testing.T) {
	r := gin.Default()
	// r.GET("/request", service.HandlerPOST)

	r.LoadHTMLFiles("../front/question_page.html")

	_, err := service.LoadQuestions("../database/data.json")

	request := service.Request{
		UserID:     "test_user",
		QuestionID: 1,
		Choice:     2,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req, _ := http.NewRequest("GET", "/request", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", w.Code)
	}

	t.Logf("Response: %s", w.Body.String())

	file := "test_result.html"

	os.WriteFile(
		file,
		w.Body.Bytes(),
		0644,
	)

	t.Log(
		"HTML saved:",
		file,
	)
}
