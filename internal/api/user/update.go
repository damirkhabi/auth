package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/arifullov/auth/internal/converter"
	desc "github.com/arifullov/auth/pkg/user_v1"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, converter.ToUserUpdateFromDesc(req))
	if err != nil {
		return nil, err
	}
	return nil, nil
}
