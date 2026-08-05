package queue_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	"github.com/ndollem/mpp/apps/api/internal/middleware"
	"github.com/ndollem/mpp/apps/api/internal/modules/core/api_key"
	"github.com/ndollem/mpp/apps/api/internal/modules/core/role"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/queue"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/queue/domain"
	"github.com/ndollem/mpp/apps/api/internal/testsupport"
)

// Seed ids. Services -003 (Dukcapil, prefix A) and -022 (BPJS, prefix C) carry no
// rows in seeders/mpp/007_demo_antrian.sql, so tests never collide with demo data.
const (
	tenantCompany = "a1000000-0000-0000-0000-000000000001"
	ownerUser     = "10000000-0000-0000-0000-000000000001"

	instansiA   = "a2000000-0000-0000-0000-000000000001"
	layananA    = "a3000000-0000-0000-0000-000000000003"
	instansiC   = "a2000000-0000-0000-0000-000000000003"
	layananC    = "a3000000-0000-0000-0000-000000000022"
	demoLayanan = "a3000000-0000-0000-0000-000000000002" // seeded with nomor_seq 1 today
)

func seqKey(layananID string, day time.Time) string {
	return fmt.Sprintf("queue:seq:%s:%s", layananID, day.Format("2006-01-02"))
}

// cleanupService removes every antrian row and Redis counter a test created for one
// service/day pair.
func cleanupService(t *testing.T, db *pgxpool.Pool, rdb *goredis.Client, layananID string, day time.Time) {
	t.Helper()
	t.Cleanup(func() {
		c := context.Background()
		_, _ = db.Exec(c, `DELETE FROM mpp.antrian WHERE jenis_layanan_id = $1 AND queue_date = $2`,
			layananID, day)
		rdb.Del(c, seqKey(layananID, day))
	})
	rdb.Del(context.Background(), seqKey(layananID, day))
}

func newPemohon(t *testing.T, db *pgxpool.Pool) string {
	t.Helper()
	var id string
	require.NoError(t, db.QueryRow(context.Background(),
		`INSERT INTO mpp.pemohon (name, phone) VALUES ($1, '08120000000') RETURNING id`,
		"Queue "+t.Name()).Scan(&id))
	t.Cleanup(func() {
		_, _ = db.Exec(context.Background(), `DELETE FROM mpp.pemohon WHERE id = $1`, id)
	})
	return id
}

// enqueue runs one EnqueueTx in its own committed transaction.
func enqueue(t *testing.T, db *pgxpool.Pool, svc *queue.QueueModule, p domain.EnqueueParams) *domain.Antrian {
	t.Helper()
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	require.NoError(t, err)
	defer func() { _ = tx.Rollback(ctx) }()

	antrian, err := svc.Service.EnqueueTx(ctx, tx, p)
	require.NoError(t, err)
	require.NoError(t, tx.Commit(ctx))
	return antrian
}

func params(pemohonID, instansiID, layananID string, day time.Time) domain.EnqueueParams {
	return domain.EnqueueParams{
		PemohonID:      pemohonID,
		InstansiID:     instansiID,
		JenisLayananID: layananID,
		QueueDate:      day,
		Source:         domain.SourceBooking,
	}
}

// BR-03: counters are keyed per service per day, so a fresh Redis must resume from
// what Postgres already holds instead of handing out number 1 again. Without this
// the seeded demo antrian (nomor_seq 1 today) would collide on first check-in.
func TestEnqueue_ResumesAfterNumbersAlreadyInPostgres(t *testing.T) {
	db, rdb := testsupport.Postgres(t), testsupport.Redis(t)
	day := time.Date(2099, 3, 1, 0, 0, 0, 0, time.UTC)
	cleanupService(t, db, rdb, layananA, day)
	mod := queue.Initialize(db, rdb)
	pemohon := newPemohon(t, db)

	// A number already exists for that service/day, and Redis knows nothing about it.
	_, err := db.Exec(context.Background(),
		`INSERT INTO mpp.antrian (pemohon_id, instansi_id, jenis_layanan_id, nomor, nomor_seq, queue_date)
		 VALUES ($1, $2, $3, 'A-007', 7, $4)`, pemohon, instansiA, layananA, day)
	require.NoError(t, err)

	antrian := enqueue(t, db, mod, params(pemohon, instansiA, layananA, day))

	require.Equal(t, 8, antrian.NomorSeq, "must resume after the highest number in Postgres")
	require.Equal(t, "A-008", antrian.Nomor)
}

func TestEnqueue_NumbersAreScopedPerServicePerDay(t *testing.T) {
	db, rdb := testsupport.Postgres(t), testsupport.Redis(t)
	day := time.Date(2099, 3, 2, 0, 0, 0, 0, time.UTC)
	next := day.AddDate(0, 0, 1)
	cleanupService(t, db, rdb, layananA, day)
	cleanupService(t, db, rdb, layananC, day)
	cleanupService(t, db, rdb, layananA, next)
	mod := queue.Initialize(db, rdb)
	pemohon := newPemohon(t, db)

	first := enqueue(t, db, mod, params(pemohon, instansiA, layananA, day))
	second := enqueue(t, db, mod, params(pemohon, instansiA, layananA, day))
	otherService := enqueue(t, db, mod, params(pemohon, instansiC, layananC, day))
	otherDay := enqueue(t, db, mod, params(pemohon, instansiA, layananA, next))

	require.Equal(t, "A-001", first.Nomor)
	require.Equal(t, "A-002", second.Nomor, "same service and day keeps counting")
	require.Equal(t, "C-001", otherService.Nomor, "another service starts over, with its own prefix")
	require.Equal(t, "A-001", otherDay.Nomor, "a new day starts over (BR-03)")
}

