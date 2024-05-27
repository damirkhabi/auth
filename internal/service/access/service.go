package access

import (
	"context"

	"github.com/arifullov/auth/internal/repository"
	"github.com/arifullov/auth/internal/service"
	"github.com/arifullov/auth/internal/sys"
	"github.com/arifullov/auth/internal/sys/codes"
	"github.com/arifullov/auth/internal/utils"
)

type serv struct {
	accessRepository     repository.AccessRepository
	accessTokenSecretKey string
}

func NewAccessService(
	accessRepository repository.AccessRepository,
	accessTokenSecretKey string,
) service.AccessService {
	return &serv{
		accessRepository:     accessRepository,
		accessTokenSecretKey: accessTokenSecretKey,
	}
}

func (s *serv) Check(ctx context.Context, accessToken string, endpointAddress string) error {
	claims, err := utils.VerifyToken(accessToken, utils.S2B(s.accessTokenSecretKey))
	if err != nil {
		return err
	}

	roles, err := s.accessRepository.GetRouteRoles(ctx, endpointAddress)
	if err != nil {
		return err
	}

	for _, role := range roles {
		if role == claims.Role {
			return nil
		}
	}

	return sys.NewCommonError(codes.PermissionDenied, "permission denied")
}
