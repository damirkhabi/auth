package access

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/arifullov/auth/internal/model"
	desc "github.com/arifullov/auth/pkg/access_v1"
)

const (
	authPrefix = "Bearer "
)

func (i *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	err := i.accessService.Check(ctx, accessToken, req.GetEndpointAddress())
	if errors.Is(err, model.ErrAccessDenied) {
		return nil, status.Errorf(codes.PermissionDenied, "access denied")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get accessible roles")
	}
	return &emptypb.Empty{}, nil
}
