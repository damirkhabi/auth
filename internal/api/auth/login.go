package auth

import (
	"context"

	desc "github.com/arifullov/auth/pkg/auth_v1"
)

func (i *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	refreshToken, err := i.authService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &desc.LoginResponse{RefreshToken: refreshToken}, nil
}
