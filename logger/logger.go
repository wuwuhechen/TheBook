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

	// fileWriter, closeFile, err := zap.Open(loggerPath)
	// if err != nil {
	// 	return nil, nil, err
	// }

	fileWriter := newRotatingWriter(loggerPath)

	errorPath := filepath.Join(filepath.Dir(loggerPath), "error.log")
	errorWriter := newRotatingWriter(errorPath)

	// errorWriter, closeErrorFile, err := zap.Open(errorPath)
	// if err != nil {
	// 	// closeFile()
	// 	return nil, nil, err
	// }

	businessPath := filepath.Join(filepath.Dir(loggerPath), "business.log")
	businessWriter := newRotatingWriter(businessPath)

	// businessWriter, closeBusinessFile, err := zap.Open(businessPath)
	// if err != nil {
	// 	// closeFile()
	// 	// closeErrorFile()
	// 	return nil, nil, err
	// }

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
		// closeFile()
		// closeErrorFile()
		// closeBusinessFile()
	}, nil
}

func newRotatingWriter(path string) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    10,   // 单个文件最大 10 MB
		MaxBackups: 7,    // 最多保留 7 个旧文件
		MaxAge:     30,   // 最多保留 30 天
		Compress:   true, // 历史文件使用 gzip 压缩
	})
}
