package user

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/arifullov/auth/internal/model"
	desc "github.com/arifullov/auth/pkg/user_v1"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.GetId())
	if errors.Is(err, model.ErrUserNotFound) {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to delete user")
	}
	log.Printf("delete user with id: %d", req.GetId())
	return nil, nil
}
