package interceptor

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
)

func InterceptorLogger(l *zap.SugaredLogger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			l.Debugw(msg, fields...)
		case logging.LevelInfo:
			l.Infow(msg, fields...)
		case logging.LevelWarn:
			l.Warnw(msg, fields...)
		case logging.LevelError:
			l.Errorw(msg, fields...)
		default:
			l.Errorw(msg, fields...)
		}
	})
}
