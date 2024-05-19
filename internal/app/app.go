package app

import (
	"context"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/arifullov/auth/internal/closer"
	"github.com/arifullov/auth/internal/config"
	"github.com/arifullov/auth/internal/interceptor"
	"github.com/arifullov/auth/internal/logger"
	"github.com/arifullov/auth/internal/metric"
	"github.com/arifullov/auth/internal/tracing"
	descAccess "github.com/arifullov/auth/pkg/access_v1"
	descAuth "github.com/arifullov/auth/pkg/auth_v1"
	descUser "github.com/arifullov/auth/pkg/user_v1"
	_ "github.com/arifullov/auth/statik"
)

type App struct {
	serviceProvider  *serviceProvider
	grpcServer       *grpc.Server
	httpServer       *http.Server
	swaggerServer    *http.Server
	prometheusServer *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	if err := app.initDeps(ctx); err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()

		if err := a.runGRPCServer(); err != nil {
			logger.Fatalf("failed to run grpc server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		if err := a.runHTTPServer(); err != nil {
			logger.Fatalf("failed to run http server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		if err := a.runSwaggerServer(); err != nil {
			logger.Fatalf("failed to run swagger server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		if err := a.runPrometheusServer(); err != nil {
			logger.Fatalf("failed to run prometheus server: %v", err)
		}
	}()

	wg.Wait()
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initMetric,
		a.initConfig,
		a.initServiceProvider,
		a.initLogger,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
		a.initPrometheusServer,
		a.initTracing,
	}
	for _, init := range inits {
		if err := init(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	if err := config.Load("./.env"); err != nil {
		return err
	}
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		logging.WithDurationField(logging.DurationToDurationField),
	}

	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptor.InterceptorLogger(logger.GetLogger()), opts...),
			interceptor.MetricsInterceptor,
			interceptor.ValidateInterceptor,
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	reflection.Register(a.grpcServer)
	descUser.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))
	descAccess.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImpl(ctx))
	descAuth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := descUser.RegisterUserV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: corsMiddleware.Handler(mux),
	}
	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:    a.serviceProvider.SwaggerConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) initPrometheusServer(_ context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.prometheusServer = &http.Server{
		Addr:    a.serviceProvider.PrometheusConfig().Address(),
		Handler: mux,
	}
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	cfg := a.serviceProvider.LoggerConfig()
	return logger.Init(cfg.Level(), cfg.Format())
}

func (a *App) initMetric(ctx context.Context) error {
	return metric.Init(ctx)
}

func (a *App) initTracing(_ context.Context) error {
	return tracing.Init(
		a.serviceProvider.JaegerConfig().CollectorEndpoint(),
		a.serviceProvider.JaegerConfig().ServiceName(),
		a.serviceProvider.JaegerConfig().DeploymentEnvironment(),
	)
}

func (a *App) runGRPCServer() error {
	logger.Infof("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())
	lis, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	if err = a.grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (a *App) runHTTPServer() error {
	logger.Infof("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address())

	if err := a.httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (a *App) runSwaggerServer() error {
	logger.Infof("Swagger server is running on %s", a.serviceProvider.SwaggerConfig().Address())

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runPrometheusServer() error {
	logger.Infof("Prometheus server is running on %s", a.serviceProvider.PrometheusConfig().Address())

	err := a.prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Infof("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		logger.Infof("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Infof("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Infof("Served swagger file: %s", path)
	}
}
