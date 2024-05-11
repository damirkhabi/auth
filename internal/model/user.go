package model

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrInvalidEmail     = errors.New("invalid email")
	ErrPasswordMismatch = errors.New("password mismatch")
)

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type Role string

type CreateUser struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	PasswordHash    string
	Role            Role
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type UpdateUser struct {
	ID    int64
	Name  sql.NullString
	Email sql.NullString
}

type User struct {
	ID        int64
	Name      string
	Email     string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}
