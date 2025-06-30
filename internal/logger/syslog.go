package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Global logger instance
	Log *zap.Logger
	// SugaredLogger for convenience methods
	Sugar *zap.SugaredLogger
)

// Config represents logger configuration
type Config struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// Initialize sets up the global logger based on configuration
func Initialize(v *viper.Viper) error {
	// Get logger configuration
	var config Config
	if err := v.UnmarshalKey("logger", &config); err != nil {
		// Use defaults if config is missing
		config = Config{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		}
	}

	// Set defaults for empty values
	if config.Level == "" {
		config.Level = "info"
	}
	if config.Format == "" {
		config.Format = "json"
	}
	if config.Output == "" {
		config.Output = "stdout"
	}

	// Parse log level
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		return fmt.Errorf("invalid log level %q: %w", config.Level, err)
	}

	// Create encoder config
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
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create encoder based on format
	var encoder zapcore.Encoder
	switch strings.ToLower(config.Format) {
	case "console":
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default: // json
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Create output writer
	var writer zapcore.WriteSyncer
	switch strings.ToLower(config.Output) {
	case "stderr":
		writer = zapcore.AddSync(os.Stderr)
	case "stdout":
		writer = zapcore.AddSync(os.Stdout)
	default:
		// Assume it's a file path
		file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file %q: %w", config.Output, err)
		}
		writer = zapcore.AddSync(file)
	}

	// Create core
	core := zapcore.NewCore(encoder, writer, level)

	// Build logger with options
	Log = zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1), // Skip one level for wrapper functions
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	// Create sugared logger
	Sugar = Log.Sugar()

	// Replace global logger
	zap.ReplaceGlobals(Log)

	return nil
}

// WithContext creates a logger with context fields
func WithContext(fields ...zap.Field) *zap.Logger {
	if Log == nil {
		// Initialize with defaults if not already initialized
		Initialize(viper.New())
	}
	return Log.With(fields...)
}

// WithRequestID creates a logger with request ID
func WithRequestID(requestID string) *zap.Logger {
	return WithContext(zap.String("request_id", requestID))
}

// WithUserID creates a logger with user ID
func WithUserID(userID string) *zap.Logger {
	return WithContext(zap.String("user_id", userID))
}

// WithError creates a logger with error field
func WithError(err error) *zap.Logger {
	return WithContext(zap.Error(err))
}

// Fields is a convenience function to create multiple fields
func Fields(fields map[string]interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

// Sync flushes any buffered log entries
func Sync() error {
	if Log != nil {
		return Log.Sync()
	}
	return nil
}

// Helper functions for common logging patterns

// Debug logs a debug message with optional fields
func Debug(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Debug(msg, fields...)
	}
}

// Info logs an info message with optional fields
func Info(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Info(msg, fields...)
	}
}

// Warn logs a warning message with optional fields
func Warn(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Warn(msg, fields...)
	}
}

// Error logs an error message with optional fields
func Error(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Error(msg, fields...)
	}
}

// Fatal logs a fatal message and exits the program
func Fatal(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Fatal(msg, fields...)
	}
	// If logger is not initialized, panic
	panic(msg)
}

// HTTPRequest creates fields for HTTP request logging
func HTTPRequest(method, path string, status int, duration int64) []zap.Field {
	return []zap.Field{
		zap.String("http.method", method),
		zap.String("http.path", path),
		zap.Int("http.status", status),
		zap.Int64("http.duration_ms", duration),
	}
}

// Performance creates fields for performance logging
func Performance(dbQueries int, dbDuration, redisDuration, apiDuration int64) []zap.Field {
	return []zap.Field{
		zap.Int("performance.db_queries", dbQueries),
		zap.Int64("performance.db_duration_ms", dbDuration),
		zap.Int64("performance.redis_duration_ms", redisDuration),
		zap.Int64("performance.external_api_duration_ms", apiDuration),
	}
}
