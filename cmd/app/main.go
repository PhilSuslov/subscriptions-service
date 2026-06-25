package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/subscriptions-service/internal/config"
	"github.com/example/subscriptions-service/internal/infrastructure/logger"
	"github.com/example/subscriptions-service/internal/repository/postgres"
	app "github.com/example/subscriptions-service/internal/service/subscription"
	httptransport "github.com/example/subscriptions-service/internal/transport/http"
	"github.com/example/subscriptions-service/internal/transport/http/handler"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	log := logger.New()
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Error("load config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		log.Error("connect postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	repo := postgres.NewSubscriptionRepository(pool)
	uc := app.NewUseCase(repo)
	h := handler.NewSubscriptionHandler(uc, log)
	router := httptransport.NewRouter(h, log)

	srv := &http.Server{Addr: cfg.HTTP.Addr, Handler: router, ReadHeaderTimeout: 5 * time.Second}
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
		log.Error("graceful shutdown failed", slog.Any("error", err))
		return
	}
	log.Info("application stopped")
}
