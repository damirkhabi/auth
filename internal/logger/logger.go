package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger
var slogger *zap.SugaredLogger

// Init initialize logger instance.
func Init(level string, format string) error {
	var config zap.Config
	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{"stdout"}
	config.DisableStacktrace = true

	loglevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return err
	}
	config.Level = zap.NewAtomicLevelAt(loglevel)

	zapLogger, err := config.Build()
	if err != nil {
		return err
	}
	globalLogger = zapLogger
	slogger = zapLogger.WithOptions(zap.AddCallerSkip(1)).Sugar()
	return nil
}

// GetLogger returns the logger instance.
// This instance is the entry point for all logging
func GetLogger() *zap.SugaredLogger {
	return slogger
}

// SetLogger sets the logger instance
// This is useful in testing as the logger can be overridden
// with a test logger
func SetLogger(l *zap.Logger) {
	slogger = l.Sugar()
}

// Reset reset overridden logger
func Reset() {
	slogger = globalLogger.WithOptions(zap.AddCallerSkip(1)).Sugar()
}

// Debug logs the provided arguments at DebugLevel.
// Spaces are added between arguments when neither is a string.
func Debug(args ...any) {
	slogger.Debug(args...)
}

// Debugf formats the message according to the format specifier
// and logs it at DebugLevel.
func Debugf(template string, args ...any) {
	slogger.Debugf(template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//
//	s.With(keysAndValues).Debug(msg)
func Debugw(msg string, keysAndValues ...any) {
	slogger.Debugw(msg, keysAndValues...)
}

// Info logs the provided arguments at InfoLevel.
// Spaces are added between arguments when neither is a string.
func Info(args ...any) {
	slogger.Info(args...)
}

// Infof formats the message according to the format specifier
// and logs it at InfoLevel.
func Infof(template string, args ...any) {
	slogger.Infof(template, args...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...any) {
	slogger.Infow(msg, keysAndValues...)
}

// Warn logs the provided arguments at WarnLevel.
// Spaces are added between arguments when neither is a string.
func Warn(args ...any) {
	slogger.Warn(args...)
}

// Warnf formats the message according to the format specifier
// and logs it at WarnLevel.
func Warnf(template string, args ...any) {
	slogger.Warnf(template, args...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...any) {
	slogger.Warnw(msg, keysAndValues...)
}

// Error logs the provided arguments at ErrorLevel.
// Spaces are added between arguments when neither is a string.
func Error(args ...any) {
	slogger.Error(args...)
}

// Errorf formats the message according to the format specifier
// and logs it at ErrorLevel.
func Errorf(template string, args ...any) {
	slogger.Errorf(template, args...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...any) {
	slogger.Errorw(msg, keysAndValues...)
}

// Fatal constructs a message with the provided arguments and calls os.Exit.
// Spaces are added between arguments when neither is a string.
func Fatal(args ...any) {
	slogger.Fatal(args...)
}

// Fatalf formats the message according to the format specifier
// and calls os.Exit.
func Fatalf(template string, args ...any) {
	slogger.Fatalf(template, args...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...any) {
	slogger.Fatalw(msg, keysAndValues...)
}

// WithOptions clones the current SugaredLogger, applies the supplied Options,
// and returns the result. It's safe to use concurrently.
func WithOptions(opts ...zap.Option) *zap.SugaredLogger {
	return slogger.WithOptions(opts...)
}

func init() {
	globalLogger, _ = zap.NewDevelopment()
	slogger = globalLogger.Sugar()
}
