package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/arifullov/auth/internal/client/db"
	txManagerMocks "github.com/arifullov/auth/internal/client/db/mocks"
	"github.com/arifullov/auth/internal/model"
	"github.com/arifullov/auth/internal/repository"
	repositoryMocks "github.com/arifullov/auth/internal/repository/mocks"
	"github.com/arifullov/auth/internal/service/user"
)

func TestGet(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager
	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Name()
		email     = gofakeit.Email()
		createdAt = gofakeit.Date()

		repositoryErr = fmt.Errorf("repository error")

		userObj = &model.User{
			ID:        id,
			Name:      name,
			Email:     email,
			Role:      model.AdminRole,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		}
	)

	tests := []struct {
		name               string
		args               args
		want               *model.User
		err                error
		userRepositoryMock userRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success get",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: userObj,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(userObj, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				return txManagerMocks.NewTxManagerMock(mc)
			},
		},
		{
			name: "failed get",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  repositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repositoryErr)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				return txManagerMocks.NewTxManagerMock(mc)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userRepositoryMock := tt.userRepositoryMock(mc)
			service := user.NewUserService(userRepositoryMock, tt.txManagerMock(mc))

			newUser, err := service.Get(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newUser)
		})
	}

}
