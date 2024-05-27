package user

import (
	"context"
	"net/mail"
	"time"

	"github.com/arifullov/auth/internal/model"
	"github.com/arifullov/auth/internal/sys/validate"
	"github.com/arifullov/auth/internal/utils"
)

func (s *serv) Create(ctx context.Context, user *model.CreateUser) (int64, error) {
	err := validate.Validate(
		ctx,
		emailIsValid(user.Email),
		passwordIsEqual(user.Password, user.PasswordConfirm),
	)
	if err != nil {
		return 0, err
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

func emailIsValid(email string) validate.Condition {
	return func(ctx context.Context) error {
		if _, err := mail.ParseAddress(email); err != nil {
			return validate.NewValidationErrors("invalid email")
		}
		return nil
	}
}

func passwordIsEqual(password string, confirmPassword string) validate.Condition {
	return func(ctx context.Context) error {
		if password != confirmPassword {
			return validate.NewValidationErrors("password mismatch")
		}
		return nil
	}
}
