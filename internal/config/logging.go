package config

import (
	"os"

	"github.com/pkg/errors"
)

const (
	logFormatEnvName = "LOG_FORMAT"
	logLevelEnvName  = "LOG_LEVEL"
)

type LoggerConfig interface {
	Format() string
	Level() string
}

type loggingConfig struct {
	format string
	level  string
}

func NewLoggingConfig() (LoggerConfig, error) {
	format := os.Getenv(logFormatEnvName)
	if format == "" {
		return nil, errors.New("log format not found")
	}

	level := os.Getenv(logLevelEnvName)
	if level == "" {
		return nil, errors.New("log level not found")
	}
	return &loggingConfig{
		format: format,
		level:  level,
	}, nil
}

func (l *loggingConfig) Format() string {
	return l.format
}

func (l *loggingConfig) Level() string {
	return l.level
}
