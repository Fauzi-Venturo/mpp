package booking_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/booking"
	"github.com/ndollem/mpp/apps/api/internal/testsupport"
)

// The MPP building seeded by seeders/mpp/002_tenant.sql. Its tenant slug lives on
// core.clients.slug (core.companies has no slug column of its own).
const (
	tenantSlug   = "owner"
	instansiSeed = "a2000000-0000-0000-0000-000000000001"
)

// fixture is a quota row created for one test, on a date far outside the seeded
// range so tests never fight each other or the demo data.
type fixture struct {
	db         *pgxpool.Pool
	instansiID string
	layananID  string
	tanggal    string
}

// newQuota inserts an agency-wide quota row (jenis_layanan_id IS NULL, the shape
// seeders/mpp/006_kuota.sql uses) and removes everything it touched afterwards.
func newQuota(t *testing.T, db *pgxpool.Pool, tanggal string, kuota, terpakai int) fixture {
	t.Helper()
	ctx := context.Background()

	var layananID string
	err := db.QueryRow(ctx,
		`SELECT id FROM mpp.jenis_layanan WHERE instansi_id = $1 AND deleted_at IS NULL ORDER BY name LIMIT 1`,
		instansiSeed).Scan(&layananID)
	require.NoError(t, err, "seed data missing — run `make db-setup`")

	_, err = db.Exec(ctx,
		`INSERT INTO mpp.kuota_booking (instansi_id, jenis_layanan_id, tanggal, kuota, terpakai)
		 VALUES ($1, NULL, $2, $3, $4)`,
		instansiSeed, tanggal, kuota, terpakai)
	require.NoError(t, err)

	t.Cleanup(func() {
		c := context.Background()
		_, _ = db.Exec(c, `DELETE FROM mpp.pemohon WHERE id IN (
			SELECT pemohon_id FROM mpp.booking WHERE tanggal = $1)`, tanggal)
		_, _ = db.Exec(c, `DELETE FROM mpp.booking WHERE tanggal = $1`, tanggal)
		_, _ = db.Exec(c, `DELETE FROM mpp.kuota_booking WHERE tanggal = $1`, tanggal)
	})

	return fixture{db: db, instansiID: instansiSeed, layananID: layananID, tanggal: tanggal}
}

func (f fixture) terpakai(t *testing.T) int {
	t.Helper()
	var n int
	require.NoError(t, f.db.QueryRow(context.Background(),
		`SELECT terpakai FROM mpp.kuota_booking WHERE tanggal = $1 AND jenis_layanan_id IS NULL`,
		f.tanggal).Scan(&n))
	return n
}

func (f fixture) bookingCount(t *testing.T) int {
	t.Helper()
	var n int
	require.NoError(t, f.db.QueryRow(context.Background(),
		`SELECT count(*) FROM mpp.booking WHERE tanggal = $1`, f.tanggal).Scan(&n))
	return n
}

func (f fixture) body(name string) string {
	return fmt.Sprintf(`{
		"instansi_id": %q,
		"jenis_layanan_id": %q,
		"tanggal": %q,
		"pemohon": {"name": %q, "phone": "08123456789"}
	}`, f.instansiID, f.layananID, f.tanggal, name)
}

func newServer(db *pgxpool.Pool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	booking.Initialize(db).SetupRoutes(r.Group("/mpp/v1"))
	return r
}

func post(r *gin.Engine, slug, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/mpp/v1/booking", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	if slug != "" {
		req.Header.Set("X-Company-Slug", slug)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

type envelope struct {
	Data struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"data"`
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

func decode(t *testing.T, w *httptest.ResponseRecorder) envelope {
	t.Helper()
	var e envelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &e), "body: %s", w.Body.String())
	return e
}

func TestCreateBooking_Success(t *testing.T) {
	db := testsupport.Postgres(t)
	f := newQuota(t, db, "2099-01-02", 10, 0)

	w := post(newServer(db), tenantSlug, f.body("Budi Santoso"))

	require.Equal(t, http.StatusCreated, w.Code, "body: %s", w.Body.String())
	e := decode(t, w)
	require.NotEmpty(t, e.Data.ID)
	require.Equal(t, "BOOKED", e.Data.Status)
	require.Equal(t, 1, f.bookingCount(t))
	require.Equal(t, 1, f.terpakai(t), "quota must be consumed exactly once")
}

func TestCreateBooking_QuotaFull(t *testing.T) {
	db := testsupport.Postgres(t)
	f := newQuota(t, db, "2099-01-03", 1, 1) // already exhausted

	w := post(newServer(db), tenantSlug, f.body("Siti Aminah"))

	require.Equal(t, http.StatusConflict, w.Code, "body: %s", w.Body.String())
	require.Equal(t, 0, f.bookingCount(t), "no booking may be written when the quota is full")
	require.Equal(t, 1, f.terpakai(t), "terpakai must not exceed kuota")
}

// The overbooking bug fails silently: the happy path stays green even when quota
// consumption is a read-then-write race. Only concurrency exposes it.
func TestCreateBooking_ConcurrentNeverOverbooks(t *testing.T) {
	db := testsupport.Postgres(t)
	const (
		kuota    = 5
		attempts = 20
	)
	f := newQuota(t, db, "2099-01-04", kuota, 0)
	r := newServer(db)

	codes := make([]int, attempts)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := range attempts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			codes[i] = post(r, tenantSlug, f.body(fmt.Sprintf("Pemohon %02d", i))).Code
		}()
	}
	close(start)
	wg.Wait()

	created, conflict := 0, 0
	for _, c := range codes {
		switch c {
		case http.StatusCreated:
			created++
		case http.StatusConflict:
			conflict++
		default:
			t.Fatalf("unexpected status %d", c)
		}
	}

	require.Equal(t, kuota, created, "exactly kuota requests may succeed")
	require.Equal(t, attempts-kuota, conflict)
	require.Equal(t, kuota, f.terpakai(t))
	require.Equal(t, kuota, f.bookingCount(t))
}

func TestCreateBooking_UnknownCompanySlug(t *testing.T) {
	db := testsupport.Postgres(t)
	f := newQuota(t, db, "2099-01-05", 10, 0)

	w := post(newServer(db), "tenant-yang-tidak-ada", f.body("Joko"))

	require.Equal(t, http.StatusNotFound, w.Code, "body: %s", w.Body.String())
	require.Equal(t, 0, f.terpakai(t), "a foreign tenant must not consume quota")
}

func TestCreateBooking_InvalidPayload(t *testing.T) {
	db := testsupport.Postgres(t)

	w := post(newServer(db), tenantSlug, `{
		"instansi_id": "`+instansiSeed+`",
		"jenis_layanan_id": "bukan-uuid",
		"tanggal": "2099-01-06",
		"pemohon": {"name": "Ani"}
	}`)

	require.Equal(t, http.StatusBadRequest, w.Code, "body: %s", w.Body.String())
	require.NotEmpty(t, decode(t, w).Errors)
}
