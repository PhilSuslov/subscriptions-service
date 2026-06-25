package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/subscriptions-service/internal/config"
	"github.com/example/subscriptions-service/internal/infrastructure/logger"
	"github.com/example/subscriptions-service/internal/repository/postgres"
	app "github.com/example/subscriptions-service/internal/service/subscription"
	httptransport "github.com/example/subscriptions-service/internal/transport/http"
	"github.com/example/subscriptions-service/internal/transport/http/handler"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	newPool = func(ctx context.Context, cfg config.PostgresConfig) (postgres.Queryer, error) {
		return newPoolAdapter(ctx, cfg)
	}
	newRepo             = postgres.NewSubscriptionRepository
	newUseCase          = app.NewUseCase
	newHandler          = handler.NewSubscriptionHandler
	newRouter           = httptransport.NewRouter
	configLoad          = config.Load
	signalNotifyContext = signal.NotifyContext
	newServerFn         = func(addr string, h http.Handler) *http.Server {
		return &http.Server{Addr: addr, Handler: h, ReadHeaderTimeout: 5 * time.Second}
	}
)

type poolCloser interface{ Close() }

type poolAdapter struct{ *pgxpool.Pool }

func newPoolAdapter(ctx context.Context, cfg config.PostgresConfig) (postgres.Queryer, error) {
	pool, err := postgres.NewPool(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return poolAdapter{Pool: pool}, nil
}

func main() {
	if err := execute(); err != nil {
		logger.New().Error("application failed", "error", err)
	}
}

func execute() error {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	log := logger.New()
	cfg, err := configLoad(*configPath)
	if err != nil {
		return err
	}

	return run(cfg, log)
}

func run(cfg *config.Config, log *slog.Logger) error {
	ctx, stop := signalNotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := newPool(ctx, cfg.Postgres)
	if err != nil {
		return err
	}
	defer func() {
		if closer, ok := pool.(poolCloser); ok {
			closer.Close()
		}
	}()

	repo := newRepo(pool)
	uc := newUseCase(repo)
	h := newHandler(uc, log)
	router := newRouter(h, log)

	srv := newServerFn(cfg.HTTP.Addr, router)
	go func() {
		log.Info("http server started", "addr", cfg.HTTP.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	log.Info("application stopped")
	return nil
}
