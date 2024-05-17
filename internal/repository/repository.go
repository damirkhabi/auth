package repository

import (
	"context"

	"github.com/arifullov/auth/internal/model"
)

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i UserRepository -o ./mocks/ -s "_minimock.go"
type UserRepository interface {
	Create(ctx context.Context, user *model.CreateUser) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.UpdateUser) error
	Delete(ctx context.Context, id int64) error
}

type AccessRepository interface {
	GetRouteRoles(ctx context.Context, route string) ([]model.Role, error)
}
