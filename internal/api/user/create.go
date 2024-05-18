package user

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/arifullov/auth/internal/converter"
	"github.com/arifullov/auth/internal/logger"
	"github.com/arifullov/auth/internal/model"
	desc "github.com/arifullov/auth/pkg/user_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.userService.Create(ctx, converter.ToUserCreateFromDesc(req))
	if errors.Is(err, model.ErrInvalidEmail) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email")
	}
	if errors.Is(err, model.ErrPasswordMismatch) {
		return nil, status.Errorf(codes.InvalidArgument, "password mismatch")
	}
	if errors.Is(err, model.ErrUserAlreadyExists) {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}
	if err != nil {
		logger.Errorf("failed to insert user: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to create user")
	}
	logger.Infof("inserted user with id: %d", id)
	return &desc.CreateResponse{Id: id}, nil
}
