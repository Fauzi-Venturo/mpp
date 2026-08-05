package serving_test

import (
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
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/serving"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/serving/domain"
	"github.com/ndollem/mpp/apps/api/internal/testsupport"
)

// Seed ids — seeders/mpp/003_instansi.sql, 004_layanan.sql, 005_loket.sql.
// Service -003 carries no seeded antrian, so tests never fight the demo rows.
const (
	tenantCompany = "a1000000-0000-0000-0000-000000000001"
	ownerUser     = "10000000-0000-0000-0000-000000000001"

	instansiA = "a2000000-0000-0000-0000-000000000001" // Dukcapil, prefix A
	layananA  = "a3000000-0000-0000-0000-000000000003"
	loket3    = "a5000000-0000-0000-0000-000000000003" // "Loket 3"
)

func newAPIKey(t *testing.T, db *pgxpool.Pool) string {
	t.Helper()
	const keyID, secret = "servingtestkey01", "servingtestsecretvalue01"
	sum := sha256.Sum256([]byte(secret))

	_, err := db.Exec(context.Background(),
		`INSERT INTO core.api_keys (user_id, company_id, key_id, secret_hash, key_prefix, name, environment)
		 VALUES ($1, $2, $3, $4, $5, 'serving test key', 'test')`,
		ownerUser, tenantCompany, keyID, hex.EncodeToString(sum[:]), "wiz_test_"+keyID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = db.Exec(context.Background(), `DELETE FROM core.api_keys WHERE key_id = $1`, keyID)
	})

	return fmt.Sprintf("wiz_test_%s_%s", keyID, secret)
}

func newServer(t *testing.T, db *pgxpool.Pool) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	roleModule := role.Initialize(db)
	middleware.SetApiKeyValidator(api_key.Initialize(db, roleModule.Repository).Service)

	r := gin.New()
	serving.Initialize(db).SetupRoutes(r.Group("/mpp/v1"))
	return r
}

// ticket is an antrian planted directly through SQL. Slice 5 owns the call/start
// endpoints that would normally produce CALLED/SERVING rows, so these tests build
// that state themselves.
type ticket struct {
	db    *pgxpool.Pool
	id    string
	nomor string
}

func newTicket(t *testing.T, db *pgxpool.Pool, status string, seq int, loketID *string, servedAt *time.Time) ticket {
	t.Helper()
	ctx := context.Background()

	var pemohonID string
	require.NoError(t, db.QueryRow(ctx,
		`INSERT INTO mpp.pemohon (name) VALUES ($1) RETURNING id`, "Serving "+t.Name()).Scan(&pemohonID))

	nomor := fmt.Sprintf("A-%03d", seq)
	var id string
	require.NoError(t, db.QueryRow(ctx,
		`INSERT INTO mpp.antrian (pemohon_id, instansi_id, jenis_layanan_id, nomor, nomor_seq,
		                          queue_date, status, loket_id, served_at, call_count)
		 VALUES ($1, $2, $3, $4, $5, CURRENT_DATE, $6, $7, $8, 1) RETURNING id`,
		pemohonID, instansiA, layananA, nomor, seq, status, loketID, servedAt).Scan(&id))

	t.Cleanup(func() {
		c := context.Background()
		_, _ = db.Exec(c, `DELETE FROM mpp.serving_session WHERE antrian_id = $1`, id)
		_, _ = db.Exec(c, `DELETE FROM mpp.antrian WHERE id = $1`, id)
		_, _ = db.Exec(c, `DELETE FROM mpp.pemohon WHERE id = $1`, pemohonID)
	})

	return ticket{db: db, id: id, nomor: nomor}
}

func (tk ticket) openSession(t *testing.T) {
	t.Helper()
	_, err := tk.db.Exec(context.Background(),
		`INSERT INTO mpp.serving_session (antrian_id, loket_id, started_at) VALUES ($1, $2, NOW())`,
		tk.id, loket3)
	require.NoError(t, err)
}

func (tk ticket) status(t *testing.T) (status string, doneAt *time.Time) {
	t.Helper()
	require.NoError(t, tk.db.QueryRow(context.Background(),
		`SELECT status, done_at FROM mpp.antrian WHERE id = $1`, tk.id).Scan(&status, &doneAt))
	return status, doneAt
}

