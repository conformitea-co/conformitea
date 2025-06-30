package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name   string
		config map[string]interface{}
		want   Config
	}{
		{
			name: "default config",
			config: map[string]interface{}{
				"logger": map[string]interface{}{},
			},
			want: Config{
				Level:  "info",
				Format: "json",
				Output: "stdout",
			},
		},
		{
			name: "custom config",
			config: map[string]interface{}{
				"logger": map[string]interface{}{
					"level":  "debug",
					"format": "console",
					"output": "stderr",
				},
			},
			want: Config{
				Level:  "debug",
				Format: "console",
				Output: "stderr",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			for k, val := range tt.config {
				v.Set(k, val)
			}

			err := Initialize(v)
			if err != nil {
				t.Fatalf("Initialize() error = %v", err)
			}
			if Log == nil {
				t.Error("Initialize() Log is nil")
			}
			if Sugar == nil {
				t.Error("Initialize() Sugar is nil")
			}
		})
	}
}

func TestInitializeInvalidLevel(t *testing.T) {
	v := viper.New()
	v.Set("logger.level", "invalid")

	err := Initialize(v)
	if err == nil {
		t.Error("Initialize() expected error for invalid level")
	}
	if !strings.Contains(err.Error(), "invalid log level") {
		t.Errorf("Initialize() error = %v, want error containing 'invalid log level'", err)
	}
}

func TestLoggingFunctions(t *testing.T) {
	// Create an observer to capture logs
	core, recorded := observer.New(zapcore.InfoLevel)
	Log = zap.New(core)
	Sugar = Log.Sugar()

	// Test Info
	Info("test info message", zap.String("key", "value"))
	if recorded.Len() != 1 {
		t.Errorf("Info() recorded %d logs, want 1", recorded.Len())
	}
	if recorded.All()[0].Message != "test info message" {
		t.Errorf("Info() message = %v, want 'test info message'", recorded.All()[0].Message)
	}
	if recorded.All()[0].Level != zapcore.InfoLevel {
		t.Errorf("Info() level = %v, want InfoLevel", recorded.All()[0].Level)
	}

	// Test Warn
	Warn("test warn message")
	if recorded.Len() != 2 {
		t.Errorf("Warn() recorded %d logs, want 2", recorded.Len())
	}
	if recorded.All()[1].Message != "test warn message" {
		t.Errorf("Warn() message = %v, want 'test warn message'", recorded.All()[1].Message)
	}
	if recorded.All()[1].Level != zapcore.WarnLevel {
		t.Errorf("Warn() level = %v, want WarnLevel", recorded.All()[1].Level)
	}

	// Test Error
	Error("test error message")
	if recorded.Len() != 3 {
		t.Errorf("Error() recorded %d logs, want 3", recorded.Len())
	}
	if recorded.All()[2].Message != "test error message" {
		t.Errorf("Error() message = %v, want 'test error message'", recorded.All()[2].Message)
	}
	if recorded.All()[2].Level != zapcore.ErrorLevel {
		t.Errorf("Error() level = %v, want ErrorLevel", recorded.All()[2].Level)
	}
}

func TestWithContext(t *testing.T) {
	// Create an observer to capture logs
	core, recorded := observer.New(zapcore.InfoLevel)
	Log = zap.New(core)

	// Test WithRequestID
	logger := WithRequestID("test-request-id")
	logger.Info("test message")

	if recorded.Len() != 1 {
		t.Errorf("WithRequestID() recorded %d logs, want 1", recorded.Len())
	}
	entry := recorded.All()[0]
	if entry.Message != "test message" {
		t.Errorf("WithRequestID() message = %v, want 'test message'", entry.Message)
	}

	// Check that request_id field exists
	requestIDField := entry.ContextMap()["request_id"]
	if requestIDField != "test-request-id" {
		t.Errorf("WithRequestID() request_id = %v, want 'test-request-id'", requestIDField)
	}
}

func TestHTTPRequest(t *testing.T) {
	fields := HTTPRequest("GET", "/api/test", 200, 45)

	if len(fields) != 4 {
		t.Errorf("HTTPRequest() returned %d fields, want 4", len(fields))
	}

	// Check field count and keys
	foundFields := make(map[string]bool)
	for _, field := range fields {
		foundFields[field.Key] = true
	}

	expectedKeys := []string{"http.method", "http.path", "http.status", "http.duration_ms"}
	for _, key := range expectedKeys {
		if !foundFields[key] {
			t.Errorf("HTTPRequest() missing field %s", key)
		}
	}
}

func TestJSONOutput(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a custom core that writes to our buffer
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	writer := zapcore.AddSync(&buf)
	core := zapcore.NewCore(encoder, writer, zapcore.InfoLevel)

	Log = zap.New(core)

	// Log a message
	Info("test json output", zap.String("field", "value"))

	// Parse the JSON output
	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if logEntry["message"] != "test json output" {
		t.Errorf("JSON message = %v, want 'test json output'", logEntry["message"])
	}
	if logEntry["level"] != "info" {
		t.Errorf("JSON level = %v, want 'info'", logEntry["level"])
	}
	if logEntry["field"] != "value" {
		t.Errorf("JSON field = %v, want 'value'", logEntry["field"])
	}
}

func TestConsoleOutput(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a custom core with console encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	writer := zapcore.AddSync(&buf)
	core := zapcore.NewCore(encoder, writer, zapcore.InfoLevel)

	Log = zap.New(core)

	// Log a message
	Info("test console output")

	// Check that output contains the message
	output := buf.String()
	if !strings.Contains(output, "test console output") {
		t.Errorf("Console output doesn't contain message, got: %v", output)
	}
	if !strings.Contains(strings.ToUpper(output), "INFO") {
		t.Errorf("Console output doesn't contain INFO level, got: %v", output)
	}
}