func TestEnqueue_StatusIsWaiting(t *testing.T) {
	db, rdb := testsupport.Postgres(t), testsupport.Redis(t)
	day := time.Date(2099, 3, 3, 0, 0, 0, 0, time.UTC)
	cleanupService(t, db, rdb, layananA, day)
	mod := queue.Initialize(db, rdb)
	pemohon := newPemohon(t, db)

	antrian := enqueue(t, db, mod, params(pemohon, instansiA, layananA, day))

	require.Equal(t, domain.StatusWaiting, antrian.Status)
	require.False(t, antrian.QueuedAt.IsZero())

	var status string
	require.NoError(t, db.QueryRow(context.Background(),
		`SELECT status FROM mpp.antrian WHERE id = $1`, antrian.ID).Scan(&status))
	require.Equal(t, "WAITING", status)
}

// ── GET /mpp/v1/queue ────────────────────────────────────────────────────────

func newAPIKey(t *testing.T, db *pgxpool.Pool) string {
	t.Helper()
	const keyID, secret = "queuetestkey0001", "queuetestsecretvalue001"
	sum := sha256.Sum256([]byte(secret))

	_, err := db.Exec(context.Background(),
		`INSERT INTO core.api_keys (user_id, company_id, key_id, secret_hash, key_prefix, name, environment)
		 VALUES ($1, $2, $3, $4, $5, 'queue test key', 'test')`,
		ownerUser, tenantCompany, keyID, hex.EncodeToString(sum[:]), "wiz_test_"+keyID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = db.Exec(context.Background(), `DELETE FROM core.api_keys WHERE key_id = $1`, keyID)
	})

	return fmt.Sprintf("wiz_test_%s_%s", keyID, secret)
}

func newServer(t *testing.T, db *pgxpool.Pool, rdb *goredis.Client) (*gin.Engine, *queue.QueueModule) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	roleModule := role.Initialize(db)
	middleware.SetApiKeyValidator(api_key.Initialize(db, roleModule.Repository).Service)

	mod := queue.Initialize(db, rdb)
	r := gin.New()
	mod.SetupRoutes(r.Group("/mpp/v1"))
	return r, mod
}

func getQueue(r *gin.Engine, apiKey, query string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/mpp/v1/queue"+query, nil)
	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// Shape mirrors the websocket `snapshot` payload — see dto.QueueStream.
type listEnvelope struct {
	Data struct {
		LayananID    string `json:"layanan_id"`
		WaitingCount int    `json:"waiting_count"`
		Waiting      []struct {
			AntrianID string `json:"antrian_id"`
			Nomor     string `json:"nomor"`
			Status    string `json:"status"`
		} `json:"waiting"`
	} `json:"data"`
}

func TestQueue_ListsWaitingForService(t *testing.T) {
	db, rdb := testsupport.Postgres(t), testsupport.Redis(t)
	today := time.Now().UTC().Truncate(24 * time.Hour)
	cleanupService(t, db, rdb, layananA, today)
	r, mod := newServer(t, db, rdb)
	key := newAPIKey(t, db)
	pemohon := newPemohon(t, db)

	first := enqueue(t, db, mod, params(pemohon, instansiA, layananA, today))
	second := enqueue(t, db, mod, params(pemohon, instansiA, layananA, today))

	w := getQueue(r, key, "?layanan_id="+layananA)

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	var e listEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &e))
	require.Len(t, e.Data.Waiting, 2)
	require.Equal(t, first.ID, e.Data.Waiting[0].AntrianID, "oldest first")
	require.Equal(t, second.Nomor, e.Data.Waiting[1].Nomor)
	require.Equal(t, "WAITING", e.Data.Waiting[0].Status)
	require.Equal(t, 2, e.Data.WaitingCount)
	require.Equal(t, layananA, e.Data.LayananID)
}

func TestQueue_ExcludesOtherServices(t *testing.T) {
	db, rdb := testsupport.Postgres(t), testsupport.Redis(t)
	today := time.Now().UTC().Truncate(24 * time.Hour)
	cleanupService(t, db, rdb, layananA, today)
	r, mod := newServer(t, db, rdb)
	key := newAPIKey(t, db)
	pemohon := newPemohon(t, db)

	mine := enqueue(t, db, mod, params(pemohon, instansiA, layananA, today))

	// demoLayanan belongs to the same instansi (prefix A) and carries its own
	// seeded ticket, so the number alone proves nothing — the id does.
	w := getQueue(r, key, "?layanan_id="+demoLayanan)

	require.Equal(t, http.StatusOK, w.Code)
	var e listEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &e))
	require.Equal(t, demoLayanan, e.Data.LayananID)
	for _, item := range e.Data.Waiting {
		require.NotEqual(t, mine.ID, item.AntrianID, "another service's stream must not leak in")
	}
}

func TestQueue_RequiresLayananID(t *testing.T) {
	db, rdb := testsupport.Postgres(t), testsupport.Redis(t)
	r, _ := newServer(t, db, rdb)
	key := newAPIKey(t, db)

	require.Equal(t, http.StatusBadRequest, getQueue(r, key, "").Code)
}

func TestQueue_RequiresAuth(t *testing.T) {
	db, rdb := testsupport.Postgres(t), testsupport.Redis(t)
	r, _ := newServer(t, db, rdb)

	require.Equal(t, http.StatusUnauthorized, getQueue(r, "", "?layanan_id="+layananA).Code)
}
