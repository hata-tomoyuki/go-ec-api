package tracing

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// MetricsMiddleware returns a chi-compatible middleware that records HTTP metrics.
// It creates two instruments:
//   - http_requests_total (Counter): counts requests by method, path, status
//   - http_request_duration_seconds (Histogram): measures latency by method, path
func MetricsMiddleware(meter metric.Meter) func(http.Handler) http.Handler {
	requestCount, _ := meter.Int64Counter("http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)
	requestDuration, _ := meter.Float64Histogram("http_request_duration_seconds",
		metric.WithDescription("HTTP request duration in seconds"),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := newResponseWriter(w)
			next.ServeHTTP(wrapped, r)

			// Resolve the chi route pattern (e.g. "/products/{id}") to avoid label explosion.
			pattern := chi.RouteContext(r.Context()).RoutePattern()
			if pattern == "" {
				pattern = r.URL.Path
			}

			duration := time.Since(start).Seconds()
			method := r.Method
			status := wrapped.statusCode

			methodAttr := attribute.String("method", method)
			pathAttr := attribute.String("path", pattern)
			statusAttr := attribute.Int("status", status)
			attrs := metric.WithAttributes(methodAttr, pathAttr, statusAttr)
			requestCount.Add(r.Context(), 1, attrs)
			requestDuration.Record(r.Context(), duration, attrs)
		})
	}
}
