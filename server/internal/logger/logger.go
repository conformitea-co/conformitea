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
	loggerConfigValues := config.GetConfig().Logger

	// Set defaults for empty values
	if loggerConfigValues.Level == "" {
		loggerConfigValues.Level = "info"
	}
	if loggerConfigValues.Format == "" {
		loggerConfigValues.Format = "json"
	}
	if loggerConfigValues.Output == "" {
		loggerConfigValues.Output = "stdout"
	}

	level, err := zapcore.ParseLevel(loggerConfigValues.Level)
	if err != nil {
		return fmt.Errorf("invalid log level %q: %w", loggerConfigValues.Level, err)
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
	switch strings.ToLower(loggerConfigValues.Format) {
	case "console":
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default: // json
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var writer zapcore.WriteSyncer
	switch strings.ToLower(loggerConfigValues.Output) {
	case "stderr":
		writer = zapcore.AddSync(os.Stderr)
	case "stdout":
		writer = zapcore.AddSync(os.Stdout)
	default:
		// Assume it's a file path
		file, err := os.OpenFile(loggerConfigValues.Output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file %q: %w", loggerConfigValues.Output, err)
		}
		writer = zapcore.AddSync(file)
	}

	core := zapcore.NewCore(encoder, writer, level)

	if config.BUILD == "development" {
		roLogger = zap.New(core,
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	} else {
		roLogger = zap.New(core)
	}

	return nil
}

func GetLogger() *zap.Logger {
	if roLogger == nil {
		panic("logger not initialized, call Initialize first")
	}

	return roLogger
}
