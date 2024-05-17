package converter

import (
	"github.com/arifullov/auth/internal/model"
	modelRepo "github.com/arifullov/auth/internal/repository/user/model"
)

func ToUserFromRepo(user modelRepo.User) *model.User {
	return &model.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Role:         model.Role(user.Role),
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
