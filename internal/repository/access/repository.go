package access

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/arifullov/auth/internal/client/db"
	"github.com/arifullov/auth/internal/model"
	"github.com/arifullov/auth/internal/repository"
	"github.com/arifullov/auth/internal/repository/access/converter"
	modelRepo "github.com/arifullov/auth/internal/repository/access/model"
)

const (
	routeAccessesTable = "route_accesses"

	roleColumn  = "role"
	routeColumn = "route"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AccessRepository {
	return &repo{
		db: db,
	}
}

func (r repo) GetRouteRoles(ctx context.Context, route string) ([]model.Role, error) {
	builderSelect := sq.Select(roleColumn).
		PlaceholderFormat(sq.Dollar).
		From(routeAccessesTable).
		Where(sq.Eq{routeColumn: route})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "access_repository.GetRouteRoles",
		QueryRaw: query,
	}

	var routeAccesses []modelRepo.RouteAccesses
	err = r.db.DB().ScanAllContext(ctx, &routeAccesses, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToRolesFromRepo(routeAccesses), nil
}
