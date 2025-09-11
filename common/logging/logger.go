package logging

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/tienhung-ho/smart-document/common/config"
)

type Logger struct {
	*zap.SugaredLogger
	config *config.LoggingConfig
}

var (
	globalLogger *Logger
)

func InitLogger(cfg *config.LoggingConfig) (*Logger, error) {
	logger, err := NewLogger(cfg)
	if err != nil {
		return nil, err
	}

	globalLogger = logger
	return logger, nil
}

func NewLogger(cfg *config.LoggingConfig) (*Logger, error) {
	level, err := parseLogLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	var core zapcore.Core

	switch strings.ToLower(cfg.Output) {
	case "stdout", "console":
		core = createConsoleCore(level, cfg.Format)
	case "file":
		core = createFileCore(level, cfg)
	case "both":
		consoleCore := createConsoleCore(level, cfg.Format)
		fileCore := createFileCore(level, cfg)
		core = zapcore.NewTee(consoleCore, fileCore)
	default:
		core = createConsoleCore(level, cfg.Format)
	}

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		SugaredLogger: zapLogger.Sugar(),
		config:        cfg,
	}, nil
}

func createConsoleCore(level zapcore.Level, format string) zapcore.Core {
	var encoder zapcore.Encoder

	if strings.ToLower(format) == "json" {
		encoder = zapcore.NewJSONEncoder(getEncoderConfig())
	} else {
		encoder = zapcore.NewConsoleEncoder(getConsoleEncoderConfig())
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
}

func createFileCore(level zapcore.Level, cfg *config.LoggingConfig) zapcore.Core {
	writer := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	encoder := zapcore.NewJSONEncoder(getEncoderConfig())
	return zapcore.NewCore(encoder, zapcore.AddSync(writer), level)
}

func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func getConsoleEncoderConfig() zapcore.EncoderConfig {
	config := getEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	return config
}

func parseLogLevel(level string) (zapcore.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unknown log level: %s", level)
	}
}

// Global logger functions
func Debug(args ...any) {
	if globalLogger != nil {
		globalLogger.Debug(args...)
	}
}

func Debugf(template string, args ...any) {
	if globalLogger != nil {
		globalLogger.Debugf(template, args...)
	}
}

func Info(args ...any) {
	if globalLogger != nil {
		globalLogger.Info(args...)
	}
}

func Infof(template string, args ...any) {
	if globalLogger != nil {
		globalLogger.Infof(template, args...)
	}
}

func Warn(args ...any) {
	if globalLogger != nil {
		globalLogger.Warn(args...)
	}
}

func Warnf(template string, args ...any) {
	if globalLogger != nil {
		globalLogger.Warnf(template, args...)
	}
}

func Error(args ...any) {
	if globalLogger != nil {
		globalLogger.Error(args...)
	}
}

func Errorf(template string, args ...any) {
	if globalLogger != nil {
		globalLogger.Errorf(template, args...)
	}
}

func Fatal(args ...any) {
	if globalLogger != nil {
		globalLogger.Fatal(args...)
	}
}

func Fatalf(template string, args ...any) {
	if globalLogger != nil {
		globalLogger.Fatalf(template, args...)
	}
}

func Panic(args ...any) {
	if globalLogger != nil {
		globalLogger.Panic(args...)
	}
}

func Panicf(template string, args ...any) {
	if globalLogger != nil {
		globalLogger.Panicf(template, args...)
	}
}

// Structured logging with fields
func WithFields(fields map[string]any) *zap.SugaredLogger {
	if globalLogger == nil {
		return nil
	}

	var zapFields []any
	for k, v := range fields {
		zapFields = append(zapFields, k, v)
	}

	return globalLogger.With(zapFields...)
}

// Context-aware logging
func WithContext(ctx map[string]any) *Logger {
	if globalLogger == nil {
		return nil
	}

	var zapFields []any
	for k, v := range ctx {
		zapFields = append(zapFields, k, v)
	}

	return &Logger{
		SugaredLogger: globalLogger.With(zapFields...),
		config:        globalLogger.config,
	}
}

// HTTP request logging
func LogHTTPRequest(method, path, userID string, statusCode int, duration float64) {
	if globalLogger != nil {
		globalLogger.Infow("HTTP Request",
			"method", method,
			"path", path,
			"user_id", userID,
			"status_code", statusCode,
			"duration_ms", duration,
		)
	}
}

// Database query logging
func LogDBQuery(query string, duration float64, err error) {
	if globalLogger == nil {
		return
	}

	if err != nil {
		globalLogger.Errorw("Database Query Failed",
			"query", query,
			"duration_ms", duration,
			"error", err.Error(),
		)
	} else {
		globalLogger.Debugw("Database Query",
			"query", query,
			"duration_ms", duration,
		)
	}
}

// Error with stack trace
func ErrorWithStack(err error, message string) {
	if globalLogger != nil {
		globalLogger.Errorw(message,
			"error", err.Error(),
			"stack", fmt.Sprintf("%+v", err),
		)
	}
}

// Sync flushes any buffered log entries
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// Close closes the logger
func Close() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}
