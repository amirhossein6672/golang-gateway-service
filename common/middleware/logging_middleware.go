package middleware

import (
	"log"
	"log/slog"
	"net/http"
	"time"
)

// wrappedWriter wraps ResponseWriter to capture the status code
type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

// Overrides WriteHeader to capture the status code
func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

// Logging middleware to log request details
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Capture request start time

		// Wrap ResponseWriter to track the status code
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r) // Proceed to the next handler

		// Log the request details and response status
		log.Println(
			"➡️",
			"handled request",
			slog.Int("statusCode", wrapped.statusCode),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Any("duration", time.Since(start)),
			slog.String("xff", r.Header.Get("X-Forwarded-For")),
		)
	})
}
