package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	// サーバーをゴルーチンで起動し、エラーをチャネルで受け取る
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("starting server", "addr", app.config.addr)
		serverErr <- srv.ListenAndServe()
	}()

	// OS シグナル（Ctrl+C / kill）を待ち受けるチャネル
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// シグナルまたはサーバーエラーを待ち受けて、グレースフルシャットダウンを実行する
	select {
	case err := <-serverErr:
		return err
	case <-quit:
		slog.Info("shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		return srv.Shutdown(ctx)
	}
}

type application struct {
	config config
	db     *pgxpool.Pool
	rdb    *redis.Client
}

type config struct {
	addr  string
	db    dbConfig
	redis redisConfig
}

type dbConfig struct {
	dsn string
}

type redisConfig struct {
	addr string
}
