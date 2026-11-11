package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/go_backend/internal/service"
)

type OrderService interface {
	CreateOrder(
		ctx context.Context,
		userID int64,
		productID int64,
		quantity int,
		requestID string,
	) (*service.Order, error)
}

type OrderHandler struct {
	service OrderService
}

func NewOrderHandler(service OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq OrderRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&orderReq); err != nil {
		writeErrorResponse(w, fmt.Errorf("decode: %w", ErrInvalidRequestBody))
		return
	}

	requestID := r.Header.Get("Idempotency-Key")

	if requestID == "" {
		writeErrorResponse(w, fmt.Errorf("missing Idempotency-Key header: %w", service.ErrInvalidRequestID))
		return
	}

	order, err := h.service.CreateOrder(r.Context(), orderReq.UserID, orderReq.ProductID, orderReq.Quantity, requestID)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	response := OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		ProductID: order.ProductID,
		Quantity:  order.Quantity,
		UnitPrice: order.UnitPrice,
		CreatedAt: order.CreatedAt,
	}

	body, err := json.Marshal(response)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("marshal: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(body)
	if err != nil {
		fmt.Printf("failed to write response: %v", err)
		return
	}

}
