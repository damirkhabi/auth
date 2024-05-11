package user

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/arifullov/auth/internal/client/db"
	"github.com/arifullov/auth/internal/model"
	"github.com/arifullov/auth/internal/repository"
	"github.com/arifullov/auth/internal/repository/user/converter"
	modelRepo "github.com/arifullov/auth/internal/repository/user/model"
)

const (
	tableName = "users"

	idColumn           = "id"
	nameColumn         = "name"
	emailColumn        = "email"
	roleColumn         = "role"
	passwordHashColumn = "password_hash"
	createdAtColumn    = "created_at"
	updatedAtColumn    = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, user *model.CreateUser) (int64, error) {
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordHashColumn, roleColumn, createdAtColumn, updatedAtColumn).
		Values(user.Name, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var userID int64
	var pgErr *pgconn.PgError
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)
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
	builderSelect := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(user), nil
}

func (r *repo) Update(ctx context.Context, user *model.UpdateUser) error {
	builderUpdate := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: user.ID})
	if user.Name.Valid {
		builderUpdate = builderUpdate.Set(nameColumn, user.Name.String)
	}
	if user.Email.Valid {
		builderUpdate = builderUpdate.Set(emailColumn, user.Email.String)
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	var pgErr *pgconn.PgError
	res, err := r.db.DB().ExecContext(ctx, q, args...)
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
	builderDelete := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return model.ErrUserNotFound
	}
	return nil
}
