// Package shared 提供 Gin Handler 集成测试可复用的隔离环境。
package shared

import (
	"TheBook/auth"
	"TheBook/logger"
	"TheBook/model"
	"TheBook/model/bank"
	"TheBook/service"
	"TheBook/sqlite"
	"TheBook/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Fixture 持有一个完整但隔离的 Gin 服务实例。
type Fixture struct {
	Router *gin.Engine
	Server *service.Server
	token  string
}

// NewFixture 创建带独立数据库、路由和认证 Token 的 Handler 测试环境。
func NewFixture(t *testing.T) *Fixture {
	t.Helper()
	gin.SetMode(gin.TestMode)
	auth.InitJWT("gin-handler-test-secret")

	db, err := sqlite.OpenDatabase(filepath.Join(t.TempDir(), "handlers.db"))
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("access test database connection: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err := db.AutoMigrate(
		&sqlite.Question{}, &sqlite.User{}, &sqlite.QuestionProgress{},
		&bank.PracticeRecordRow{}, &bank.PracticeAnswer{}, &bank.WrongQuestionRow{},
	); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	questions := []sqlite.Question{
		{ID: 1, Category: "测试", Question: "第一题", OptionA: "A1", OptionB: "B1", OptionC: "C1", OptionD: "D1", Answer: 1, Explanation: "第一题解析"},
		{ID: 2, Category: "测试", Question: "第二题", OptionA: "A2", OptionB: "B2", OptionC: "C2", OptionD: "D2", Answer: 2, Explanation: "第二题解析"},
		{ID: 3, Category: "测试", Question: "第三题", OptionA: "A3", OptionB: "B3", OptionC: "C3", OptionD: "D3", Answer: 3, Explanation: "第三题解析"},
	}
	if err := db.Create(&questions).Error; err != nil {
		t.Fatalf("seed questions: %v", err)
	}

	server := &service.Server{
		DB: model.NewSQLiteQuestionBank(db), PM: model.NewPracticeBank(db),
		WQ: model.NewWrongQuestionBank(db), RS: model.NewRandomSessionBank(),
		UM: model.NewUserBankSQLite(db), QS: model.NewQuestionProgressBankSQLite(db),
		Log: &logger.Logger{App: zap.NewNop(), Business: zap.NewNop()},
	}
	root, err := utils.FindProjectRoot()
	if err != nil {
		t.Fatalf("find project root: %v", err)
	}
	router, err := service.GinInit(filepath.Join(root, "front_end", "*.html"), server, server.Log)
	if err != nil {
		t.Fatalf("initialize router: %v", err)
	}
	server.Router = router

	token, err := auth.GenerateToken(1, "test_user", "user", "Test User", time.Hour)
	if err != nil {
		t.Fatalf("create auth token: %v", err)
	}
	return &Fixture{Router: router, Server: server, token: token}
}

// Request 创建使用默认测试用户认证的 HTTP 请求。
func (f *Fixture) Request(method, target string, body []byte) *http.Request {
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+f.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

// RequestAs 创建使用指定测试用户认证的 HTTP 请求。
func (f *Fixture) RequestAs(t *testing.T, userID uint, method, target string, body []byte) *http.Request {
	t.Helper()
	token, err := auth.GenerateToken(userID, "user", "user", "User", time.Hour)
	if err != nil {
		t.Fatalf("create auth token: %v", err)
	}
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

// JSONBody 将测试请求值编码为 JSON 字节。
func JSONBody(t *testing.T, value any) []byte {
	t.Helper()
	body, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal JSON: %v", err)
	}
	return body
}

// Serve 使用测试路由执行请求并返回响应记录器。
func Serve(f *Fixture, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	f.Router.ServeHTTP(w, req)
	return w
}

// NewPractice 创建并注册指定题目的测试套题。
func NewPractice(t *testing.T, f *Fixture, id int, userID uint, questionIDs ...int) *model.Practice {
	t.Helper()
	p := (&model.Practice{}).NewPractice()
	p.ID, p.UserID = id, userID
	p.Questions = append([]int(nil), questionIDs...)
	p.TotalQuestions, p.PracticeSize = len(questionIDs), len(questionIDs)
	p.StartTime = time.Now()
	if err := f.Server.PM.Create(p); err != nil {
		t.Fatalf("create practice: %v", err)
	}
	t.Cleanup(func() { _ = f.Server.PM.Delete(p) })
	return p
}
