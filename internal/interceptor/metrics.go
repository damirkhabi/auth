package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/arifullov/auth/internal/metric"
)

func MetricsInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	metric.IncRequestCounter()

	timeStart := time.Now()

	res, err := handler(ctx, req)
	diffTime := time.Since(timeStart)

	if err != nil {
		metric.IncResponseCounter("error", info.FullMethod)
		metric.HistogramResponseTimeObserve("error", diffTime.Seconds())
	} else {
		metric.IncResponseCounter("success", info.FullMethod)
		metric.HistogramResponseTimeObserve("success", diffTime.Seconds())
	}

	return res, err
}
