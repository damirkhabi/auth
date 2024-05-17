package auth

import (
	"context"
	"time"

	"github.com/arifullov/auth/internal/model"
	"github.com/arifullov/auth/internal/utils"
)

func (s *serv) Login(ctx context.Context, username string, password string) (string, error) {
	user, err := s.userRepository.GetByEmail(ctx, username)
	if err != nil {
		return "", model.ErrUserNotFound
	}
	isPasswordEqual, err := utils.CheckPbkdf2SHA256(password, user.PasswordHash)
	if err != nil {
		return "", err
	}
	if !isPasswordEqual {
		return "", model.ErrWrongCredentials
	}

	refreshToken, err := generateRefreshToken(user, utils.S2B(s.tokenConfig.RefreshTokenSecretKey()), s.tokenConfig.RefreshTokenExpiration())
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func generateRefreshToken(user *model.User, secretKey []byte, duration time.Duration) (string, error) {
	return utils.GenerateToken(user, secretKey, duration)
}

func generateAccessToken(user *model.User, secretKey []byte, duration time.Duration) (string, error) {
	return utils.GenerateToken(user, secretKey, duration)
}
