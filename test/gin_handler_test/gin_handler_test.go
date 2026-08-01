package gin_handler_test

import (
	"TheBook/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	router         *gin.Engine
	questionServer *service.QuestionServer
)

func TestMain(M *testing.M) {
	r, qs, err := service.InitSystem()
	if err != nil {
		panic(err)
	}

	router = r
	questionServer = qs

	code := M.Run()

	os.Exit(code)
}

func TestInitQuestionPage(t *testing.T) {

	request := service.Request{UserID: "test_user", QuestionID: 1, Choice: 5}
	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", "/request", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf(
			"expected 302 got %d",
			w.Code,
		)
	}

	location :=
		w.Header().Get("Location")

	// 第二次请求 GET
	req2, _ := http.NewRequest(
		"GET",
		location,
		nil,
	)

	w2 := httptest.NewRecorder()

	router.ServeHTTP(
		w2,
		req2,
	)

	if w2.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			w2.Code,
		)
	}

	os.WriteFile(
		"test_result.html",
		w2.Body.Bytes(),
		0644,
	)

}

func TestPostQuestionRedirect(t *testing.T) {

	request := service.Request{
		UserID:     "test_user",
		QuestionID: 1,
		Choice:     5,
	}

	jsonData, err := json.Marshal(request)

	if err != nil {
		t.Fatalf(
			"Failed to marshal request: %v",
			err,
		)
	}

	req, err := http.NewRequest(
		"POST",
		"/request",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(
		w,
		req,
	)

	// POST应该返回302
	if w.Code != http.StatusFound {

		t.Fatalf(
			"Expected status 302, got %d",
			w.Code,
		)
	}

	// 检查重定向地址

	location :=
		w.Header().Get("Location")

	expected :=
		"/question?question_id=1"

	if location != expected {

		t.Fatalf(
			"Expected Location %s, got %s",
			expected,
			location,
		)
	}

	t.Logf(
		"Redirect success: %s",
		location,
	)
}

func TestGetQuestionPage(t *testing.T) {

	req, err := http.NewRequest(
		"GET",
		"/question?question_id=1",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	router.ServeHTTP(
		w,
		req,
	)

	if w.Code != http.StatusOK {

		t.Fatalf(
			"Expected status 200, got %d",
			w.Code,
		)
	}

	body :=
		w.Body.String()

	// 简单检查页面内容

	if !strings.Contains(
		body,
		"题目",
	) {

		t.Fatal(
			"Response does not contain question content",
		)
	}

	// 保存HTML方便查看

	err = os.WriteFile(
		"test_result.html",
		w.Body.Bytes(),
		0644,
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(
		"HTML saved: test_result.html",
	)
}

func TestCheckAnswerTrue(t *testing.T) {
	request := service.Request{
		UserID:     "test_user",
		QuestionID: 1,
		Choice:     1,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", "/check_answer", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	//返回是json文件，进行解析
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["correct"] == nil {
		t.Fatalf("Expected 'correct' field in response, got nil")
	}

	if response["correct"].(bool) != true {
		t.Fatalf("Expected 'correct' to be true, got %v", response["correct"])
	}
}

func TestCheckAnswerFalse(t *testing.T) {
	request := service.Request{
		UserID:     "test_user",
		QuestionID: 1,
		Choice:     2,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", "/check_answer", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	//返回是json文件，进行解析
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["correct"] == nil {
		t.Fatalf("Expected 'correct' field in response, got nil")
	}

	if response["correct"].(bool) != false {
		t.Fatalf("Expected 'correct' to be false, got %v", response["correct"])
	}
}

func TestGetHomePage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", w.Code)
	}

	file := "test_home_page.html"

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
