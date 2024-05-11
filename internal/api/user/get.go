package user

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/arifullov/auth/internal/converter"
	"github.com/arifullov/auth/internal/model"
	desc "github.com/arifullov/auth/pkg/user_v1"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	userObj, err := i.userService.Get(ctx, req.GetId())
	if errors.Is(err, model.ErrUserNotFound) {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to get user")
	}
	log.Printf("get user: %v", userObj)
	return converter.ToUserFromService(userObj), nil
}
