package interceptor

import (
	"context"
	"errors"

	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/arifullov/auth/internal/sys"
	"github.com/arifullov/auth/internal/sys/codes"
	"github.com/arifullov/auth/internal/sys/validate"
)

type CircuitBreakerInterceptor struct {
	circuitBreaker *gobreaker.CircuitBreaker[any]
}

func NewCircuitBreakerInterceptor(cb *gobreaker.CircuitBreaker[any]) *CircuitBreakerInterceptor {
	return &CircuitBreakerInterceptor{
		circuitBreaker: cb,
	}
}

func (i *CircuitBreakerInterceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	res, err := handler(ctx, req)
	if validate.IsValidationError(err) {
		return nil, err
	}
	if sys.IsCommonError(err) {
		commErr := sys.GetCommonError(err)
		switch commErr.Code() {
		// Skip this errors
		case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists,
			codes.PermissionDenied, codes.Unauthenticated:
			res, _ = i.circuitBreaker.Execute(func() (any, error) {
				return res, nil
			})
			return nil, err
		default:
		}
	}

	res, err = i.circuitBreaker.Execute(func() (any, error) {
		return res, err
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return nil, status.Error(grpcCodes.Unavailable, "service unavailable")
		}
		return nil, err
	}

	return res, nil
}
