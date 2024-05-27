package model

import (
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	ID           int64
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Role     Role   `json:"role"`
}
