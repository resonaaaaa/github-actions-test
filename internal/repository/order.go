package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/go_backend/internal/service"
)

type MySQLOrderRepository struct {
	db DBTX
}

func NewOrderRepository(db DBTX) *MySQLOrderRepository {
	return &MySQLOrderRepository{db: db}
}

func (s *MySQLOrderRepository) Create(
	ctx context.Context,
	order *service.Order,
) error {
	result, err := s.db.ExecContext(ctx, "INSERT INTO orders (user_id, product_id, quantity, unit_price, request_id) VALUES (?, ?, ?, ?, ?)",
		order.UserID,
		order.ProductID,
		order.Quantity,
		order.UnitPrice,
		order.RequestID,
	)
	if err != nil {
		if IsDuplicateKey(err) {
			return fmt.Errorf("%w: user ID: %d, request ID: %s", service.ErrIdempotencyConflict, order.UserID, order.RequestID)
		}
		return fmt.Errorf("failed to create order: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID for order: %w", err)
	}
	order.ID = id
	return nil
}

func (s *MySQLOrderRepository) GetByRequestID(
	ctx context.Context,
	userID int64,
	requestID string,
) (*service.Order, error) {
	var order service.Order
	row := s.db.QueryRowContext(ctx, "SELECT id, user_id, product_id, quantity, unit_price, request_id, created_at FROM orders WHERE user_id = ? AND request_id = ?", userID, requestID)
	err := row.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.UnitPrice, &order.RequestID, &order.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order by request ID: %w", err)
	}
	return &order, nil
}

type MySQLOrderReader struct {
	db DBTX
}

func NewOrderReader(db DBTX) *MySQLOrderReader {
	return &MySQLOrderReader{db: db}
}

func (r *MySQLOrderReader) GetByRequestID(
	ctx context.Context,
	userID int64,
	requestID string,
) (*service.Order, error) {
	var order service.Order
	row := r.db.QueryRowContext(ctx, "SELECT id, user_id, product_id, quantity, unit_price, request_id, created_at FROM orders WHERE user_id = ? AND request_id = ?", userID, requestID)
	err := row.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.UnitPrice, &order.RequestID, &order.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order by request ID: %w", err)
	}
	return &order, nil
}
