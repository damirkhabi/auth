package converter

import (
	"database/sql"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/arifullov/auth/internal/model"
	desc "github.com/arifullov/auth/pkg/user_v1"
)

func ToUserFromService(user *model.User) *desc.GetResponse {
	role := desc.UserRole_USER
	if user.Role == model.AdminRole {
		role = desc.UserRole_ADMIN
	}
	return &desc.GetResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      role,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func ToUserCreateFromDesc(user *desc.CreateRequest) *model.CreateUser {
	role := model.UserRole
	if user.Role == desc.UserRole_ADMIN {
		role = model.AdminRole
	}
	return &model.CreateUser{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            role,
	}
}

func ToUserUpdateFromDesc(user *desc.UpdateRequest) *model.UpdateUser {
	userUpdate := &model.UpdateUser{
		ID: user.GetId(),
	}
	if user.GetName() != nil {
		userUpdate.Name = sql.NullString{String: user.GetName().Value, Valid: true}
	}
	if user.GetEmail() != nil {
		userUpdate.Email = sql.NullString{String: user.GetEmail().Value, Valid: true}
	}
	return userUpdate
}
