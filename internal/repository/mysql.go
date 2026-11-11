package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/go_backend/internal/service"

	"github.com/go-sql-driver/mysql"
)

type UnitOfWork struct {
	db *sql.DB
}

func (uow *UnitOfWork) DoWithinTx(
	ctx context.Context,
	fn func(ctx context.Context, repos *service.Repositories) error,
) error {
	tx, err := uow.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	repos := &service.Repositories{
		Product: NewProductRepository(tx),
		Order:   NewOrderRepository(tx),
	}
	if err := fn(ctx, repos); err != nil {
		return fmt.Errorf("failed to execute function within transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func IsDuplicateKey(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}

func IsDeadlock(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1213
}

func (uow *UnitOfWork) DoWithinTxWithRetry(
	ctx context.Context,
	fn func(ctx context.Context, repos *service.Repositories) error,
	maxRetries int,
) error {
	for retry := 0; ; retry++ {
		err := uow.DoWithinTx(ctx, fn)
		if err == nil {
			return nil
		}
		if !IsDeadlock(err) {
			return fmt.Errorf("failed to execute function within transaction: %w", err)
		}
		if retry >= maxRetries {
			return fmt.Errorf("failed to execute function within transaction after %d retries: %w: %w", maxRetries, service.ErrMaxRetries, err)
		}

	}
}

func NewUnitOfWork(db *sql.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}
