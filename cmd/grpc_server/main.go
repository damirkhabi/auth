package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/mail"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/arifullov/auth/pkg/user_v1"
)

const grpcPort = 50052
const dbDSN = "host=localhost port=5001 dbname=auth_db user=auth password=secret_pass sslmode=disable"
const Iter = 20000

type server struct {
	desc.UnimplementedUserV1Server

	dbPool *pgxpool.Pool
}

const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const SaltSize = 12

func getRandomSalt(size int) []byte {
	salt := make([]byte, size)
	l := len(allowedChars)
	for i := range salt {
		salt[i] = allowedChars[rand.Intn(l)]
	}
	return salt
}

func checkPbkdf2SHA256(password, encoded string) (bool, error) {
	parts := strings.SplitN(encoded, "$", 4)
	if len(parts) != 4 {
		return false, errors.New("Hash must consist of 4 segments")
	}
	iter, err := strconv.Atoi(parts[1])
	if err != nil {
		return false, fmt.Errorf("Wrong number of iterations: %v", err)
	}
	salt := []byte(parts[2])
	k, err := base64.StdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, fmt.Errorf("Wrong hash encoding: %v", err)
	}
	dk := pbkdf2.Key([]byte(password), salt, iter, sha256.Size, sha256.New)
	return bytes.Equal(k, dk), nil
}

func makePbkdf2SHA256(password string) string {
	salt := getRandomSalt(SaltSize)
	dk := pbkdf2.Key([]byte(password), salt, Iter, sha256.Size, sha256.New)
	b64Hash := base64.StdEncoding.EncodeToString(dk)
	return fmt.Sprintf("%s$%d$%s$%s", "pbkdf2_sha256", Iter, salt, b64Hash)
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	now := time.Now()
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email")
	}

	if req.Password != req.PasswordConfirm {
		return nil, status.Errorf(codes.InvalidArgument, "password mismatch")
	}
	passwordHash := makePbkdf2SHA256(req.Password)
	role := "user"
	if req.Role == desc.UserRole_ADMIN {
		role = "admin"
	}

	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password_hash", "role", "created_at", "updated_at").
		Values(req.Name, req.Email, passwordHash, role, now, now).
		Suffix("RETURNING id")
	log.Printf("create user: %v", req)
	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to build query")
	}

	var userID int64
	err = s.dbPool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Printf("failed to insert user: %v", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, status.Errorf(codes.AlreadyExists, "user already exists")
			}
		}
		return nil, status.Errorf(codes.Unavailable, "failed to insert user")
	}
	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("get user: %v", req)
	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
	}

	userResponse := &desc.GetResponse{}
	var role string
	var createdAt, updatedAt time.Time
	err = s.dbPool.QueryRow(ctx, query, args...).Scan(
		&userResponse.Id, &userResponse.Name, &userResponse.Email,
		&role, &createdAt, &updatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if err != nil {
		log.Printf("failed to select notes: %v", err)
		return nil, status.Errorf(codes.Unavailable, "user select error")
	}
	userResponse.CreatedAt = timestamppb.New(createdAt)
	userResponse.UpdatedAt = timestamppb.New(updatedAt)
	if role == "admin" {
		userResponse.Role = desc.UserRole_ADMIN
	}

	return userResponse, nil
}
func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": req.Id})
	if req.Name != nil {
		builderUpdate = builderUpdate.Set("name", req.Name.Value)
	}
	if req.Email != nil {
		if _, err := mail.ParseAddress(req.Email.Value); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid email")
		}
		builderUpdate = builderUpdate.Set("email", req.Email.Value)
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to update user")
	}

	res, err := s.dbPool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to update user")
	}
	if res.RowsAffected() == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return nil, nil
}
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	builderDelete := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to delete user")
	}

	res, err := s.dbPool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to update user")
	}
	if res.RowsAffected() == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	log.Printf("delete user: %v", req)
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed db connect: %v", err)
	}
	if err = dbPool.Ping(ctx); err != nil {
		log.Fatalf("failed db connect: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{dbPool: dbPool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
