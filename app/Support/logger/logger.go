package logger

import (
	"os"
	"path/filepath"
	"sync"

	"mall/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Init 初始化日志实例
func Init() {
	once.Do(func() {
		log = newLogger()
	})
}

// newLogger 创建 zap logger, 支持日志轮转
func newLogger() *zap.Logger {
	cfg := config.Log()

	// 确保日志目录存在
	if err := os.MkdirAll(filepath.Dir(cfg.Path), 0755); err != nil {
		panic("无法创建日志目录: " + err.Error())
	}

	// 日志轮转配置 (lumberjack)
	writer := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Path, "mall.log"),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   true,
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 日志级别
	level := parseLevel(cfg.Level)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(writer), zapcore.AddSync(os.Stdout)),
		level,
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func parseLevel(s string) zapcore.Level {
	switch s {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}

// ---------- 便捷方法 ----------

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// Sync 刷新日志缓冲区
func Sync() {
	_ = log.Sync()
}
