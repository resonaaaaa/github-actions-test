package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"example.com/go_backend/internal/service"
)

var ErrInvalidRequestBody = errors.New("invalid request body")

type Error struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

func MapHTTPError(err error) (int, Error) {
	switch {
	case errors.Is(err, service.ErrIdempotencyConflict):
		return http.StatusConflict, Error{
			Message: "idempotency conflict",
			Code:    "IDEMPOTENCY_CONFLICT",
		}
	case errors.Is(err, service.ErrInvalidRequestID):
		return http.StatusBadRequest, Error{
			Message: "invalid request ID",
			Code:    "INVALID_REQUEST_ID",
		}
	case errors.Is(err, service.ErrInvalidUserID):
		return http.StatusBadRequest, Error{
			Message: "invalid user ID",
			Code:    "INVALID_USER_ID",
		}
	case errors.Is(err, service.ErrInvalidProductID):
		return http.StatusBadRequest, Error{
			Message: "invalid product ID",
			Code:    "INVALID_PRODUCT_ID",
		}
	case errors.Is(err, service.ErrInsufficientStock):
		return http.StatusConflict, Error{
			Message: "insufficient stock",
			Code:    "INSUFFICIENT_STOCK",
		}
	case errors.Is(err, service.ErrInvalidQuantity):
		return http.StatusBadRequest, Error{
			Message: "invalid quantity",
			Code:    "INVALID_QUANTITY",
		}
	case errors.Is(err, service.ErrProductNotFound):
		return http.StatusNotFound, Error{
			Message: "product not found",
			Code:    "PRODUCT_NOT_FOUND",
		}
	case errors.Is(err, service.ErrOrderNotFound):
		return http.StatusNotFound, Error{
			Message: "order not found",
			Code:    "ORDER_NOT_FOUND",
		}
	case errors.Is(err, ErrInvalidRequestBody):
		return http.StatusBadRequest, Error{
			Message: "invalid request body",
			Code:    "INVALID_REQUEST_BODY",
		}

	default:
		return http.StatusInternalServerError, Error{
			Message: "internal server error",
			Code:    "INTERNAL_SERVER_ERROR",
		}
	}
}

type StatusRecorder struct {
	http.ResponseWriter
	StatusCode int
	written    bool
}

func (sr *StatusRecorder) WriteHeader(code int) {
	if sr.written {
		return
	}
	sr.written = true
	sr.StatusCode = code
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *StatusRecorder) Write(b []byte) (int, error) {
	if !sr.written {
		sr.WriteHeader(http.StatusOK)
	}
	return sr.ResponseWriter.Write(b)
}

func writeErrorResponse(w http.ResponseWriter, err error) {
	status, errorBody := MapHTTPError(err)
	response := ErrorResponse{
		Error: errorBody,
	}
	err_json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(err_json)
	if err != nil {
		fmt.Printf("failed to write error response: %v", err)
		return
	}
}
