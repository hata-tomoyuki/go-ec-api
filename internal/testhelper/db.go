package testhelper

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/jackc/pgx/v5/stdlib" // goose に pgx ドライバを登録
)

// NewTestDB はテスト用 PostgreSQL コンテナを起動し、マイグレーション済みの pgxpool を返す。
// テスト終了時にコンテナを自動的に破棄する。
func NewTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("ecom_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	})

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// goose でマイグレーション適用
	if err := runMigrations(dsn); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to create pgxpool: %v", err)
	}
	t.Cleanup(func() { pool.Close() })

	return pool
}

// migrationsDir はこのファイルから migrations ディレクトリへの相対パスを解決する。
func migrationsDir() string {
	// このファイルの絶対パスを取得
	_, filename, _, _ := runtime.Caller(0)
	// testhelper/ → internal/ → ecommerce/ → migrations/
	return filepath.Join(filepath.Dir(filename), "..", "adapters", "postgresql", "migrations")
}

func runMigrations(dsn string) error {
	// goose は database/sql を使うため pgx の stdlib ドライバを使う
	db, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open db for goose: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, migrationsDir()); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}
	return nil
}
