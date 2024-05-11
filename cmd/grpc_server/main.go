package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userAPI "github.com/arifullov/auth/internal/api/user"
	userRepository "github.com/arifullov/auth/internal/repository/user"
	userService "github.com/arifullov/auth/internal/service/user"
	desc "github.com/arifullov/auth/pkg/user_v1"
)

const grpcPort = 50052
const dbDSN = "host=localhost port=5001 dbname=auth_db user=auth password=secret_pass sslmode=disable"

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

	userRepo := userRepository.NewRepository(dbPool)
	userSrv := userService.NewUserService(userRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, userAPI.NewImplementation(userSrv))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
