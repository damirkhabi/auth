package auth

import (
	"context"

	"github.com/arifullov/auth/internal/utils"
)

func (s *serv) GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	claims, err := utils.VerifyToken(oldRefreshToken, utils.S2B(s.tokenConfig.RefreshTokenSecretKey()))
	if err != nil {
		return "", err
	}

	user, err := s.userRepository.GetByEmail(ctx, claims.Username)
	if err != nil {
		return "", err
	}

	refreshToken, err := generateRefreshToken(user, utils.S2B(s.tokenConfig.RefreshTokenSecretKey()), s.tokenConfig.RefreshTokenExpiration())
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}
