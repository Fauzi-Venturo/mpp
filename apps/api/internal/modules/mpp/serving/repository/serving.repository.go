package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/serving/domain"
	"github.com/ndollem/mpp/apps/api/pkg/logger"
)

var (
	// ErrAntrianNotFound means no such ticket exists.
	ErrAntrianNotFound = errors.New("antrian not found")
	// ErrNotServing means the ticket is not in a state that can be finished.
	ErrNotServing = errors.New("antrian is not being served")
)

type ServingRepository struct {
	db *pgxpool.Pool
}

func NewServingRepository(db *pgxpool.Pool) *ServingRepository {
	return &ServingRepository{db: db}
}

// finishQuery is guarded by status = 'SERVING', the only legal source of DONE
// (queue-state-machine.md). The guard is what makes a double finish impossible:
// two operators both reach the UPDATE, only the one seeing the row still SERVING
// changes it. A read-then-write here would let both through.
const finishQuery = `
	UPDATE mpp.antrian
	SET status = 'DONE', done_at = NOW(), updated_at = NOW()
	WHERE id = $1 AND status = 'SERVING'
	RETURNING id, nomor, status, loket_id, served_at, done_at`

// Finish closes the ticket, its serving session and the loket's idle clock in one
// transaction, so a half-finished service can never be observed.
func (r *ServingRepository) Finish(ctx context.Context, antrianID string) (*domain.Completion, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.Error("Failed to begin finish transaction", logger.Err(err))
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var (
		completion domain.Completion
		loketID    *string
		servedAt   *time.Time
	)

	err = tx.QueryRow(ctx, finishQuery, antrianID).Scan(
		&completion.AntrianID, &completion.Nomor, &completion.Status,
		&loketID, &servedAt, &completion.DoneAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, r.classify(ctx, antrianID)
		}
		logger.Error("Failed to finish the antrian", logger.Err(err))
		return nil, err
	}

	// Close the open session if there is one. Slice 5 owns /start, so a ticket may
	// legitimately have none — the authoritative timestamps live on the antrian row
	// itself, and there is no DB constraint tying the two together.
	if _, err := tx.Exec(ctx,
		`UPDATE mpp.serving_session SET ended_at = NOW(), outcome = 'DONE'
		 WHERE antrian_id = $1 AND ended_at IS NULL`, antrianID); err != nil {
		logger.Error("Failed to close the serving session", logger.Err(err))
		return nil, err
	}

	// Marking the loket idle again is what feeds idle-longest allocation (BR-12).
	if loketID != nil {
		if _, err := tx.Exec(ctx,
			`UPDATE mpp.loket SET last_idle_at = NOW(), updated_at = NOW() WHERE id = $1`,
			*loketID); err != nil {
			logger.Error("Failed to mark the loket idle", logger.Err(err))
			return nil, err
		}
	}

	if servedAt != nil {
		seconds := int(completion.DoneAt.Sub(*servedAt).Round(time.Second).Seconds())
		completion.DurationSeconds = &seconds
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("Failed to commit the finish transaction", logger.Err(err))
		return nil, err
	}

	return &completion, nil
}

// classify tells "no such ticket" (404) apart from "wrong state" (409) after the
// guarded UPDATE matched nothing.
func (r *ServingRepository) classify(ctx context.Context, antrianID string) error {
	var status string
	err := r.db.QueryRow(ctx, `SELECT status FROM mpp.antrian WHERE id = $1`, antrianID).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAntrianNotFound
		}
		return err
	}
	return ErrNotServing
}

// currentCallQuery finds what the TV should be showing: the newest CALLED or
// SERVING ticket of the day for this instansi. CURRENT_DATE scopes it because
// numbers reset at midnight (BR-03).
const currentCallQuery = `
	SELECT a.id, a.nomor, a.status, a.call_count,
	       COALESCE(NULLIF(l.name, ''), l.code, ''),
	       COALESCE(a.called_at, a.queued_at)
	FROM mpp.antrian a
	LEFT JOIN mpp.loket l ON l.id = a.loket_id AND l.deleted_at IS NULL
	WHERE a.instansi_id = $1 AND a.queue_date = CURRENT_DATE
	  AND a.status IN ('CALLED', 'SERVING')
	ORDER BY COALESCE(a.called_at, a.queued_at) DESC
	LIMIT 1`

const nextUpQuery = `
	SELECT id, nomor
	FROM mpp.antrian
	WHERE instansi_id = $1 AND queue_date = CURRENT_DATE AND status = 'WAITING'
	ORDER BY nomor_seq
	LIMIT $2`

// CurrentCall returns the ticket on screen, or nil between calls.
func (r *ServingRepository) CurrentCall(ctx context.Context, instansiID string) (*domain.CurrentCall, error) {
	var call domain.CurrentCall

	err := r.db.QueryRow(ctx, currentCallQuery, instansiID).Scan(
		&call.AntrianID, &call.Nomor, &call.Status, &call.CallCount, &call.Loket, &call.CalledAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.Error("Failed to read the current call", logger.Err(err))
		return nil, err
	}

	return &call, nil
}

// NextUp returns the numbers queued behind the current call.
func (r *ServingRepository) NextUp(ctx context.Context, instansiID string, limit int) ([]domain.NextUp, error) {
	rows, err := r.db.Query(ctx, nextUpQuery, instansiID, limit)
	if err != nil {
		logger.Error("Failed to read the next-up list", logger.Err(err))
		return nil, err
	}
	defer rows.Close()

	list := []domain.NextUp{}
	for rows.Next() {
		var n domain.NextUp
		if err := rows.Scan(&n.AntrianID, &n.Nomor); err != nil {
			logger.Error("Failed to scan a next-up row", logger.Err(err))
			return nil, err
		}
		list = append(list, n)
	}

	return list, rows.Err()
}