func (tk ticket) session(t *testing.T) (endedAt *time.Time, outcome *string) {
	t.Helper()
	require.NoError(t, tk.db.QueryRow(context.Background(),
		`SELECT ended_at, outcome FROM mpp.serving_session WHERE antrian_id = $1`,
		tk.id).Scan(&endedAt, &outcome))
	return endedAt, outcome
}

func loketIdleAt(t *testing.T, db *pgxpool.Pool) *time.Time {
	t.Helper()
	var at *time.Time
	require.NoError(t, db.QueryRow(context.Background(),
		`SELECT last_idle_at FROM mpp.loket WHERE id = $1`, loket3).Scan(&at))
	return at
}

func done(r *gin.Engine, apiKey, antrianID string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/mpp/v1/antrian/"+antrianID+"/done", nil)
	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func display(r *gin.Engine, apiKey, query string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/mpp/v1/display"+query, nil)
	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

type doneEnvelope struct {
	Data struct {
		AntrianID       string `json:"antrian_id"`
		Status          string `json:"status"`
		DurationSeconds *int   `json:"duration_seconds"`
	} `json:"data"`
}

type displayEnvelope struct {
	Data struct {
		InstansiID string `json:"instansi_id"`
		Current    *struct {
			AntrianID string `json:"antrian_id"`
			Nomor     string `json:"nomor"`
			Loket     string `json:"loket"`
			Status    string `json:"status"`
			TTSText   string `json:"tts_text"`
		} `json:"current"`
		NextUp []struct {
			AntrianID string `json:"antrian_id"`
			Nomor     string `json:"nomor"`
		} `json:"next_up"`
	} `json:"data"`
}

// ── POST /mpp/v1/antrian/{id}/done ───────────────────────────────────────────

func TestDone_FinishesAndClosesTheSession(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	key := newAPIKey(t, db)
	loket := loket3
	served := time.Now().UTC().Add(-2 * time.Minute)
	tk := newTicket(t, db, "SERVING", 801, &loket, &served)
	tk.openSession(t)
	idleBefore := loketIdleAt(t, db)

	w := done(r, key, tk.id)

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	var e doneEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &e))
	require.Equal(t, "DONE", e.Data.Status)

	status, doneAt := tk.status(t)
	require.Equal(t, "DONE", status)
	require.NotNil(t, doneAt, "done_at must be stamped")

	endedAt, outcome := tk.session(t)
	require.NotNil(t, endedAt, "the serving session must be closed (BR-19)")
	require.NotNil(t, outcome)
	require.Equal(t, "DONE", *outcome)

	// last_idle_at feeds slice 5's idle-longest allocation (BR-12).
	idleAfter := loketIdleAt(t, db)
	require.NotNil(t, idleAfter)
	if idleBefore != nil {
		require.True(t, idleAfter.After(*idleBefore), "the loket must be marked idle again")
	}
}

func TestDone_RecordsServiceDuration(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	key := newAPIKey(t, db)
	loket := loket3
	served := time.Now().UTC().Add(-5 * time.Minute)
	tk := newTicket(t, db, "SERVING", 802, &loket, &served)
	tk.openSession(t)

	w := done(r, key, tk.id)

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	var e doneEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &e))
	require.NotNil(t, e.Data.DurationSeconds, "BR-19 wants the service duration recorded")
	require.InDelta(t, 300, *e.Data.DurationSeconds, 10)
}

func TestDone_RejectsWhenNotServing(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	key := newAPIKey(t, db)
	tk := newTicket(t, db, "WAITING", 803, nil, nil)

	w := done(r, key, tk.id)

	require.Equal(t, http.StatusConflict, w.Code, "body: %s", w.Body.String())
	status, doneAt := tk.status(t)
	require.Equal(t, "WAITING", status, "DONE is only reachable from SERVING")
	require.Nil(t, doneAt)
}

