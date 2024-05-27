package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/arifullov/auth/pkg/user_v1"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return nil, nil
}
