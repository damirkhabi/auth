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
)

type GRPCConfig interface {
	Address() string
	RateLimit() int
	PeriodRateLimit() time.Duration
}

type grpcConfig struct {
	host string
	port string

	rateLimit       int
	periodRateLimit time.Duration
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

	return &grpcConfig{
		host:            host,
		port:            port,
		rateLimit:       rateLimit,
		periodRateLimit: period,
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
