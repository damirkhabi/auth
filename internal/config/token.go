package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	refreshTokenSecretKeyEnvName  = "REFRESH_TOKEN_SECRET_KEY"
	accessTokenSecretKeyEnvName   = "ACCESS_TOKEN_SECRET_KEY"
	refreshTokenExpirationEnvName = "REFRESH_TOKEN_EXPIRATION"
	accessTokenExpirationEnvName  = "ACCESS_TOKEN_EXPIRATION"
)

type TokenConfig interface {
	RefreshTokenSecretKey() string
	AccessTokenSecretKey() string
	RefreshTokenExpiration() time.Duration
	AccessTokenExpiration() time.Duration
}

type tokenConfig struct {
	refreshTokenSecretKey string
	accessTokenSecretKey  string

	refreshTokenExpiration time.Duration
	accessTokenExpiration  time.Duration
}

func NewTokenConfig() (TokenConfig, error) {
	refreshTokenSecretKey := os.Getenv(refreshTokenSecretKeyEnvName)
	if refreshTokenSecretKey == "" {
		return nil, errors.New("refresh token secret key not found")
	}

	accessTokenSecretKey := os.Getenv(accessTokenSecretKeyEnvName)
	if accessTokenSecretKey == "" {
		return nil, errors.New("access token secret key not found")
	}

	refreshTokenExpirationStr := os.Getenv(refreshTokenExpirationEnvName)
	if refreshTokenExpirationStr == "" {
		return nil, errors.New("refresh token expiration key not found")
	}
	refreshTokenExpiration, err := time.ParseDuration(refreshTokenExpirationStr)
	if err != nil {
		return nil, errors.New("invalid refresh token expiration")
	}

	accessTokenExpirationStr := os.Getenv(accessTokenExpirationEnvName)
	if accessTokenExpirationStr == "" {
		return nil, errors.New("access token expiration key not found")
	}
	accessTokenExpiration, err := time.ParseDuration(accessTokenExpirationStr)
	if err != nil {
		return nil, errors.New("invalid access token expiration")
	}

	return &tokenConfig{
		refreshTokenSecretKey:  refreshTokenSecretKey,
		accessTokenSecretKey:   accessTokenSecretKey,
		refreshTokenExpiration: refreshTokenExpiration,
		accessTokenExpiration:  accessTokenExpiration,
	}, nil
}

func (t *tokenConfig) RefreshTokenSecretKey() string {
	return t.refreshTokenSecretKey
}

func (t *tokenConfig) AccessTokenSecretKey() string {
	return t.accessTokenSecretKey
}

func (t *tokenConfig) RefreshTokenExpiration() time.Duration {
	return t.refreshTokenExpiration
}

func (t *tokenConfig) AccessTokenExpiration() time.Duration {
	return t.accessTokenExpiration
}
