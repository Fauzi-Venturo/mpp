// Package testsupport wires integration tests to the real Postgres and Redis
// started by `make up`. Tests skip (not fail) when the infra is down, so
// `go test ./...` stays green on a machine without docker running.
package testsupport

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/ndollem/mpp/apps/api/internal/config"
)

const dialTimeout = 3 * time.Second

// Postgres returns a pool against the test database, closed when the test ends.
func Postgres(t *testing.T) *pgxpool.Pool {
	t.Helper()

	cfg := config.Load()
	pool, err := pgxpool.New(context.Background(), cfg.Database.GetDSN())
	if err != nil {
		t.Skipf("postgres config invalid: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Skipf("postgres unreachable (%v) — run `make up && make db-setup`", err)
	}

	t.Cleanup(pool.Close)
	return pool
}

// Redis returns a client against the test Redis, closed when the test ends.
func Redis(t *testing.T) *redis.Client {
	t.Helper()

	cfg := config.Load()
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		t.Skipf("redis unreachable (%v) — run `make up`", err)
	}

	t.Cleanup(func() { _ = client.Close() })
	return client
}
