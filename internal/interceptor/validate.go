package interceptor

import (
	"context"

	"google.golang.org/grpc"

	"github.com/arifullov/auth/internal/sys/validate"
)

type validator interface {
	Validate() error
}

func ValidateInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if val, ok := req.(validator); ok {
		if err := val.Validate(); err != nil {
			return nil, validate.NewValidationErrors(err.Error())
		}
	}

	return handler(ctx, req)
}
