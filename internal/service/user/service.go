package user

import (
	"github.com/arifullov/auth/internal/client/db"
	"github.com/arifullov/auth/internal/repository"
	"github.com/arifullov/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

func NewUserService(
	userRepository repository.UserRepository,
	txManager db.TxManager,
) service.UserService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
