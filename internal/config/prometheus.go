package config

import (
	"net"
	"os"

	"github.com/pkg/errors"
)

const (
	prometheusHostEnvName = "PROMETHEUS_HOST"
	prometheusPortEnvName = "PROMETHEUS_PORT"
)

type PrometheusConfig interface {
	Address() string
}

type prometheusConfig struct {
	host string
	port string
}

func NewPrometheusConfig() (PrometheusConfig, error) {
	host := os.Getenv(prometheusHostEnvName)
	if host == "" {
		return nil, errors.New("prometheus host not found")
	}

	port := os.Getenv(prometheusPortEnvName)
	if port == "" {
		return nil, errors.New("prometheus port not found")
	}

	return &prometheusConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *prometheusConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
