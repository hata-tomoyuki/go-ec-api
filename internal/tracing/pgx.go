package tracing

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "pgx"

type ctxKey struct{}

// PgxTracer implements pgx.QueryTracer to add OpenTelemetry spans to all SQL queries.
type PgxTracer struct{}

func NewPgxTracer() *PgxTracer {
	return &PgxTracer{}
}

func (t *PgxTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, "db.query", trace.WithAttributes(
		attribute.String("db.statement", data.SQL),
	))
	return context.WithValue(ctx, ctxKey{}, span)
}

func (t *PgxTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	span, ok := ctx.Value(ctxKey{}).(trace.Span)
	if !ok {
		return
	}
	defer span.End()

	if data.Err != nil {
		span.RecordError(data.Err)
		span.SetStatus(codes.Error, data.Err.Error())
	}
}
