package logger

import (
	"fmt"
	"os"
	"strings"

	"conformitea/infrastructure/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Initialize(loggerConfigValues config.LoggerConfig) (*zap.Logger, error) {
	if err := loggerConfigValues.Validate(); err != nil {
		return nil, fmt.Errorf("invalid logger configuration: %w", err)
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
			return nil, fmt.Errorf("failed to open log file %q: %w", loggerConfigValues.Output, err)
		}
		writer = zapcore.AddSync(file)
	}

	level, _ := zapcore.ParseLevel(loggerConfigValues.Level)
	core := zapcore.NewCore(encoder, writer, level)

	if config.BUILD == "development" {
		return zap.New(core,
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		), nil
	} else {
		return zap.New(core), nil
	}
}
