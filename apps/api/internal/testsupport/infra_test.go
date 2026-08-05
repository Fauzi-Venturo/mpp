package testsupport_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ndollem/mpp/apps/api/internal/testsupport"
)

// Smoke test for the dev engine: migrations applied and Redis writable.
// Fails loudly if `make db-setup` was never run, skips if infra is down.
func TestPostgresHasMigratedSchemas(t *testing.T) {
	pool := testsupport.Postgres(t)
	ctx := context.Background()

	var schemas int
	err := pool.QueryRow(ctx,
		`SELECT count(*) FROM information_schema.schemata WHERE schema_name IN ('core', 'mpp')`,
	).Scan(&schemas)
	require.NoError(t, err)
	assert.Equal(t, 2, schemas, "schema core+mpp must exist — run `make db-setup`")

	var instansi int
	err = pool.QueryRow(ctx, `SELECT count(*) FROM mpp.instansi`).Scan(&instansi)
	require.NoError(t, err)
	assert.Positive(t, instansi, "seeder mpp must have inserted instansi")
}

func TestRedisReadWrite(t *testing.T) {
	client := testsupport.Redis(t)
	ctx := context.Background()

	key := "mpp:testsupport:smoke"
	require.NoError(t, client.Set(ctx, key, "ok", 0).Err())
	t.Cleanup(func() { client.Del(ctx, key) })

	got, err := client.Get(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, "ok", got)
}
