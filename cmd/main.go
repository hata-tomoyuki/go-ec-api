package main

import (
	"context"
	"log/slog"
	"os"

	"example.com/ecommerce/internal/env"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("your-secret-key"), nil)
}

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
