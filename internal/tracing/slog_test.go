package tracing_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"example.com/ecommerce/internal/tracing"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestTraceHandler_WithValidSpan(t *testing.T) {
	// Set up an in-memory span exporter so we get a real, valid SpanContext.
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	defer tp.Shutdown(context.Background())

	ctx, span := tp.Tracer("test").Start(context.Background(), "test-span")
	defer span.End()

	sc := span.SpanContext()
	if !sc.IsValid() {
		t.Fatal("expected a valid SpanContext from the SDK tracer")
	}

	var buf bytes.Buffer
	inner := slog.NewTextHandler(&buf, nil)
	logger := slog.New(tracing.NewTraceHandler(inner))

	logger.ErrorContext(ctx, "something failed", "error", "not found")

	output := buf.String()

	if !bytes.Contains([]byte(output), []byte("trace_id=")) {
		t.Errorf("expected trace_id in log output, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("span_id=")) {
		t.Errorf("expected span_id in log output, got: %s", output)
	}

	wantTraceID := sc.TraceID().String()
	if !bytes.Contains([]byte(output), []byte(wantTraceID)) {
		t.Errorf("expected trace_id=%s in log output, got: %s", wantTraceID, output)
	}

	wantSpanID := sc.SpanID().String()
	if !bytes.Contains([]byte(output), []byte(wantSpanID)) {
		t.Errorf("expected span_id=%s in log output, got: %s", wantSpanID, output)
	}
}

func TestTraceHandler_WithoutSpan(t *testing.T) {
	var buf bytes.Buffer
	inner := slog.NewTextHandler(&buf, nil)
	logger := slog.New(tracing.NewTraceHandler(inner))

	// No span in context — trace_id / span_id should NOT appear.
	logger.Error("something failed", "error", "not found")

	output := buf.String()

	if bytes.Contains([]byte(output), []byte("trace_id=")) {
		t.Errorf("expected no trace_id in log output, got: %s", output)
	}
	if bytes.Contains([]byte(output), []byte("span_id=")) {
		t.Errorf("expected no span_id in log output, got: %s", output)
	}
}
