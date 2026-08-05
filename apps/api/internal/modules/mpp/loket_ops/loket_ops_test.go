package loket_ops_test

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
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/loket_ops"
	"github.com/ndollem/mpp/apps/api/internal/testsupport"
)

// Seed ids. Service -012 (Imigrasi, prefix B) is untouched by 007_demo_antrian.sql
// and by the checkin/queue test packages (-001, -002, -003, -021, -022), so these
// tests never fight the other suites that run in parallel. It also has
// requires_fo_verification = FALSE, keeping the FO rule (BR-24) out of the picture.
const (
	tenantCompany = "a1000000-0000-0000-0000-000000000001"
	ownerUser     = "10000000-0000-0000-0000-000000000001"

	instansiB = "a2000000-0000-0000-0000-000000000002"
	layananB  = "a3000000-0000-0000-0000-000000000012"
	loketB1   = "a5000000-0000-0000-0000-000000000011"
	loketB2   = "a5000000-0000-0000-0000-000000000012"
	loketA1   = "a5000000-0000-0000-0000-000000000001" // Dukcapil — never serves layananB
)

func today() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

// resetService wipes every antrian this suite may have written, before and after.
func resetService(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	drop := func() {
		c := context.Background()
		_, _ = db.Exec(c, `DELETE FROM mpp.serving_session WHERE antrian_id IN (
			SELECT id FROM mpp.antrian WHERE jenis_layanan_id = $1)`, layananB)
		_, _ = db.Exec(c, `DELETE FROM mpp.antrian WHERE jenis_layanan_id = $1`, layananB)
	}
	drop()
	t.Cleanup(drop)
}

// keepLoket snapshots the shared seed state a test is about to mutate. last_idle_at
// and status belong to the seeded lokets, not to this test, so they must go back.
func keepLoket(t *testing.T, db *pgxpool.Pool, loketIDs ...string) {
	t.Helper()
	type snap struct {
		status string
		idle   time.Time
	}
	saved := make(map[string]snap, len(loketIDs))

	for _, id := range loketIDs {
		var s snap
		require.NoError(t, db.QueryRow(context.Background(),
			`SELECT status, last_idle_at FROM mpp.loket WHERE id = $1`, id).Scan(&s.status, &s.idle),
			"seed data missing — run `make db-setup`")
		saved[id] = s
	}

	t.Cleanup(func() {
		for id, s := range saved {
			_, _ = db.Exec(context.Background(),
				`UPDATE mpp.loket SET status = $2, last_idle_at = $3 WHERE id = $1`, id, s.status, s.idle)
		}
	})
}

func setIdle(t *testing.T, db *pgxpool.Pool, loketID string, at time.Time) {
	t.Helper()
	_, err := db.Exec(context.Background(),
		`UPDATE mpp.loket SET last_idle_at = $2 WHERE id = $1`, loketID, at)
	require.NoError(t, err)
}

func idleOf(t *testing.T, db *pgxpool.Pool, loketID string) time.Time {
	t.Helper()
	var at time.Time
	require.NoError(t, db.QueryRow(context.Background(),
		`SELECT last_idle_at FROM mpp.loket WHERE id = $1`, loketID).Scan(&at))
	return at
}

func newPemohon(t *testing.T, db *pgxpool.Pool) string {
	t.Helper()
	var id string
	require.NoError(t, db.QueryRow(context.Background(),
		`INSERT INTO mpp.pemohon (name, phone) VALUES ($1, '08120000000') RETURNING id`,
		"Loket "+t.Name()).Scan(&id))
	t.Cleanup(func() {
		_, _ = db.Exec(context.Background(), `DELETE FROM mpp.pemohon WHERE id = $1`, id)
	})
	return id
}

