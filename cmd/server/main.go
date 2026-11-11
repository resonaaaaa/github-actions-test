package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/go_backend/internal/application"
	"example.com/go_backend/internal/config"
	"example.com/go_backend/internal/repository"
	"example.com/go_backend/internal/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func OpenDB(cfg config.Config) (*sql.DB, error) {
	dsn := cfg.MySQL.DSN
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	db.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MySQL.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func run() error {
	// 接收系统信号
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	if err := godotenv.Load(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		return err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := OpenDB(cfg)
	if err != nil {
		logger.Error("failed to open database", "error", err)
		return err
	}
	defer db.Close()

	// Initialize repositories
	uow := repository.NewUnitOfWork(db)
	orderReader := repository.NewOrderRepository(db)

	// Initialize services
	orderService := service.NewOrderService(uow, orderReader)

	// Initialize http handlers
	orderHandler := application.NewOrderHandler(orderService)

	// Initialize http server
	mux := http.NewServeMux()
	mux.Handle("/orders", application.LoggingMiddleware(logger)(http.HandlerFunc(orderHandler.CreateOrder)))

	// Start http server
	server := &http.Server{
		Addr:              cfg.HTTP.Addr,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		Handler:           mux,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("starting server", "addr", cfg.HTTP.Addr)
		errCh <- server.ListenAndServe()
	}()

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		stop()

		logger.Info("shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("graceful shutdown failed", "error", err)

			if err := server.Close(); err != nil {
				logger.Error("server close failed", "error", err)
				return fmt.Errorf("server close failed: %w", err)
			}

			return fmt.Errorf("server shutdown failed: %w", err)
		}

	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
			return fmt.Errorf("server error: %w", err)
		}
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		os.Exit(1)
	}
}
