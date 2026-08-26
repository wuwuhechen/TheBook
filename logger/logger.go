package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	App      *zap.Logger
	Business *zap.Logger
}

// InitLogger 初始化应用日志和错误日志。
// 普通日志写入 loggerPath，Error 及以上级别的日志额外写入同目录的 error.log。
func InitLogger(loggerPath string, development bool) (*Logger, func(), error) {
	if err := os.MkdirAll(filepath.Dir(loggerPath), 0755); err != nil {
		return nil, nil, err
	}

	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var encoder zapcore.Encoder

	if development {
		encoder = zapcore.NewConsoleEncoder(encodeConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encodeConfig)
	}

	errorPath := filepath.Join(filepath.Dir(loggerPath), "error.log")
	businessPath := filepath.Join(filepath.Dir(loggerPath), "business.log")
	fileRotator := newRotatingWriter(loggerPath)
	errorRotator := newRotatingWriter(errorPath)
	businessRotator := newRotatingWriter(businessPath)
	fileWriter := zapcore.AddSync(fileRotator)
	errorWriter := zapcore.AddSync(errorRotator)
	businessWriter := zapcore.AddSync(businessRotator)

	core := zapcore.NewTee(
		zapcore.NewCore(
			encoder,
			fileWriter,
			zap.InfoLevel,
		),
		zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			zap.InfoLevel,
		),
		zapcore.NewCore(
			encoder,
			errorWriter,
			zap.LevelEnablerFunc(func(level zapcore.Level) bool {
				return level >= zapcore.ErrorLevel
			}),
		),
	)

	log := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	businessLog := zap.New(
		zapcore.NewCore(encoder, businessWriter, zap.InfoLevel),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return &Logger{App: log, Business: businessLog}, func() {
		_ = log.Sync()
		_ = businessLog.Sync()
		_ = fileRotator.Close()
		_ = errorRotator.Close()
		_ = businessRotator.Close()
	}, nil
}

func newRotatingWriter(path string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    10,   // 单个文件最大 10 MB
		MaxBackups: 7,    // 最多保留 7 个旧文件
		MaxAge:     30,   // 最多保留 30 天
		Compress:   true, // 历史文件使用 gzip 压缩
	}
}
