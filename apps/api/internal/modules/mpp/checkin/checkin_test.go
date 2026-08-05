package checkin_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
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

	"github.com/ndollem/mpp/apps/api/internal/middleware"
	"github.com/ndollem/mpp/apps/api/internal/modules/core/api_key"
	"github.com/ndollem/mpp/apps/api/internal/modules/core/role"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/checkin"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/queue"
	"github.com/ndollem/mpp/apps/api/internal/testsupport"
)

// Seeded by seeders/mpp/002_tenant.sql and seeders/core/002_users.sql.
const (
	tenantCompany = "a1000000-0000-0000-0000-000000000001"
	ownerUser     = "10000000-0000-0000-0000-000000000001"
	instansiSeed  = "a2000000-0000-0000-0000-000000000001"
	layananSeed   = "a3000000-0000-0000-0000-000000000001"
)

// newServer wires what an API-key request needs plus the queue module that hands
// out numbers. No router.Setup: middleware.JWTAuth() resolves X-API-Key straight
// against Postgres.
func newServer(t *testing.T, db *pgxpool.Pool) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	roleModule := role.Initialize(db)
	middleware.SetApiKeyValidator(api_key.Initialize(db, roleModule.Repository).Service)

	checkinModule := checkin.Initialize(db)
	checkinModule.Service.SetEnqueuer(queue.Initialize(db, testsupport.Redis(t)).Service)

	r := gin.New()
	checkinModule.SetupRoutes(r.Group("/mpp/v1"))
	return r
}

// resetCounter clears the Redis sequence for today's stream so number assertions
// are repeatable, and drops the antrian rows a test enqueued.
func resetCounter(t *testing.T, db *pgxpool.Pool, layananID string) {
	t.Helper()
	rdb := testsupport.Redis(t)
	day := time.Now().UTC().Format("2006-01-02")
	key := fmt.Sprintf("queue:seq:%s:%s", layananID, day)

	drop := func() {
		c := context.Background()
		_, _ = db.Exec(c, `DELETE FROM mpp.antrian WHERE jenis_layanan_id = $1 AND queue_date = CURRENT_DATE`, layananID)
		rdb.Del(c, key)
	}
	drop()
	t.Cleanup(drop)
}

// newAPIKey inserts a kiosk device key. Secrets are stored as a plain SHA-256 hex
// digest (api_key.service.go hashSecret), so the fixture can build one directly.
func newAPIKey(t *testing.T, db *pgxpool.Pool) string {
	t.Helper()
	const (
		keyID  = "kiosktestkey0001"
		secret = "kiosktestsecretvalue0001"
	)
	sum := sha256.Sum256([]byte(secret))

	_, err := db.Exec(context.Background(),
		`INSERT INTO core.api_keys (user_id, company_id, key_id, secret_hash, key_prefix, name, environment)
		 VALUES ($1, $2, $3, $4, $5, 'kiosk test key', 'test')`,
		ownerUser, tenantCompany, keyID, hex.EncodeToString(sum[:]), "wiz_test_"+keyID)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = db.Exec(context.Background(), `DELETE FROM core.api_keys WHERE key_id = $1`, keyID)
	})

	return fmt.Sprintf("wiz_test_%s_%s", keyID, secret)
}

// booked is a BOOKED booking carrying a QR token, written straight through SQL so
// these tests do not depend on the booking module's in-flight slice-2 changes.
type booked struct {
	db        *pgxpool.Pool
	id        string
	token     string
	pemohonID string
}

func newBooking(t *testing.T, db *pgxpool.Pool, tanggal string, expiresAt time.Time) booked {
	t.Helper()
	ctx := context.Background()

	token := fmt.Sprintf("%064x", time.Now().UnixNano()+int64(len(t.Name())))

	var pemohonID string
	require.NoError(t, db.QueryRow(ctx,
		`INSERT INTO mpp.pemohon (name, phone) VALUES ($1, '08123456789') RETURNING id`,
		"Checkin "+t.Name()).Scan(&pemohonID))

	var id string
	require.NoError(t, db.QueryRow(ctx,
		`INSERT INTO mpp.booking (pemohon_id, instansi_id, jenis_layanan_id, tanggal, channel, status,
		                          qr_token, qr_expires_at)
		 VALUES ($1, $2, $3, $4, 'WEB', 'BOOKED', $5, $6) RETURNING id`,
		pemohonID, instansiSeed, layananSeed, tanggal, token, expiresAt).Scan(&id))

	t.Cleanup(func() {
		c := context.Background()
		_, _ = db.Exec(c, `DELETE FROM mpp.antrian WHERE pemohon_id = $1`, pemohonID)
		_, _ = db.Exec(c, `DELETE FROM mpp.booking WHERE id = $1`, id)
		_, _ = db.Exec(c, `DELETE FROM mpp.pemohon WHERE id = $1`, pemohonID)
	})

	return booked{db: db, id: id, token: token, pemohonID: pemohonID}
}

