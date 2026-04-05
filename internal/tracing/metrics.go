package tracing

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	promexporter "go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// MeterShutdown wraps the MeterProvider so callers can shut it down gracefully.
type MeterShutdown struct {
	provider *sdkmetric.MeterProvider
}

// Shutdown stops the MeterProvider and flushes remaining metrics.
func (m *MeterShutdown) Shutdown(ctx context.Context) error {
	return m.provider.Shutdown(ctx)
}

// InitMeter sets up the OpenTelemetry MeterProvider with a Prometheus exporter.
// Unlike tracing, the Prometheus exporter is pull-based and always active
// — no external endpoint configuration is needed.
func InitMeter(serviceName string) (*MeterShutdown, error) {
	exporter, err := promexporter.New()
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
	)
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	slog.Info("metrics enabled (prometheus exporter)")
	return &MeterShutdown{provider: mp}, nil
}
