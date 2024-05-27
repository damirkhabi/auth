package auth

import (
	"context"

	"github.com/arifullov/auth/internal/sys"
	"github.com/arifullov/auth/internal/sys/codes"
	"github.com/arifullov/auth/internal/utils"
)

func (s *serv) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, utils.S2B(s.tokenConfig.RefreshTokenSecretKey()))
	if err != nil {
		return "", sys.NewCommonError(codes.Unauthenticated, err.Error())
	}

	user, err := s.userRepository.GetByEmail(ctx, claims.Username)
	if err != nil {
		return "", err
	}

	accessToken, err := generateAccessToken(user, utils.S2B(s.tokenConfig.AccessTokenSecretKey()), s.tokenConfig.AccessTokenExpiration())
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
