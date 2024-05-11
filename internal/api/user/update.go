package user

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/arifullov/auth/internal/converter"
	"github.com/arifullov/auth/internal/model"
	desc "github.com/arifullov/auth/pkg/user_v1"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, converter.ToUserUpdateFromDesc(req))
	if errors.Is(err, model.ErrUserNotFound) {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if errors.Is(err, model.ErrUserAlreadyExists) {
		return nil, status.Errorf(codes.InvalidArgument, "email address is already in use")
	}
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to update user")
	}
	log.Printf("update user: %v", req)
	return nil, nil
}
