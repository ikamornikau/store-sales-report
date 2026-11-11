package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ikamornikau/store-sales-report/backend/internal/adapters/postgres"
	"github.com/ikamornikau/store-sales-report/backend/internal/config"
	"github.com/ikamornikau/store-sales-report/backend/internal/handlers"
	"github.com/ikamornikau/store-sales-report/backend/internal/handlers/reports"
	"github.com/ikamornikau/store-sales-report/backend/internal/handlers/stores"
	"github.com/ikamornikau/store-sales-report/backend/internal/logs"
	"github.com/ikamornikau/store-sales-report/backend/internal/middlewares"
)

func main() {
	slog.SetDefault(logs.NewDefaultLogger())
	if err := run(); err != nil {
		slog.Error("failed to run server", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load[Config]("/secrets/config.yaml")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.Postgres.URL)
	if err != nil {
		return fmt.Errorf("failed to create postgres pool: %w", err)
	}
	defer pool.Close()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(timeoutCtx); err != nil {
		return fmt.Errorf("failed to ping postgres: %w", err)
	}

	serveMux := http.NewServeMux()
	postgres := postgres.New(pool)

	middlewares := []func(http.Handler) http.Handler{
		middlewares.WithRequestID,
		middlewares.WithStatusCodes,
		middlewares.WithRequestDuration,
		middlewares.WithRecover,
	}

	storesHandler := stores.NewHandler(postgres)
	storesHandler.Mount("/api/v1/stores", serveMux, middlewares...)

	reportsHandler := reports.NewHandler(postgres)
	reportsHandler.Mount("/api/v1/stores", serveMux, middlewares...)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serveMux.Handle("GET /ready", handlers.NewReadyHandler(ctx))
	serveMux.Handle("GET /health", handlers.NewHealthHandler())
	serveMux.Handle("GET /metrics", promhttp.Handler())

	server := &http.Server{
		Addr:              cfg.Server.Address,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
		Handler:           serveMux,
	}

	errChan := make(chan error, 1)
	go func() {
		slog.Info("starting server", slog.String("address", cfg.Server.Address))
		errChan <- server.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to listen and serve: %w", err)
		}

	case <-ctx.Done():
		slog.Info("shutting down server")
		if err := server.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}

		slog.Info("server shutdown complete")
	}

	return nil
}