func (b booked) state(t *testing.T) (status string, checkedInAt *time.Time) {
	t.Helper()
	require.NoError(t, b.db.QueryRow(context.Background(),
		`SELECT status, checked_in_at FROM mpp.booking WHERE id = $1`, b.id).Scan(&status, &checkedInAt))
	return status, checkedInAt
}

func scan(r *gin.Engine, apiKey, token string) *httptest.ResponseRecorder {
	body := fmt.Sprintf(`{"qr_token":%q}`, token)
	req := httptest.NewRequest(http.MethodPost, "/mpp/v1/checkin", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

type envelope struct {
	Data struct {
		BookingID   string     `json:"booking_id"`
		Status      string     `json:"status"`
		CheckedInAt *time.Time `json:"checked_in_at"`
		Antrian     *struct {
			ID       string `json:"id"`
			Nomor    string `json:"nomor"`
			NomorSeq int    `json:"nomor_seq"`
			Status   string `json:"status"`
		} `json:"antrian"`
	} `json:"data"`
	Message string `json:"message"`
}

func today() string { return time.Now().UTC().Format("2006-01-02") }

func endOfToday() time.Time {
	n := time.Now().UTC()
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC).Add(24 * time.Hour)
}

func TestCheckin_Success(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	key := newAPIKey(t, db)
	b := newBooking(t, db, today(), endOfToday())

	w := scan(r, key, b.token)

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	var e envelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &e))
	require.Equal(t, b.id, e.Data.BookingID)
	require.Equal(t, "CHECKED_IN", e.Data.Status)

	status, checkedInAt := b.state(t)
	require.Equal(t, "CHECKED_IN", status)
	require.NotNil(t, checkedInAt, "checked_in_at must be stamped")
}

func TestCheckin_TokenReuseRejected(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	key := newAPIKey(t, db)
	b := newBooking(t, db, today(), endOfToday())

	require.Equal(t, http.StatusOK, scan(r, key, b.token).Code)
	_, firstCheckIn := b.state(t)

	second := scan(r, key, b.token)

	require.Equal(t, http.StatusConflict, second.Code, "body: %s", second.Body.String())
	status, secondCheckIn := b.state(t)
	require.Equal(t, "CHECKED_IN", status)
	require.Equal(t, firstCheckIn.UnixNano(), secondCheckIn.UnixNano(),
		"a replayed scan must not restamp checked_in_at")
}

func TestCheckin_ExpiredTokenRejected(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	key := newAPIKey(t, db)
	// Booked yesterday, window already closed.
	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	b := newBooking(t, db, yesterday.Format("2006-01-02"), yesterday)

	w := scan(r, key, b.token)

	require.Equal(t, http.StatusConflict, w.Code, "body: %s", w.Body.String())
	status, checkedInAt := b.state(t)
	require.Equal(t, "BOOKED", status, "an expired scan must not transition the booking")
	require.Nil(t, checkedInAt)
}

// A reuse guard that reads-then-writes still passes the sequential reuse test.
// Only a concurrent double scan proves the transition is single-use.
func TestCheckin_ConcurrentDoubleScanChecksInOnce(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	key := newAPIKey(t, db)
	b := newBooking(t, db, today(), endOfToday())

	const attempts = 10
	codes := make([]int, attempts)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := range attempts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			codes[i] = scan(r, key, b.token).Code
		}()
	}
	close(start)
	wg.Wait()

	ok, conflict := 0, 0
	for _, c := range codes {
		switch c {
		case http.StatusOK:
			ok++
		case http.StatusConflict:
			conflict++
		default:
			t.Fatalf("unexpected status %d", c)
		}
	}

	require.Equal(t, 1, ok, "exactly one scan may check the booking in")
	require.Equal(t, attempts-1, conflict)
}

func TestCheckin_UnknownToken(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	key := newAPIKey(t, db)

	w := scan(r, key, fmt.Sprintf("%064d", 0))

	require.Equal(t, http.StatusNotFound, w.Code, "body: %s", w.Body.String())
}

// Slice 4 — check-in also enqueues: CHECKED_IN booking, WAITING antrian with a number.