// A read-then-write guard passes every sequential test and still lets two
// operators close the same ticket twice.
func TestDone_ConcurrentFinishesOnce(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	key := newAPIKey(t, db)
	loket := loket3
	served := time.Now().UTC().Add(-time.Minute)
	tk := newTicket(t, db, "SERVING", 804, &loket, &served)
	tk.openSession(t)

	const attempts = 10
	codes := make([]int, attempts)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := range attempts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			codes[i] = done(r, key, tk.id).Code
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
	require.Equal(t, 1, ok, "only one request may finish the ticket")
	require.Equal(t, attempts-1, conflict)
}

func TestDone_UnknownAntrian(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	key := newAPIKey(t, db)

	w := done(r, key, "00000000-0000-0000-0000-000000000000")

	require.Equal(t, http.StatusNotFound, w.Code, "body: %s", w.Body.String())
}

// Slice 5 owns /start, which opens the session. Until it lands — and whenever an
// operator finishes a ticket that never got a session row — /done must still work.
func TestDone_WithoutOpenSessionStillFinishes(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	key := newAPIKey(t, db)
	loket := loket3
	served := time.Now().UTC().Add(-time.Minute)
	tk := newTicket(t, db, "SERVING", 805, &loket, &served)

	w := done(r, key, tk.id)

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	status, _ := tk.status(t)
	require.Equal(t, "DONE", status)
}

func TestDone_RequiresAuth(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	loket := loket3
	tk := newTicket(t, db, "SERVING", 806, &loket, nil)

	require.Equal(t, http.StatusUnauthorized, done(r, "", tk.id).Code)
	status, _ := tk.status(t)
	require.Equal(t, "SERVING", status)
}

// ── GET /mpp/v1/display ──────────────────────────────────────────────────────

func TestDisplay_ShowsCurrentCallAndNextUp(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	key := newAPIKey(t, db)
	loket := loket3
	called := newTicket(t, db, "CALLED", 810, &loket, nil)
	waiting := newTicket(t, db, "WAITING", 811, nil, nil)

	w := display(r, key, "?instansi_id="+instansiA)

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	var e displayEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &e))
	require.NotNil(t, e.Data.Current, "the TV must show the number being called")
	require.Equal(t, called.id, e.Data.Current.AntrianID)
	require.Equal(t, called.nomor, e.Data.Current.Nomor)
	require.Equal(t, "Loket 3", e.Data.Current.Loket, "the TV shows a label, not an id")
	require.NotEmpty(t, e.Data.Current.TTSText)

	var found bool
	for _, n := range e.Data.NextUp {
		if n.AntrianID == waiting.id {
			found = true
		}
	}
	require.True(t, found, "waiting tickets must appear in next_up")
}

// Nothing is CALLED or SERVING in the seed data, so this is the default path a
// real TV hits, not an edge case: an idle screen, never an error.
func TestDisplay_EmptyWhenNothingIsCalled(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	key := newAPIKey(t, db)
	// Imigrasi has loket but no seeded antrian at all.
	const instansiB = "a2000000-0000-0000-0000-000000000002"

	w := display(r, key, "?instansi_id="+instansiB)

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	var e displayEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &e))
	require.Nil(t, e.Data.Current)
	require.NotNil(t, e.Data.NextUp, "next_up is an empty list, never null")
	require.Empty(t, e.Data.NextUp)
}

func TestDisplay_RequiresInstansiID(t *testing.T) {
	db := testsupport.Postgres(t)
	r := newServer(t, db)
	key := newAPIKey(t, db)

	require.Equal(t, http.StatusBadRequest, display(r, key, "").Code)
}

// ── tts_text ─────────────────────────────────────────────────────────────────

// The backend owns pronunciation so every TV says the same thing
// (docs/05-integrations/tts-voice-calling.md:54-59).
func TestBuildTTSText(t *testing.T) {
	cases := []struct {
		nomor, loket, want string
	}{
		{"A-014", "Loket 3", "Nomor antrian A - nol satu empat, silakan menuju loket tiga"},
		{"C-001", "Loket 1", "Nomor antrian C - nol nol satu, silakan menuju loket satu"},
		{"B-205", "Loket 2", "Nomor antrian B - dua nol lima, silakan menuju loket dua"},
	}

	for _, c := range cases {
		require.Equal(t, c.want, domain.BuildTTSText(c.nomor, c.loket), "nomor %s", c.nomor)
	}
}
