package auth

import (
	"github.com/arifullov/auth/internal/client/db"
	"github.com/arifullov/auth/internal/config"
	"github.com/arifullov/auth/internal/repository"
	"github.com/arifullov/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
	tokenConfig    config.TokenConfig
}

func NewAuthService(
	userRepository repository.UserRepository,
	txManager db.TxManager,
	tokenConfig config.TokenConfig,
) service.AuthService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
		tokenConfig:    tokenConfig,
	}
}
