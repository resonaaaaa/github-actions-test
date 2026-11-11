package application

import (
	"time"
)

type OrderRequest struct {
	UserID    int64 `json:"user_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// 我们直接把RequestID放在请求头里, 而不是写入请求体。
type OrderResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice int64     `json:"unit_price"`
	CreatedAt time.Time `json:"created_at"`
}
