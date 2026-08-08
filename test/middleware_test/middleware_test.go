package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"TheBook/logger"

	"github.com/gin-gonic/gin"
)

func TestRequestLoggerWritesErrorLog(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 使用临时目录，避免污染项目中的 logs 文件夹。
	logDir := t.TempDir()
	appLogPath := filepath.Join(logDir, "app.log")
	errorLogPath := filepath.Join(logDir, "error.log")

	log, closeLogger, err := logger.InitLogger(appLogPath, false)
	if err != nil {
		t.Fatalf("初始化日志器失败: %v", err)
	}
	defer closeLogger()

	router := gin.New()
	router.Use(logger.RequestLogger(log))

	// 模拟一个后端异常接口。
	router.GET("/test-error", func(c *gin.Context) {
		_ = c.Error(errors.New("test internal error"))
		c.Status(http.StatusInternalServerError)
	})

	request := httptest.NewRequest(http.MethodGet, "/test-error", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("期望状态码为 500，实际为 %d", response.Code)
	}

	// 先刷新日志缓冲，再读取文件。
	_ = log.Sync()

	data, err := os.ReadFile(errorLogPath)
	if err != nil {
		t.Fatalf("读取 error.log 失败: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "http request failed") {
		t.Fatalf("error.log 未记录请求失败日志，实际内容：%s", content)
	}

	if !strings.Contains(content, "test internal error") {
		t.Fatalf("error.log 未记录错误详情，实际内容：%s", content)
	}
}
