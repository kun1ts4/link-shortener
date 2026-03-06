package main

import (
	"context"
	"errors"
	"fmt"
	"link-shortener/internal/config"
	"link-shortener/internal/domain"
	"link-shortener/internal/handlers"
	"link-shortener/internal/logger"
	"link-shortener/internal/random"
	"link-shortener/internal/repo/memory"
	"link-shortener/internal/repo/postgres"
	"link-shortener/internal/repo/postgres/migrations"
	"link-shortener/internal/usecase"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const configPath = "./config/config.yaml"

func main() {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}
	log := logger.NewLogger(cfg.Env)

	log.Info("starting link-shortener",
		"env", cfg.Env,
		"storage", cfg.Storage.Type,
	)

	linkGen := random.NewLinkGen(cfg.Shortener)

	var repo domain.Repository
	switch cfg.Storage.Type {
	case "postgres":
		repo, err = waitForPostgres(log, cfg.Storage.Postgres)
		if err != nil {
			log.Error("failed to connect to postgres", "error", err)
			os.Exit(1)
		}
		log.Info("connected to postgres")

		log.Info("running migrations")
		if err = migrations.Migrate(cfg.Storage.Postgres.DSN); err != nil {
			log.Error("migration failed", "error", err)
			os.Exit(1)
		}
		log.Info("migrations applied")

	case "memory":
		repo = memory.NewMemRepository(cfg.Storage.Memory)
		log.Info("using in-memory storage", "max_size", cfg.Storage.Memory.MaxSize)
	}

	uc := usecase.NewShortenerUseCase(linkGen, repo)

	handler := handlers.NewHandler(uc, log)
	r := handler.NewRouter()

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	readTimeout, writeTimeout, idleTimeout := parseHTTPTimeouts(log, cfg.HTTP)

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("server listening", "addr", addr)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Error("shutdown error", "error", err)
	} else {
		log.Info("server stopped gracefully")
	}

	type closer interface{ Close() }
	if c, ok := repo.(closer); ok {
		c.Close()
		log.Info("repository closed")
	}
}

func parseHTTPTimeouts(log *slog.Logger, cfg config.HTTPConfig) (read, write, idle time.Duration) {
	parse := func(s, name string) time.Duration {
		d, err := time.ParseDuration(s)
		if err != nil {
			log.Error("invalid http timeout", "field", name, "value", s, "error", err)
			os.Exit(1)
		}
		return d
	}
	return parse(cfg.ReadTimeout, "read_timeout"),
		parse(cfg.WriteTimeout, "write_timeout"),
		parse(cfg.IdleTimeout, "idle_timeout")
}

func waitForPostgres(log *slog.Logger, cfg config.PostgresConfig) (domain.Repository, error) {
	retryDelay, err := time.ParseDuration(cfg.RetryDelay)
	if err != nil {
		log.Error("invalid retry_delay", "value", cfg.RetryDelay, "error", err)
		os.Exit(1)
	}

	var repo *postgres.PgRepository
	for i := 1; i <= cfg.RetryCount; i++ {
		log.Info("connecting to postgres", "attempt", i, "max", cfg.RetryCount)
		repo, err = postgres.NewPostgresRepository(cfg)
		if err == nil {
			return repo, nil
		}
		log.Warn("postgres not ready", "attempt", i, "error", err)
		time.Sleep(retryDelay)
	}
	return nil, fmt.Errorf("postgres unavailable: %w", err)
}
