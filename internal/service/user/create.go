package user

import (
	"context"
	"net/mail"
	"time"

	"github.com/arifullov/auth/internal/model"
	"github.com/arifullov/auth/internal/utils"
)

func (s *serv) Create(ctx context.Context, user *model.CreateUser) (int64, error) {
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return 0, model.ErrInvalidEmail
	}
	if user.Password != user.PasswordConfirm {
		return 0, model.ErrPasswordMismatch
	}
	passwordHash := utils.MakePbkdf2SHA256(user.Password)
	now := time.Now()
	user.PasswordHash = passwordHash
	user.CreatedAt = now
	user.UpdatedAt = now
	id, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return 0, err
	}
	return id, nil
}
