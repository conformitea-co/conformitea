package config

import (
	"errors"
	"slices"
	"strings"

	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

func (l *LoggerConfig) Validate() error {
	var errs []error

	if _, err := zapcore.ParseLevel(l.Level); err != nil {
		errs = append(errs, errors.New("logger.level is invalid"))
	}

	if !slices.Contains([]string{"json", "console"}, strings.ToLower(l.Format)) {
		errs = append(errs, errors.New("logger.format must be one of: json, console"))
	}

	if l.Output == "" {
		errs = append(errs, errors.New("logger.output is required"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
