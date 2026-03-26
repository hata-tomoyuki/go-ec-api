package main

import (
	"context"
	"log/slog"
	"os"

	"example.com/ecommerce/internal/env"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=127.0.0.1 port=5433 user=postgres password=postgres dbname=ecom sslmode=disable"),
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("connected to database", "dsn", cfg.db.dsn)

	api := application{
		config: cfg,
		db:     conn,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
