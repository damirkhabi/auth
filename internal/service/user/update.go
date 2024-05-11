package user

import (
	"context"

	"github.com/arifullov/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, user *model.UpdateUser) error {
	if err := s.userRepository.Update(ctx, user); err != nil {
		return err
	}
	return nil
}
