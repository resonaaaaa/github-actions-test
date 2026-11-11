package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/go_backend/internal/service"
)

type MySQLProductRepository struct {
	db DBTX
}

func NewProductRepository(db DBTX) *MySQLProductRepository {
	return &MySQLProductRepository{db: db}
}

func (p *MySQLProductRepository) GetByID(
	ctx context.Context,
	id int64,
) (*service.Product, error) {
	var product service.Product
	err := p.db.QueryRowContext(ctx, "SELECT id, name, price_cents, stock FROM products WHERE id = ?", id).Scan(&product.ID, &product.Name, &product.Price_cents, &product.Stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %d", service.ErrProductNotFound, id)
		}
		return nil, fmt.Errorf("failed to get product by ID: %d: %w", id, err)
	}
	return &product, nil
}

func (p *MySQLProductRepository) DecreaseStock(
	ctx context.Context,
	id int64,
	quantity int,
) error {
	if quantity <= 0 {
		return fmt.Errorf("%w: %d is invalid", service.ErrInvalidQuantity, quantity)
	}

	result, err := p.db.ExecContext(ctx, "UPDATE products SET stock = stock - ? WHERE id = ? AND stock >= ?", quantity, id, quantity)
	if err != nil {
		return fmt.Errorf("failed to decrease stock for product ID: %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for product ID: %d: %w", id, err)
	}

	switch rowsAffected {
	case 0:
		return fmt.Errorf("%w: product ID: %d, requested quantity: %d", service.ErrInsufficientStock, id, quantity)
	case 1:
		return nil
	default:
		return fmt.Errorf("unexpected rows affected for product ID: %d: %d", id, rowsAffected)
	}
}
