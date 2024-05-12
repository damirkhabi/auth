package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/arifullov/auth/internal/api/user"
	"github.com/arifullov/auth/internal/model"
	"github.com/arifullov/auth/internal/service"
	serviceMocks "github.com/arifullov/auth/internal/service/mocks"
	desc "github.com/arifullov/auth/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, false, false, 8)

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            desc.UserRole_ADMIN,
		}

		userObj = &model.CreateUser{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            model.AdminRole,
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success create",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, userObj).Return(id, nil)
				return mock
			},
		},
		{
			name: "failed create",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  status.Errorf(codes.Unavailable, "failed to create user"),
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, userObj).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceMock := tt.userServiceMock(mc)
			api := user.NewImplementation(userServiceMock)

			newID, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}

}
