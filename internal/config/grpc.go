package config

import (
	"net"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"

	grpcRateLimitEnvName       = "GRPC_RATE_LIMIT"
	grpcPeriodRateLimitEnvName = "GRPC_PERIOD_RATE_LIMIT"

	grpcCircuitBreakerMaxRequestsEnvName = "GRPC_CIRCUIT_BREAKER_MAX_REQUESTS"
	grpcCircuitTimeoutEnvName            = "GRPC_CIRCUIT_BREAKER_TIMEOUT"
	grpcCircuitFailureRatioEnvName       = "GRPC_CIRCUIT_FAILURE_RATIO"
)

type GRPCConfig interface {
	Address() string
	RateLimit() int
	PeriodRateLimit() time.Duration
	CircuitBreakerMaxRequests() uint32
	CircuitBreakerTimeout() time.Duration
	FailureRatio() float64
}

type grpcConfig struct {
	host string
	port string

	rateLimit       int
	periodRateLimit time.Duration

	circuitBreakerMaxRequests uint32
	circuitBreakerTimeout     time.Duration
	failureRatio              float64
}

func NewGRPCConfig() (GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if host == "" {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if port == "" {
		return nil, errors.New("grpc port not found")
	}

	rateLimitStr := os.Getenv(grpcRateLimitEnvName)
	if rateLimitStr == "" {
		return nil, errors.New("rate limit not found")
	}
	rateLimit, err := strconv.Atoi(rateLimitStr)
	if err != nil {
		return nil, errors.New("rate limit error")
	}

	periodStr := os.Getenv(grpcPeriodRateLimitEnvName)
	if periodStr == "" {
		return nil, errors.New("period limit not found")
	}
	period, err := time.ParseDuration(periodStr)
	if err != nil {
		return nil, errors.New("invalid period")
	}

	circuitBreakerMaxRequestsStr := os.Getenv(grpcCircuitBreakerMaxRequestsEnvName)
	if circuitBreakerMaxRequestsStr == "" {
		return nil, errors.New("circuit breaker max requests not found")
	}
	circuitBreakerMaxRequests, err := strconv.Atoi(circuitBreakerMaxRequestsStr)
	if err != nil {
		return nil, errors.New("invalid circuit breaker max requests")
	}

	circuitTimeoutStr := os.Getenv(grpcCircuitTimeoutEnvName)
	if circuitTimeoutStr == "" {
		return nil, errors.New("circuit timeout not found")
	}
	circuitTimeout, err := time.ParseDuration(circuitTimeoutStr)
	if err != nil {
		return nil, errors.New("invalid circuit timeout")
	}

	circuitFailureRatioStr := os.Getenv(grpcCircuitFailureRatioEnvName)
	if circuitFailureRatioStr == "" {
		return nil, errors.New("circuit failure ratio not found")
	}
	circuitFailureRatio, err := strconv.ParseFloat(circuitFailureRatioStr, 64)
	if err != nil {
		return nil, errors.New("invalid circuit failure ratio")
	}

	return &grpcConfig{
		host:                      host,
		port:                      port,
		rateLimit:                 rateLimit,
		periodRateLimit:           period,
		circuitBreakerMaxRequests: uint32(circuitBreakerMaxRequests),
		circuitBreakerTimeout:     circuitTimeout,
		failureRatio:              circuitFailureRatio,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *grpcConfig) RateLimit() int {
	return cfg.rateLimit
}

func (cfg *grpcConfig) PeriodRateLimit() time.Duration {
	return cfg.periodRateLimit
}

func (cfg *grpcConfig) CircuitBreakerMaxRequests() uint32 {
	return cfg.circuitBreakerMaxRequests
}

func (cfg *grpcConfig) CircuitBreakerTimeout() time.Duration {
	return cfg.circuitBreakerTimeout
}

func (cfg *grpcConfig) FailureRatio() float64 {
	return cfg.failureRatio
}
