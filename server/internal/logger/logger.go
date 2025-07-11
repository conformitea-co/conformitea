package logger

import (
	"fmt"
	"os"
	"strings"

	"conformitea/server/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var roLogger *zap.Logger

func Initialize() error {
	config := config.GetConfig().Logger

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

	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		return fmt.Errorf("invalid log level %q: %w", config.Level, err)
	}

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

	core := zapcore.NewCore(encoder, writer, level)

	// Build logger with options
	roLogger = zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1), // Skip one level for wrapper functions
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return nil
}

func GetLogger() *zap.Logger {
	if roLogger == nil {
		panic("logger not initialized, call Initialize first")
	}

	return roLogger
}
