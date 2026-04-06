package cache

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// Store wraps a Redis client and provides type-safe cache operations.
// All errors are logged but not returned — on Redis failure the caller
// simply falls through to the database (cache-aside pattern).
type Store struct {
	rdb *redis.Client
}

func NewStore(rdb *redis.Client) *Store {
	return &Store{rdb: rdb}
}

// Get retrieves a cached value by key and unmarshals it into T.
// Returns (zero, false) on miss or error.
func Get[T any](ctx context.Context, s *Store, key string) (T, bool) {
	var zero T
	if s == nil {
		return zero, false
	}

	data, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err != redis.Nil {
			slog.Warn("cache get error", "key", key, "error", err)
		}
		return zero, false
	}

	var val T
	if err := json.Unmarshal(data, &val); err != nil {
		slog.Warn("cache unmarshal error", "key", key, "error", err)
		return zero, false
	}

	slog.Debug("cache hit", "key", key)
	return val, true
}

// Set marshals val to JSON and stores it in Redis with the given TTL.
// Errors are logged only.
func Set[T any](ctx context.Context, s *Store, key string, val T, ttl time.Duration) {
	if s == nil {
		return
	}

	data, err := json.Marshal(val)
	if err != nil {
		slog.Warn("cache marshal error", "key", key, "error", err)
		return
	}

	if err := s.rdb.Set(ctx, key, data, ttl).Err(); err != nil {
		slog.Warn("cache set error", "key", key, "error", err)
	}
}

// Delete removes a single key from the cache.
func (s *Store) Delete(ctx context.Context, key string) {
	if s == nil {
		return
	}
	if err := s.rdb.Del(ctx, key).Err(); err != nil {
		slog.Warn("cache delete error", "key", key, "error", err)
	}
}

// InvalidateByPrefix removes all keys matching the given prefixes.
// Uses SCAN (not KEYS) to avoid blocking Redis on large keyspaces.
func (s *Store) InvalidateByPrefix(ctx context.Context, prefixes ...string) {
	if s == nil {
		return
	}
	for _, prefix := range prefixes {
		var cursor uint64
		for {
			keys, next, err := s.rdb.Scan(ctx, cursor, prefix+"*", 100).Result()
			if err != nil {
				slog.Warn("cache scan error", "prefix", prefix, "error", err)
				break
			}
			if len(keys) > 0 {
				if err := s.rdb.Del(ctx, keys...).Err(); err != nil {
					slog.Warn("cache batch delete error", "prefix", prefix, "error", err)
				}
			}
			cursor = next
			if cursor == 0 {
				break
			}
		}
	}
}
