package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/go_backend/internal/service"
)

type FakeOrderService struct {
	Order *service.Order
	Err   error

	Called    bool
	UserID    int64
	ProductID int64
	Quantity  int
	RequestID string
}

func (f *FakeOrderService) CreateOrder(
	ctx context.Context,
	userID int64,
	productID int64,
	quantity int,
	requestID string,
) (*service.Order, error) {

	f.Called = true
	f.UserID = userID
	f.ProductID = productID
	f.Quantity = quantity
	f.RequestID = requestID

	return f.Order, f.Err
}

func TestCreateOrderHandler(t *testing.T) {
	fake := &FakeOrderService{
		Order: &service.Order{
			ID:        101,
			UserID:    1,
			ProductID: 10,
			Quantity:  2,
			UnitPrice: 9990,
			CreatedAt: time.Now(),
		},
	}

	handler := NewOrderHandler(fake)

	body := strings.NewReader(`{
        "user_id": 1,
        "product_id": 10,
        "quantity": 2
    }`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		body,
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "abc123")

	rec := httptest.NewRecorder()

	handler.CreateOrder(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			resp.StatusCode,
		)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected content type")
	}

	var result OrderResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if result.ID != 101 {
		t.Fatalf("expected order ID 101, got %d", result.ID)
	}

	if !fake.Called {
		t.Fatal("expected service to be called")
	}

	if fake.RequestID != "abc123" {
		t.Fatalf(
			"expected request ID abc123, got %s",
			fake.RequestID,
		)
	}
}

func TestCreateOrderHandler_InsufficientStock(t *testing.T) {
	fake := &FakeOrderService{
		Err: service.ErrInsufficientStock,
	}

	handler := NewOrderHandler(fake)

	body := strings.NewReader(`{
		"user_id": 1,
		"product_id": 10,
		"quantity": 2
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		body,
	)

	rec := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "abc123")

	handler.CreateOrder(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected content type")
	}

	var data ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if data.Error.Code != "INSUFFICIENT_STOCK" {
		t.Fatalf("expected error code INSUFFICIENT_STOCK, got %s", data.Error.Code)
	}
	if data.Error.Message != "insufficient stock" {
		t.Fatalf("expected error message 'insufficient stock', got %s", data.Error.Message)
	}

}

func TestCreateOrderHandler_LackOfIdempotencyKey(t *testing.T) {
	fake := &FakeOrderService{
		Order: &service.Order{
			ID:        101,
			UserID:    1,
			ProductID: 10,
			Quantity:  2,
			UnitPrice: 9990,
			CreatedAt: time.Now(),
		},
	}

	handler := NewOrderHandler(fake)

	body := strings.NewReader(`{
		"user_id": 1,
		"product_id": 10,
		"quantity": 2
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		body,
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.CreateOrder(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected content type")
	}

	var data ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if data.Error.Code != "INVALID_REQUEST_ID" {
		t.Fatalf("expected error code INVALID_REQUEST_ID, got %s", data.Error.Code)
	}
	if data.Error.Message != "invalid request ID" {
		t.Fatalf("expected error message 'invalid request ID', got %s", data.Error.Message)
	}
}

func TestCreateOrderHandler_InvalidRequestBody(t *testing.T) {
	fake := &FakeOrderService{
		Order: &service.Order{
			ID:        101,
			UserID:    1,
			ProductID: 10,
			Quantity:  2,
			UnitPrice: 9990,
			CreatedAt: time.Now(),
		},
	}

	handler := NewOrderHandler(fake)

	body := strings.NewReader(`{
		"user_id": 1,
		"product_id": 10,
		"quantity": 2,
		"request_id": "abc123"
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		body,
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "abc123")

	rec := httptest.NewRecorder()

	handler.CreateOrder(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected content type")
	}

	var data ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if data.Error.Code != "INVALID_REQUEST_BODY" {
		t.Fatalf("expected error code INVALID_REQUEST_BODY, got %s", data.Error.Code)
	}
	if data.Error.Message != "invalid request body" {
		t.Fatalf("expected error message 'invalid request body', got %s", data.Error.Message)
	}

}
