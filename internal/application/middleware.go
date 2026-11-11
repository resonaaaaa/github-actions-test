package application

import (
	"log/slog"
	"net/http"
)

func LoggingMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Incoming request",
				"method", r.Method,
				"url", r.URL.String(),
				"remote_addr", r.RemoteAddr,
			)
			next.ServeHTTP(w, r)
		})
	}
}
