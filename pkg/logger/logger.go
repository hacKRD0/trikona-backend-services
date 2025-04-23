package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log is the global logger instance
	Log *zap.Logger
)

// InitLogger initializes the global logger
func InitLogger() error {
	// Configure encoder
	encoderConfig := zapcore.EncoderConfig{
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

	// Set log level based on environment
	var level zapcore.Level
	if os.Getenv("ENV") == "production" {
		level = zapcore.InfoLevel
	} else {
		level = zapcore.DebugLevel
	}

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	// Create logger
	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

// Info logs an info message with fields
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Error logs an error message with fields
func Error(msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	Log.Error(msg, fields...)
}

// Debug logs a debug message with fields
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Warn logs a warning message with fields
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Fatal logs a fatal message with fields and exits
func Fatal(msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	Log.Fatal(msg, fields...)
}

// WithRequestID creates a logger with request ID
func WithRequestID(requestID string) *zap.Logger {
	return Log.With(zap.String("request_id", requestID))
}

// WithDuration creates a logger with duration
func WithDuration(start time.Time) zap.Field {
	return zap.Duration("duration", time.Since(start))
} 