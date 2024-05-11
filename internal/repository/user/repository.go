package user

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/arifullov/auth/internal/model"
	"github.com/arifullov/auth/internal/repository"
	"github.com/arifullov/auth/internal/repository/user/converter"
	modelRepo "github.com/arifullov/auth/internal/repository/user/model"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, user *model.CreateUser) (int64, error) {
	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password_hash", "role", "created_at", "updated_at").
		Values(user.Name, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	var userID int64
	var pgErr *pgconn.PgError
	err = r.db.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil && errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return 0, model.ErrUserAlreadyExists
		}
	}
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"id": id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	var user modelRepo.User
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(user), nil
}

func (r *repo) Update(ctx context.Context, user *model.UpdateUser) error {
	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": user.ID})
	if user.Name.Valid {
		builderUpdate = builderUpdate.Set("name", user.Name.String)
	}
	if user.Email.Valid {
		builderUpdate = builderUpdate.Set("email", user.Email.String)
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	var pgErr *pgconn.PgError
	res, err := r.db.Exec(ctx, query, args...)
	if err != nil && errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return model.ErrUserAlreadyExists
		}
	}
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builderDelete := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return model.ErrUserNotFound
	}
	return nil
}
