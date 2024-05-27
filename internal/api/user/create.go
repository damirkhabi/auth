package user

import (
	"context"

	"github.com/arifullov/auth/internal/converter"
	desc "github.com/arifullov/auth/pkg/user_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.userService.Create(ctx, converter.ToUserCreateFromDesc(req))
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{Id: id}, nil
}
