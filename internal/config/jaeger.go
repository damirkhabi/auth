package config

import (
	"errors"
	"os"
)

const (
	jaegerCollectorEndpointEnvName     = "JAEGER_COLLECTOR_ENDPOINT"
	jaegerServiceNameEnvName           = "JAEGER_SERVICE_NAME"
	jaegerDeploymentEnvironmentEnvName = "JAEGER_DEPLOYMENT_ENVIRONMENT"
)

type JaegerConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	DeploymentEnvironment() string
}

type jaegerConfig struct {
	collectorEndpoint     string
	serviceName           string
	deploymentEnvironment string
}

func NewJaegerConfig() (JaegerConfig, error) {
	collectorEndpoint := os.Getenv(jaegerCollectorEndpointEnvName)
	if collectorEndpoint == "" {
		return nil, errors.New("jaeger collector endpoint not found")
	}
	serviceName := os.Getenv(jaegerServiceNameEnvName)
	if serviceName == "" {
		return nil, errors.New("jaeger service name not found")
	}
	deploymentEnvironment := os.Getenv(jaegerDeploymentEnvironmentEnvName)
	if deploymentEnvironment == "" {
		return nil, errors.New("jaeger deployment environment not found")
	}
	return &jaegerConfig{
		collectorEndpoint:     collectorEndpoint,
		serviceName:           serviceName,
		deploymentEnvironment: deploymentEnvironment,
	}, nil
}

func (j *jaegerConfig) CollectorEndpoint() string {
	return j.collectorEndpoint
}

func (j *jaegerConfig) ServiceName() string {
	return j.serviceName
}

func (j *jaegerConfig) DeploymentEnvironment() string {
	return j.deploymentEnvironment
}
