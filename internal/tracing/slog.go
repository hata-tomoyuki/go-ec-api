package tracing

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

// TraceHandler is a slog.Handler decorator that enriches log records
// with trace_id and span_id from the OpenTelemetry span in the context.
type TraceHandler struct {
	inner slog.Handler
}

// NewTraceHandler wraps an existing slog.Handler so that every log record
// automatically includes trace_id / span_id when a valid OTel span exists
// in the context.
func NewTraceHandler(inner slog.Handler) *TraceHandler {
	return &TraceHandler{inner: inner}
}

func (h *TraceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

// Handle extracts trace_id and span_id from the OTel span in ctx and adds
// them as slog attributes before delegating to the inner handler.
func (h *TraceHandler) Handle(ctx context.Context, record slog.Record) error {
	span := trace.SpanFromContext(ctx)

	if !span.SpanContext().IsValid() {
		return h.inner.Handle(ctx, record)
	}

	sc := span.SpanContext()
	record.AddAttrs(
		slog.String("trace_id", sc.TraceID().String()),
		slog.String("span_id", sc.SpanID().String()),
	)

	return h.inner.Handle(ctx, record)
}

func (h *TraceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TraceHandler{inner: h.inner.WithAttrs(attrs)}
}

func (h *TraceHandler) WithGroup(name string) slog.Handler {
	return &TraceHandler{inner: h.inner.WithGroup(name)}
}
