package app

import (
	"context"

	"github.com/arifullov/auth/internal/api/access"
	"github.com/arifullov/auth/internal/api/auth"
	"github.com/arifullov/auth/internal/api/user"
	"github.com/arifullov/auth/internal/client/db"
	"github.com/arifullov/auth/internal/client/db/pg"
	"github.com/arifullov/auth/internal/client/db/transaction"
	"github.com/arifullov/auth/internal/closer"
	"github.com/arifullov/auth/internal/config"
	"github.com/arifullov/auth/internal/logger"
	"github.com/arifullov/auth/internal/repository"
	"github.com/arifullov/auth/internal/service"

	accessRepository "github.com/arifullov/auth/internal/repository/access"
	userRepository "github.com/arifullov/auth/internal/repository/user"
	userService "github.com/arifullov/auth/internal/service/user"

	accessService "github.com/arifullov/auth/internal/service/access"

	authService "github.com/arifullov/auth/internal/service/auth"
)

type serviceProvider struct {
	pgConfig         config.PGConfig
	grpcConfig       config.GRPCConfig
	httpConfig       config.HTTPConfig
	swaggerConfig    config.SwaggerConfig
	prometheusConfig config.HTTPConfig
	tokenConfig      config.TokenConfig
	loggerConfig     config.LoggerConfig
	jaegerConfig     config.JaegerConfig

	dbClient         db.Client
	txManager        db.TxManager
	userRepository   repository.UserRepository
	accessRepository repository.AccessRepository

	userService   service.UserService
	accessService service.AccessService
	authService   service.AuthService

	userImpl  *user.Implementation
	authImpl  *auth.Implementation
	accessImp *access.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		pgConfig, err := config.NewPGConfig()
		if err != nil {
			logger.Fatalf("failed to get pg config: %s", err.Error())
		}
		s.pgConfig = pgConfig
	}
	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		grpcConfig, err := config.NewGRPCConfig()
		if err != nil {
			logger.Fatalf("failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = grpcConfig
	}
	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		httpConfig, err := config.NewHTTPConfig()
		if err != nil {
			logger.Fatalf("failed to get http config: %s", err.Error())
		}
		s.httpConfig = httpConfig
	}
	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			logger.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) PrometheusConfig() config.PrometheusConfig {
	if s.prometheusConfig == nil {
		prometheusConfig, err := config.NewPrometheusConfig()
		if err != nil {
			logger.Fatalf("failed to get prometheus config: %s", err.Error())
		}
		s.prometheusConfig = prometheusConfig
	}
	return s.prometheusConfig
}

func (s *serviceProvider) TokenConfig() config.TokenConfig {
	if s.tokenConfig == nil {
		cfg, err := config.NewTokenConfig()
		if err != nil {
			logger.Fatalf("failed to get token config: %s", err.Error())
		}

		s.tokenConfig = cfg
	}

	return s.tokenConfig
}

func (s *serviceProvider) LoggerConfig() config.LoggerConfig {
	if s.loggerConfig == nil {
		cfg, err := config.NewLoggingConfig()
		if err != nil {
			logger.Fatalf("failed to get logging config: %s", err.Error())
		}

		s.loggerConfig = cfg
	}

	return s.loggerConfig
}

func (s *serviceProvider) JaegerConfig() config.JaegerConfig {
	if s.jaegerConfig == nil {
		cfg, err := config.NewJaegerConfig()
		if err != nil {
			logger.Fatalf("failed to get jaeger config: %s", err.Error())
		}
		s.jaegerConfig = cfg
	}
	return s.jaegerConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			logger.Fatalf("failed to create db client: %v", err)
		}

		if err = cl.DB().Ping(ctx); err != nil {
			logger.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}
	return s.dbClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}
	return s.userRepository
}

func (s *serviceProvider) AccessRepository(ctx context.Context) repository.AccessRepository {
	if s.accessRepository == nil {
		s.accessRepository = accessRepository.NewRepository(s.DBClient(ctx))
	}
	return s.accessRepository
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}
	return s.txManager
}

func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewAccessService(
			s.AccessRepository(ctx),
			s.TokenConfig().AccessTokenSecretKey(),
		)
	}
	return s.accessService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewAuthService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
			s.TokenConfig(),
		)
	}
	return s.authService
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewUserService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}
	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}
	return s.userImpl
}

func (s *serviceProvider) AccessImpl(ctx context.Context) *access.Implementation {
	if s.accessImp == nil {
		s.accessImp = access.NewImplementation(s.AccessService(ctx))
	}
	return s.accessImp
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}
	return s.authImpl
}
