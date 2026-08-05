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
	"time"

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
		ID          string    `json:"id"`
		Status      string    `json:"status"`
		QRToken     string    `json:"qr_token"`
		QRExpiresAt time.Time `json:"qr_expires_at"`
	} `json:"data"`
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

func get(r *gin.Engine, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	return w
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

// Slice 2 — BR-09: booking must hand out a single-use, time-bound check-in token.

func (f fixture) storedToken(t *testing.T, bookingID string) string {
	t.Helper()
	var tok *string
	require.NoError(t, f.db.QueryRow(context.Background(),
		`SELECT qr_token FROM mpp.booking WHERE id = $1`, bookingID).Scan(&tok))
	require.NotNil(t, tok, "qr_token must be persisted, not only returned")
	return *tok
}

func TestCreateBooking_IssuesAQrTokenValidUntilTheEndOfTheBookingDay(t *testing.T) {
	db := testsupport.Postgres(t)
	f := newQuota(t, db, "2099-01-07", 10, 0)

	w := post(newServer(db), tenantSlug, f.body("Rina"))

	require.Equal(t, http.StatusCreated, w.Code, "body: %s", w.Body.String())
	e := decode(t, w)
	require.NotEmpty(t, e.Data.QRToken, "booking response must carry the QR token")
	require.Equal(t, e.Data.QRToken, f.storedToken(t, e.Data.ID))
	// Midnight WIB closing the booking day, i.e. 17:00 UTC the day before —
	// expiring at 00:00 UTC would keep the token alive until 07:00 local.
	require.Equal(t,
		time.Date(2099, 1, 7, 17, 0, 0, 0, time.UTC),
		e.Data.QRExpiresAt.UTC(),
		"token must expire when the booking day ends in the service zone")
}

func TestCreateBooking_GivesEveryBookingItsOwnToken(t *testing.T) {
	db := testsupport.Postgres(t)
	f := newQuota(t, db, "2099-01-08", 10, 0)
	r := newServer(db)

	first := decode(t, post(r, tenantSlug, f.body("Andi")))
	second := decode(t, post(r, tenantSlug, f.body("Bayu")))

	require.NotEmpty(t, first.Data.QRToken)
	require.NotEqual(t, first.Data.QRToken, second.Data.QRToken)
}

func TestGetBooking_ReturnsTheSameToken(t *testing.T) {
	db := testsupport.Postgres(t)
	f := newQuota(t, db, "2099-01-09", 10, 0)
	r := newServer(db)
	created := decode(t, post(r, tenantSlug, f.body("Citra")))

	w := get(r, "/mpp/v1/booking/"+created.Data.ID)

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	e := decode(t, w)
	require.Equal(t, created.Data.ID, e.Data.ID)
	require.Equal(t, created.Data.QRToken, e.Data.QRToken)
}

func TestGetBooking_UnknownIDIsNotFound(t *testing.T) {
	db := testsupport.Postgres(t)

	w := get(newServer(db), "/mpp/v1/booking/aaaaaaaa-0000-0000-0000-000000000000")

	require.Equal(t, http.StatusNotFound, w.Code, "body: %s", w.Body.String())
	// The envelope proves the handler answered — gin's own 404 is plain text, so
	// this test cannot pass merely because the route is missing.
	require.NotEmpty(t, decode(t, w).Message)
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
