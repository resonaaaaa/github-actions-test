package service

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrIdempotencyConflict = errors.New("idempotency key reused with different request")
	ErrInvalidRequestID    = errors.New("invalid request ID")
	ErrInvalidUserID       = errors.New("invalid user ID")
	ErrInvalidProductID    = errors.New("invalid product ID")
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrInvalidQuantity     = errors.New("invalid quantity")
	ErrProductNotFound     = errors.New("product not found")
	ErrMaxRetries          = errors.New("failed to retry")
	ErrOrderNotFound       = errors.New("order not found")
)

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice int64     `json:"unit_price"`
	RequestID string    `json:"request_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Price_cents int64  `json:"price_cents"`
	Stock       int    `json:"stock"`
}

type ProductRepository interface {
	GetByID(
		ctx context.Context,
		id int64,
	) (*Product, error)

	DecreaseStock(
		ctx context.Context,
		id int64,
		quantity int,
	) error
}

type OrderRepository interface {
	Create(
		ctx context.Context,
		order *Order,
	) error

	GetByRequestID(
		ctx context.Context,
		userID int64,
		requestID string,
	) (*Order, error)
}

// 事务中使用的仓库集合
type Repositories struct {
	Product ProductRepository
	Order   OrderRepository
}

// 用于在事务中执行操作的接口
type UnitOfWork interface {
	DoWithinTx(
		ctx context.Context,
		fn func(ctx context.Context, repos *Repositories) error,
	) error
	DoWithinTxWithRetry(
		ctx context.Context,
		fn func(ctx context.Context, repos *Repositories) error,
		maxRetries int,
	) error
}

type OrderReader interface {
	GetByRequestID(
		ctx context.Context,
		userID int64,
		requestID string,
	) (*Order, error)
}

type OrderService struct {
	uow         UnitOfWork
	orderReader OrderReader
}

func NewOrderService(uow UnitOfWork, orderReader OrderReader) *OrderService {
	return &OrderService{
		uow:         uow,
		orderReader: orderReader,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID int64, productID int64, quantity int, requestID string) (*Order, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if productID <= 0 {
		return nil, ErrInvalidProductID
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	if requestID == "" || len(requestID) > 64 {
		return nil, ErrInvalidRequestID
	}
	var createdOrder *Order

	maxRetries := 5
	if err := s.uow.DoWithinTxWithRetry(ctx, func(ctx context.Context, repos *Repositories) error {
		order := &Order{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
			RequestID: requestID,
		}

		product, err := repos.Product.GetByID(ctx, productID)
		if err != nil {
			return err
		}
		if product == nil {
			return ErrProductNotFound
		}

		createdOrder = order
		createdOrder.UnitPrice = product.Price_cents

		if err := repos.Product.DecreaseStock(ctx, productID, quantity); err != nil {
			return fmt.Errorf("failed to decrease stock: %w", err)
		}

		if err := repos.Order.Create(ctx, createdOrder); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		return nil
	}, maxRetries); err != nil {
		if errors.Is(err, ErrIdempotencyConflict) {
			existingOrder, err := s.orderReader.GetByRequestID(ctx, userID, requestID)
			if err != nil && !errors.Is(err, ErrOrderNotFound) {
				return nil, fmt.Errorf("failed to check existing order: %w", err)
			}
			if existingOrder != nil {
				if existingOrder.ProductID != productID || existingOrder.Quantity != quantity {
					return nil, ErrIdempotencyConflict
				} else {
					return existingOrder, nil
				}
			}

		}
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return createdOrder, nil
}
