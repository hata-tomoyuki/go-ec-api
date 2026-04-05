package main

import (
	"context"
	"log/slog"
	"os"

	"example.com/ecommerce/internal/env"
	"example.com/ecommerce/internal/tracing"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(env.MustGetString("JWT_SECRET")), nil)
}

func main() {
	ctx := context.Background()

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=127.0.0.1 port=5433 user=postgres password=postgres dbname=ecom sslmode=disable"),
		},
	}

	logger := slog.New(tracing.NewTraceHandler(slog.NewTextHandler(os.Stdout, nil)))
	slog.SetDefault(logger)

	// --- OpenTelemetry ---
	otlpEndpoint := env.GetString("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	shutdownTracer, err := tracing.InitTracer(ctx, "ecommerce", otlpEndpoint)
	if err != nil {
		slog.Error("failed to init tracer", "error", err)
		os.Exit(1)
	}
	defer shutdownTracer(ctx)

	// --- Database pool with query tracing ---
	poolCfg, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	poolCfg.ConnConfig.Tracer = tracing.NewPgxTracer()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	logger.Info("connected to database")

	api := application{
		config: cfg,
		db:     pool,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