// waiting seeds one WAITING ticket for today and returns its id.
func waiting(t *testing.T, db *pgxpool.Pool, seq int) string {
	t.Helper()
	var id string
	require.NoError(t, db.QueryRow(context.Background(),
		`INSERT INTO mpp.antrian (pemohon_id, instansi_id, jenis_layanan_id, nomor, nomor_seq, queue_date)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		newPemohon(t, db), instansiB, layananB, fmt.Sprintf("B-%03d", seq), seq, today()).Scan(&id))
	return id
}

type row struct {
	Status    string
	LoketID   *string
	CallCount int
	CalledAt  *time.Time
}

func load(t *testing.T, db *pgxpool.Pool, antrianID string) row {
	t.Helper()
	var r row
	require.NoError(t, db.QueryRow(context.Background(),
		`SELECT status, loket_id, call_count, called_at FROM mpp.antrian WHERE id = $1`,
		antrianID).Scan(&r.Status, &r.LoketID, &r.CallCount, &r.CalledAt))
	return r
}

// ── HTTP plumbing (same recipe as queue_test.go / checkin_test.go) ────────────

func newAPIKey(t *testing.T, db *pgxpool.Pool) string {
	t.Helper()
	// key_id is VARCHAR(16); the format parser also forbids underscores inside parts.
	const keyID, secret = "loketopskey00001", "loketopstestsecret001"
	sum := sha256.Sum256([]byte(secret))

	_, err := db.Exec(context.Background(),
		`INSERT INTO core.api_keys (user_id, company_id, key_id, secret_hash, key_prefix, name, environment)
		 VALUES ($1, $2, $3, $4, $5, 'loket ops test key', 'test')`,
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

	middleware.SetApiKeyValidator(api_key.Initialize(db, role.Initialize(db).Repository).Service)

	r := gin.New()
	loket_ops.Initialize(db).SetupRoutes(r.Group("/mpp/v1"))
	return r
}

func post(r *gin.Engine, apiKey, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func callNext(r *gin.Engine, apiKey, loketID string) *httptest.ResponseRecorder {
	body := fmt.Sprintf(`{"jenis_layanan_id": %q}`, layananB)
	if loketID != "" {
		body = fmt.Sprintf(`{"jenis_layanan_id": %q, "loket_id": %q}`, layananB, loketID)
	}
	return post(r, apiKey, "/mpp/v1/queue/next", body)
}

// Field names mirror the websocket `call.created` payload (websocket-events.md:35).
type envelope struct {
	Data struct {
		AntrianID string `json:"antrian_id"`
		Nomor     string `json:"nomor"`
		Status    string `json:"status"`
		LoketID   string `json:"loket_id"`
		CallCount int    `json:"call_count"`
	} `json:"data"`
	Message string `json:"message"`
}

func decode(t *testing.T, w *httptest.ResponseRecorder) envelope {
	t.Helper()
	var e envelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &e), "body: %s", w.Body.String())
	return e
}

// ── Call next ────────────────────────────────────────────────────────────────

// queue-state-machine.md:68 — the first call sets call_count to 1, not 0.
func TestCallNext_CallsTheOldestTicketForTheRequestedLoket(t *testing.T) {
	db := testsupport.Postgres(t)
	resetService(t, db)
	keepLoket(t, db, loketB1)
	r, key := newServer(t, db), newAPIKey(t, db)
	first, second := waiting(t, db, 1), waiting(t, db, 2)

	w := callNext(r, key, loketB1)

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	e := decode(t, w)
	require.Equal(t, first, e.Data.AntrianID, "FIFO: the oldest ticket goes first")
	require.Equal(t, "CALLED", e.Data.Status)
	require.Equal(t, 1, e.Data.CallCount)

	got := load(t, db, first)
	require.Equal(t, loketB1, *got.LoketID)
	require.NotNil(t, got.CalledAt)
	require.Equal(t, "WAITING", load(t, db, second).Status, "only one ticket may be called")
}

// BR-12 / FR-QUE-02: without an explicit loket the server picks the eligible loket
// that has been idle the longest.
func TestCallNext_WithoutLoketPicksTheLongestIdleOne(t *testing.T) {
	db := testsupport.Postgres(t)
	resetService(t, db)
	keepLoket(t, db, loketB1, loketB2)
	r, key := newServer(t, db), newAPIKey(t, db)
	ticket := waiting(t, db, 1)

	setIdle(t, db, loketB1, time.Now().UTC().Add(-2*time.Hour)) // idle longest
	setIdle(t, db, loketB2, time.Now().UTC())

	w := callNext(r, key, "")

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	require.Equal(t, loketB1, *load(t, db, ticket).LoketID)
}

func TestCallNext_EmptyQueueIsNotFound(t *testing.T) {
	db := testsupport.Postgres(t)
	resetService(t, db)
	keepLoket(t, db, loketB1)
	r, key := newServer(t, db), newAPIKey(t, db)

	w := callNext(r, key, loketB1)

	require.Equal(t, http.StatusNotFound, w.Code, "body: %s", w.Body.String())
	require.NotEmpty(t, decode(t, w).Message)
}

// BR-15: a CLOSED loket receives no allocation.
func TestCallNext_ClosedLoketIsRejected(t *testing.T) {
	db := testsupport.Postgres(t)
	resetService(t, db)
	keepLoket(t, db, loketB1)
	r, key := newServer(t, db), newAPIKey(t, db)
	ticket := waiting(t, db, 1)

	_, err := db.Exec(context.Background(), `UPDATE mpp.loket SET status = 'CLOSED' WHERE id = $1`, loketB1)
	require.NoError(t, err)

	w := callNext(r, key, loketB1)

	require.Equal(t, http.StatusConflict, w.Code, "body: %s", w.Body.String())
	require.Equal(t, "WAITING", load(t, db, ticket).Status)
}

// BR-12: eligibility also means the loket actually serves that service.
func TestCallNext_LoketOfAnotherAgencyIsRejected(t *testing.T) {
	db := testsupport.Postgres(t)
	resetService(t, db)
	keepLoket(t, db, loketA1)
	r, key := newServer(t, db), newAPIKey(t, db)
	ticket := waiting(t, db, 1)

	w := callNext(r, key, loketA1)

	require.Equal(t, http.StatusConflict, w.Code, "body: %s", w.Body.String())
	require.Equal(t, "WAITING", load(t, db, ticket).Status)
}

// Two operators pressing "next" at the same moment must never be handed the same
// ticket — a race that stays invisible on the happy path.
func TestCallNext_ConcurrentCallersNeverShareATicket(t *testing.T) {
	db := testsupport.Postgres(t)
	resetService(t, db)
	keepLoket(t, db, loketB1, loketB2)
	r, key := newServer(t, db), newAPIKey(t, db)
	waiting(t, db, 1) // exactly one ticket for two lokets

	codes := make([]int, 2)
	lokets := []string{loketB1, loketB2}
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := range codes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			codes[i] = callNext(r, key, lokets[i]).Code
		}()
	}
	close(start)
	wg.Wait()

	ok, empty := 0, 0
	for _, c := range codes {
		switch c {
		case http.StatusOK:
			ok++
		case http.StatusNotFound:
			empty++
		default:
			t.Fatalf("unexpected status %d", c)
		}
	}
	require.Equal(t, 1, ok, "exactly one operator may win the ticket")
	require.Equal(t, 1, empty)
}

// ── Recall ───────────────────────────────────────────────────────────────────

func TestRecall_CountsEachCall(t *testing.T) {
	db := testsupport.Postgres(t)
	resetService(t, db)
	keepLoket(t, db, loketB1)
	r, key := newServer(t, db), newAPIKey(t, db)
	ticket := waiting(t, db, 1)
	require.Equal(t, http.StatusOK, callNext(r, key, loketB1).Code)

	w := post(r, key, "/mpp/v1/antrian/"+ticket+"/recall", "")

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	require.Equal(t, 2, decode(t, w).Data.CallCount)
	require.Equal(t, "CALLED", load(t, db, ticket).Status)
}

// BR-16 / FR-OPR-03: call at most 3x. The DB CHECK would raise 23514 on a 4th call,
// so the service must refuse it cleanly with 409 first.
func TestRecall_FourthCallIsRejected(t *testing.T) {
	db := testsupport.Postgres(t)
	resetService(t, db)
	keepLoket(t, db, loketB1)
	r, key := newServer(t, db), newAPIKey(t, db)
	ticket := waiting(t, db, 1)
	require.Equal(t, http.StatusOK, callNext(r, key, loketB1).Code) // call 1
	require.Equal(t, http.StatusOK, post(r, key, "/mpp/v1/antrian/"+ticket+"/recall", "").Code)
	require.Equal(t, http.StatusOK, post(r, key, "/mpp/v1/antrian/"+ticket+"/recall", "").Code)

	w := post(r, key, "/mpp/v1/antrian/"+ticket+"/recall", "") // the 4th

	require.Equal(t, http.StatusConflict, w.Code, "body: %s", w.Body.String())
	require.Equal(t, 3, load(t, db, ticket).CallCount, "call_count must stop at 3")
}

// ── Start ────────────────────────────────────────────────────────────────────

func TestStart_MovesToServingAndOpensASession(t *testing.T) {
	db := testsupport.Postgres(t)
	resetService(t, db)
	keepLoket(t, db, loketB1)
	r, key := newServer(t, db), newAPIKey(t, db)
	ticket := waiting(t, db, 1)
	require.Equal(t, http.StatusOK, callNext(r, key, loketB1).Code)

	w := post(r, key, "/mpp/v1/antrian/"+ticket+"/start", "")

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	require.Equal(t, "SERVING", load(t, db, ticket).Status)

	var sessions int
	require.NoError(t, db.QueryRow(context.Background(),
		`SELECT count(*) FROM mpp.serving_session WHERE antrian_id = $1 AND ended_at IS NULL`,
		ticket).Scan(&sessions))
	require.Equal(t, 1, sessions, "serving must open exactly one session")
}

// NFR-DATA-03: illegal transitions are refused server-side with 409.
func TestStart_FromWaitingIsRejected(t *testing.T) {
	db := testsupport.Postgres(t)
	resetService(t, db)
	r, key := newServer(t, db), newAPIKey(t, db)
	ticket := waiting(t, db, 1) // never called

	w := post(r, key, "/mpp/v1/antrian/"+ticket+"/start", "")

	require.Equal(t, http.StatusConflict, w.Code, "body: %s", w.Body.String())
	require.Equal(t, "WAITING", load(t, db, ticket).Status)
}

// ── Skip ─────────────────────────────────────────────────────────────────────

// queue-state-machine.md:71,83 — a no-show closes the ticket and frees the loket,
// which means refreshing last_idle_at so the fair allocator keeps working.
func TestSkip_MarksNoShowAndFreesTheLoket(t *testing.T) {
	db := testsupport.Postgres(t)
	resetService(t, db)
	keepLoket(t, db, loketB1)
	r, key := newServer(t, db), newAPIKey(t, db)
	ticket := waiting(t, db, 1)
	require.Equal(t, http.StatusOK, callNext(r, key, loketB1).Code)

	setIdle(t, db, loketB1, time.Now().UTC().Add(-time.Hour))
	before := idleOf(t, db, loketB1)

	w := post(r, key, "/mpp/v1/antrian/"+ticket+"/skip", "")

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	require.Equal(t, "SKIPPED", load(t, db, ticket).Status)
	require.True(t, idleOf(t, db, loketB1).After(before), "the loket must be marked idle again")
}

func TestSkip_RequiresAuthentication(t *testing.T) {
	db := testsupport.Postgres(t)
	resetService(t, db)
	r := newServer(t, db)
	ticket := waiting(t, db, 1)

	w := post(r, "", "/mpp/v1/antrian/"+ticket+"/skip", "")

	require.Equal(t, http.StatusUnauthorized, w.Code, "body: %s", w.Body.String())
}