func (b booked) antrianCount(t *testing.T) int {
	t.Helper()
	var n int
	require.NoError(t, b.db.QueryRow(context.Background(),
		`SELECT count(*) FROM mpp.antrian WHERE booking_id = $1`, b.id).Scan(&n))
	return n
}

func TestCheckin_AllocatesQueueNumber(t *testing.T) {
	db := testsupport.Postgres(t)
	resetCounter(t, db, layananSeed)
	r := newServer(t, db)
	key := newAPIKey(t, db)
	b := newBooking(t, db, today(), endOfToday())

	w := scan(r, key, b.token)

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	var e envelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &e))
	require.NotNil(t, e.Data.Antrian, "check-in must hand back a queue number")
	require.Regexp(t, `^A-\d{3}$`, e.Data.Antrian.Nomor, "prefix comes from mpp.instansi.prefix (BR-01)")
	require.Equal(t, "WAITING", e.Data.Antrian.Status)
	require.Equal(t, 1, b.antrianCount(t))
}

// Read-then-write numbering passes every sequential test and still hands two
// citizens the same number under load. Only parallel scans catch it.
func TestCheckin_ConcurrentScansNeverDuplicateNumbers(t *testing.T) {
	db := testsupport.Postgres(t)
	resetCounter(t, db, layananSeed)
	r := newServer(t, db)
	key := newAPIKey(t, db)

	const scans = 10
	bookings := make([]booked, scans)
	for i := range bookings {
		bookings[i] = newBooking(t, db, today(), endOfToday())
	}

	numbers := make([]string, scans)
	codes := make([]int, scans)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := range scans {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			w := scan(r, key, bookings[i].token)
			codes[i] = w.Code
			var e envelope
			if err := json.Unmarshal(w.Body.Bytes(), &e); err == nil && e.Data.Antrian != nil {
				numbers[i] = e.Data.Antrian.Nomor
			}
		}()
	}
	close(start)
	wg.Wait()

	seen := map[string]bool{}
	for i, code := range codes {
		require.Equal(t, http.StatusOK, code, "scan %d failed", i)
		require.NotEmpty(t, numbers[i])
		require.False(t, seen[numbers[i]], "number %s handed out twice", numbers[i])
		seen[numbers[i]] = true
	}
	require.Len(t, seen, scans)

	var rows int
	require.NoError(t, db.QueryRow(context.Background(),
		`SELECT count(*) FROM mpp.antrian WHERE jenis_layanan_id = $1 AND queue_date = CURRENT_DATE`,
		layananSeed).Scan(&rows))
	require.Equal(t, scans, rows)
}

// Check-in writes twice now (booking + antrian). If the second write fails the
// first must not stand, or the citizen is checked in with no number and a rescan
// only earns a 409.
func TestCheckin_FailedEnqueueRollsBackTheCheckin(t *testing.T) {
	db := testsupport.Postgres(t)
	resetCounter(t, db, layananSeed)
	r := newServer(t, db)
	key := newAPIKey(t, db)
	b := newBooking(t, db, today(), endOfToday())

	// Redis drifted behind Postgres: the next INCR yields a number already taken,
	// so the antrian INSERT trips antrian_service_day_seq_key.
	_, err := db.Exec(context.Background(),
		`INSERT INTO mpp.antrian (pemohon_id, instansi_id, jenis_layanan_id, nomor, nomor_seq, queue_date)
		 VALUES ($1, $2, $3, 'A-005', 5, CURRENT_DATE)`, b.pemohonID, instansiSeed, layananSeed)
	require.NoError(t, err)
	testsupport.Redis(t).Set(context.Background(),
		fmt.Sprintf("queue:seq:%s:%s", layananSeed, today()), 4, time.Hour)

	w := scan(r, key, b.token)

	require.Equal(t, http.StatusInternalServerError, w.Code, "body: %s", w.Body.String())
	status, checkedInAt := b.state(t)
	require.Equal(t, "BOOKED", status, "the check-in must roll back with the enqueue")
	require.Nil(t, checkedInAt)
	require.Equal(t, 0, b.antrianCount(t))
}

func TestCheckin_RequiresApiKey(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	b := newBooking(t, db, today(), endOfToday())

	w := scan(r, "", b.token)

	require.Equal(t, http.StatusUnauthorized, w.Code, "body: %s", w.Body.String())
	status, _ := b.state(t)
	require.Equal(t, "BOOKED", status, "an unauthenticated scan must change nothing")
}
