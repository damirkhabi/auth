package interceptor

import (
	"context"
	"errors"

	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	res, err := i.circuitBreaker.Execute(func() (any, error) {
		return handler(ctx, req)
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return nil, status.Error(codes.Unavailable, "service unavailable")
		}
		return nil, err
	}

	return res, nil
}
