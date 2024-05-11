package user

import (
	"github.com/arifullov/auth/internal/repository"
	"github.com/arifullov/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
}

func NewUserService(
	userRepository repository.UserRepository,
) service.UserService {
	return &serv{
		userRepository: userRepository,
	}
}
