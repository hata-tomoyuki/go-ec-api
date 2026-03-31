package main

import (
	"context"
	"fmt"
	"log"

	"example.com/ecommerce/internal/env"
	"example.com/ecommerce/internal/seed"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	dsn := env.GetString(
		"GOOSE_DBSTRING",
		"host=127.0.0.1 port=5433 user=postgres password=postgres dbname=ecom sslmode=disable",
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("DB接続に失敗:", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("DB Pingに失敗:", err)
	}
	fmt.Println("DB接続成功")

	if err := seed.Run(ctx, pool); err != nil {
		log.Fatal("シード投入に失敗:", err)
	}

	fmt.Println("シードデータの投入が完了しました")
}
