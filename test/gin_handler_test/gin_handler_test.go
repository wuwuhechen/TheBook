package gin_handler_test

import (
	"TheBook/auth"
	"TheBook/logger"
	"TheBook/model"
	"TheBook/service"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	router         *gin.Engine
	questionServer *service.Server
	authToken      string
)

func TestMain(M *testing.M) {
	log, closeLogger, err := logger.InitLogger("logs_test/app.log", false)
	if err != nil {
		panic(err)
	}
	defer closeLogger()

	server, err := service.InitSystem(log)
	if err != nil {
		panic(err)
	}

	router = server.Router
	questionServer = server
	authToken, err = auth.GenerateToken(1, "test_user", "user", "Test User", time.Hour)
	if err != nil {
		panic(err)
	}

	code := M.Run()

	os.Exit(code)
}

func newAuthenticatedRequest(method, target string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, target, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	return req, nil
}

func TestInitQuestionPage(t *testing.T) {

	request := model.Request{UserID: "test_user", QuestionID: 1, Choice: 5}
	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req, err := newAuthenticatedRequest("POST", "/question/request", bytes.NewBuffer(jsonData))
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
	req2, _ := newAuthenticatedRequest(
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

	request := model.Request{
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

	req, err := newAuthenticatedRequest(
		"POST",
		"/question/request",
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

	req, err := newAuthenticatedRequest(
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
	request := model.Request{
		UserID:     "test_user",
		QuestionID: 1,
		Choice:     1,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req, err := newAuthenticatedRequest("POST", "/question/check_answer", bytes.NewBuffer(jsonData))
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
	request := model.Request{
		UserID:     "test_user",
		QuestionID: 1,
		Choice:     2,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req, err := newAuthenticatedRequest("POST", "/question/check_answer", bytes.NewBuffer(jsonData))
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
	req, err := newAuthenticatedRequest("GET", "/", nil)
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

func TestPostRandomQuestion(t *testing.T) {
	req, err := newAuthenticatedRequest("POST", "/question/random", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("Expected status code 302, got %d", w.Code)
	}
	if !strings.HasPrefix(w.Header().Get("Location"), "/question/random/") {
		t.Fatalf("Expected random-session redirect, got %q", w.Header().Get("Location"))
	}
}

func TestGenerateExam(t *testing.T) {
	request := model.Request{
		UserID:       "test_user",
		PracticeSize: 5,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req, err := newAuthenticatedRequest("POST", "/practice/init", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("Expected status code 302, got %d", w.Code)
	}

	t.Logf("GenerateExam response: %s", w.Body.String())
}

func TestHandlerPostSubmitAnswer(t *testing.T) {
	const practiceID = 10001
	question, err := questionServer.DB.GetQuestion(1)
	if err != nil {
		t.Fatalf("Failed to get question: %v", err)
	}

	practice := (&model.Practice{}).NewPractice()
	practice.ID = practiceID
	practice.Questions = []int{question.ID}
	if err := questionServer.PM.Create(practice); err != nil {
		t.Fatalf("Failed to create practice: %v", err)
	}
	t.Cleanup(func() {
		if err := questionServer.PM.Delete(practice); err != nil {
			t.Errorf("Failed to delete practice: %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		request := model.Request{
			PracticeID: practiceID,
			QuestionID: question.ID,
			Choice:     question.Answer,
		}
		jsonData, err := json.Marshal(request)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		req, err := newAuthenticatedRequest("POST", "/practice/answer", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d: %s", w.Code, w.Body.String())
		}

		var response struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response.Message != "saved" {
			t.Fatalf("Expected message saved, got %q", response.Message)
		}
		if got := practice.Answers[question.ID]; got != question.Answer {
			t.Fatalf("Expected saved answer %d, got %d", question.Answer, got)
		}
	})

	t.Run("practice not found", func(t *testing.T) {
		body := `{"practice_id":999999,"question_id":1,"choice":1}`
		req, err := newAuthenticatedRequest("POST", "/practice/answer", strings.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected status code 404, got %d", w.Code)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req, err := newAuthenticatedRequest("POST", "/practice/answer", strings.NewReader(`{"practice_id":`))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status code 400, got %d", w.Code)
		}
	})
}

func TestHandlerSubmitPractice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const practiceID = 10002
		question, err := questionServer.DB.GetQuestion(1)
		if err != nil {
			t.Fatalf("Failed to get question: %v", err)
		}

		practice := (&model.Practice{}).NewPractice()
		practice.ID = practiceID
		practice.UserID = 1
		practice.Questions = []int{question.ID}
		practice.Answers[question.ID] = question.Answer
		if err := questionServer.PM.Create(practice); err != nil {
			t.Fatalf("Failed to create practice: %v", err)
		}
		originalRecordPath := questionServer.RecordPath
		questionServer.RecordPath = t.TempDir() + "/practice_records.json"
		t.Cleanup(func() {
			questionServer.RecordPath = originalRecordPath
			if err := questionServer.PM.Delete(practice); err != nil {
				t.Errorf("Failed to delete practice: %v", err)
			}
		})

		req, err := newAuthenticatedRequest("POST", "/practice/10002/submit", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d: %s", w.Code, w.Body.String())
		}

		var response model.PracticeResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response.Total != 1 || response.CorrectCount != 1 || response.WrongCount != 0 {
			t.Fatalf("Unexpected practice result: %+v", response)
		}
		if len(response.Details) != 1 || !response.Details[0].Correct {
			t.Fatalf("Expected one correct answer detail, got %+v", response.Details)
		}
		if !practice.Completed {
			t.Fatal("Expected practice to be marked completed")
		}
	})

	t.Run("invalid practice ID", func(t *testing.T) {
		req, err := newAuthenticatedRequest("POST", "/practice/not-a-number/submit", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status code 400, got %d", w.Code)
		}
	})

	t.Run("practice not found", func(t *testing.T) {
		req, err := newAuthenticatedRequest("POST", "/practice/999999/submit", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected status code 404, got %d", w.Code)
		}
	})
}

func TestHandlerGetPracticePage(t *testing.T) {
	const practiceID = 10003
	question, err := questionServer.DB.GetQuestion(1)
	if err != nil {
		t.Fatalf("Failed to get question: %v", err)
	}

	practice := (&model.Practice{}).NewPractice()
	practice.ID = practiceID
	practice.Questions = []int{question.ID}
	practice.TotalQuestions = questionServer.DB.GetTotalCount()
	practice.Duration = 5 * time.Minute
	practice.StartTime = time.Now()
	if err := questionServer.PM.Create(practice); err != nil {
		t.Fatalf("Failed to create practice: %v", err)
	}
	t.Cleanup(func() {
		if err := questionServer.PM.Delete(practice); err != nil {
			t.Errorf("Failed to delete practice: %v", err)
		}
	})

	req, err := newAuthenticatedRequest("GET", "/practice/10003", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d: %s", w.Code, w.Body.String())
	}

	body := w.Body.String()
	if !strings.Contains(body, question.Question) {
		t.Fatalf("Response does not contain question content")
	}
	if !strings.Contains(body, `id="practice_id" value="10003"`) {
		t.Fatalf("Response does not contain practice ID")
	}
	if !strings.Contains(body, `id="duration" value="300"`) &&
		!strings.Contains(body, `id="duration" value="299"`) {
		t.Fatalf("Response does not contain duration in seconds")
	}

	pagePath := "test_practice_page.html"
	if err := os.WriteFile(pagePath, w.Body.Bytes(), 0644); err != nil {
		t.Fatalf("Failed to save practice page: %v", err)
	}
	t.Logf("Practice page saved: %s", pagePath)
}

func TestHandlerGetPracticeResultPage(t *testing.T) {
	const practiceID = 10004
	question, err := questionServer.DB.GetQuestion(1)
	if err != nil {
		t.Fatalf("Failed to get question: %v", err)
	}

	practice := (&model.Practice{}).NewPractice()
	practice.ID = practiceID
	practice.Questions = []int{question.ID}
	practice.Answers[question.ID] = question.Answer
	practice.Completed = true
	if err := questionServer.PM.Create(practice); err != nil {
		t.Fatalf("Failed to create practice: %v", err)
	}
	t.Cleanup(func() {
		if err := questionServer.PM.Delete(practice); err != nil {
			t.Errorf("Failed to delete practice: %v", err)
		}
	})

	req, err := newAuthenticatedRequest("GET", "/practice/10004/result", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d: %s", w.Code, w.Body.String())
	}

	body := w.Body.String()
	for _, expected := range []string{
		"练习结果",
		question.Question,
		"回答正确",
		question.Explanation,
	} {
		if !strings.Contains(body, expected) {
			t.Fatalf("Response does not contain %q", expected)
		}
	}

	t.Run("practice not completed", func(t *testing.T) {
		practice.Completed = false
		defer func() { practice.Completed = true }()

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status code 400, got %d", w.Code)
		}
	})
}

func TestHandlerGetRandomQuestionPage(t *testing.T) {
	const sessionID = 10005
	questionIDs := questionServer.DB.GetALLQuestionIDs()
	if len(questionIDs) < 2 {
		t.Fatal("Expected at least two questions")
	}

	session := &model.RandomSession{
		ID:        sessionID,
		Questions: questionIDs[:2],
	}
	questionServer.RS.Create(session)

	t.Cleanup(func() { questionServer.RS.Delete(sessionID) })

	req, err := newAuthenticatedRequest("GET", "/question/random/10005", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d: %s", w.Code, w.Body.String())
	}
	if session.CurrentIndex != 0 {
		t.Fatalf("Expected initial index 0, got %d", session.CurrentIndex)
	}

	req, err = newAuthenticatedRequest("GET", "/question/random/10005?direction=next", nil)
	if err != nil {
		t.Fatalf("Failed to create next request: %v", err)
	}
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d: %s", w.Code, w.Body.String())
	}
	if session.CurrentIndex != 1 {
		t.Fatalf("Expected next index 1, got %d", session.CurrentIndex)
	}
	if !strings.Contains(w.Body.String(), `name="direction" value="last"`) {
		t.Fatal("Response does not contain random-session previous navigation")
	}
}

func TestHandlerPostLogin(t *testing.T) {
	username := fmt.Sprintf("login_test_%d", time.Now().UnixNano())
	password := "test-password"
	user, err := questionServer.UM.RegisterUser(model.RegisterRequest{
		Username: username,
		Password: password,
		Nickname: "Login Test",
	})
	if err != nil {
		t.Fatalf("Failed to create login test user: %v", err)
	}
	t.Cleanup(func() {
		if err := questionServer.UM.DeleteUser(user.Username); err != nil {
			t.Errorf("Failed to delete login test user: %v", err)
		}
	})

	payload, err := json.Marshal(model.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Failed to marshal login request: %v", err)
	}
	req, err := http.NewRequest("POST", "/auth/login", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("Failed to create login request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("Expected status code 302, got %d: %s", w.Code, w.Body.String())
	}
	if location := w.Header().Get("Location"); location != "/" {
		t.Fatalf("Expected redirect to /, got %q", location)
	}
	if cookie := w.Header().Get("Set-Cookie"); !strings.Contains(cookie, "access_token=") {
		t.Fatalf("Expected access_token cookie, got %q", cookie)
	}
}

func TestHandlerPostRegisterSQLite(t *testing.T) {
	username := fmt.Sprintf("register_test_%d", time.Now().UnixNano())
	t.Cleanup(func() {
		if err := questionServer.UM.DeleteUser(username); err != nil {
			t.Errorf("Failed to delete registration test user: %v", err)
		}
	})

	payload, err := json.Marshal(model.RegisterRequest{
		Username: username,
		Password: "Register123!",
		Nickname: "Register Test",
	})
	if err != nil {
		t.Fatalf("Failed to marshal register request: %v", err)
	}
	req, err := http.NewRequest("POST", "/auth/register", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("Failed to create register request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d: %s", w.Code, w.Body.String())
	}
	if _, exists := questionServer.UM.GetUser(username); !exists {
		t.Fatal("Expected registered user to be persisted in SQLite")
	}
	if cookie := w.Header().Get("Set-Cookie"); !strings.Contains(cookie, "access_token=") {
		t.Fatalf("Expected access_token cookie, got %q", cookie)
	}
}
