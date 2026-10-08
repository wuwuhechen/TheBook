package logger

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewRotatingWriterConfiguration(t *testing.T) {
	writer := newRotatingWriter("logs/app.log")

	if writer.Filename != "logs/app.log" {
		t.Errorf("Filename = %q, want %q", writer.Filename, "logs/app.log")
	}
	if writer.MaxSize != 10 || writer.MaxBackups != 7 || writer.MaxAge != 30 || !writer.Compress {
		t.Errorf("rotating writer = %+v, want MaxSize=10, MaxBackups=7, MaxAge=30, Compress=true", writer)
	}
}

func TestInitLoggerWritesExpectedLogFiles(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "logs", "app.log")
	log, closeLog, err := InitLogger(logPath, false)
	if err != nil {
		t.Fatalf("InitLogger() error = %v", err)
	}
	if log == nil || log.App == nil || log.Business == nil {
		t.Fatalf("InitLogger() = %+v, want initialized App and Business loggers", log)
	}

	log.App.Info("application event")
	log.App.Error("application error")
	log.Business.Info("business event")
	closeLog()

	for _, path := range []string{
		logPath,
		filepath.Join(filepath.Dir(logPath), "error.log"),
		filepath.Join(filepath.Dir(logPath), "business.log"),
	} {
		info, statErr := os.Stat(path)
		if statErr != nil {
			t.Errorf("expected log file %q: %v", path, statErr)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("log file %q is empty", path)
		}
	}
}

func TestRequestLoggerRecordsStatusLevelAndFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		status    int
		wantLevel zapcoreLevel
		wantMsg   string
	}{
		{name: "success", status: http.StatusOK, wantLevel: zapInfo, wantMsg: "http request processed"},
		{name: "client error", status: http.StatusBadRequest, wantLevel: zapWarn, wantMsg: "http request rejected"},
		{name: "server error", status: http.StatusInternalServerError, wantLevel: zapError, wantMsg: "http request failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, observed := observer.New(zap.DebugLevel)
			log := &Logger{App: zap.New(core)}
			router := gin.New()
			router.Use(RequestLogger(log))
			router.GET("/health", func(c *gin.Context) {
				c.Status(tt.status)
			})

			request := httptest.NewRequest(http.MethodGet, "/health", nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			entries := observed.All()
			if len(entries) != 1 {
				t.Fatalf("observed log entries = %d, want 1", len(entries))
			}

			entry := entries[0]
			if entry.Level.String() != tt.wantLevel.String() {
				t.Errorf("log level = %s, want %s", entry.Level, tt.wantLevel)
			}
			if entry.Message != tt.wantMsg {
				t.Errorf("log message = %q, want %q", entry.Message, tt.wantMsg)
			}

			fields := entry.ContextMap()
			if fields["source"] != "backend" || fields["method"] != http.MethodGet || fields["path"] != "/health" {
				t.Errorf("log fields = %#v, want source, method, and path", fields)
			}
			if fields["status"] != int64(tt.status) {
				t.Errorf("status field = %#v, want %d", fields["status"], tt.status)
			}
			if _, ok := fields["duration"]; !ok {
				t.Errorf("log fields = %#v, want duration", fields)
			}
		})
	}
}

// zapcoreLevel keeps the table concise while retaining zap's concrete levels.
type zapcoreLevel = zapcore.Level

const (
	zapInfo  = zapcore.InfoLevel
	zapWarn  = zapcore.WarnLevel
	zapError = zapcore.ErrorLevel
)
