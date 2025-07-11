package sections

import (
	"errors"
	"strings"
)

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

func (l *LoggerConfig) Validate() error {
	var errs []string
	if l.Level == "" {
		errs = append(errs, "logger.level is required")
	}

	if l.Format == "" {
		errs = append(errs, "logger.format is required")
	}

	if l.Output == "" {
		errs = append(errs, "logger.output is required")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}
