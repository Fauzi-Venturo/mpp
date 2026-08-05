package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/loket_ops/domain"
	"github.com/ndollem/mpp/apps/api/pkg/logger"
)

var (
	// ErrNoWaiting means the service has nobody left to call today.
	ErrNoWaiting = errors.New("no waiting antrian for this service")
	// ErrLoketNotEligible means the loket is closed, inactive, or does not serve
	// that service (BR-12/BR-15).
	ErrLoketNotEligible = errors.New("loket is not eligible for this service")
	// ErrAntrianNotFound means no antrian carries that id.
	ErrAntrianNotFound = errors.New("antrian not found")
)

type LoketOpsRepository struct {
	db *pgxpool.Pool
}

func NewLoketOpsRepository(db *pgxpool.Pool) *LoketOpsRepository {
	return &LoketOpsRepository{db: db}
}

// eligibleLoketQuery validates the loket an operator claims to be sitting at.
const eligibleLoketQuery = `
	SELECT l.id, l.code
	FROM mpp.loket l
	JOIN mpp.loket_layanan ll ON ll.loket_id = l.id AND ll.jenis_layanan_id = $2
	WHERE l.id = $1 AND l.status = 'OPEN' AND l.is_active AND l.deleted_at IS NULL`

// idleLongestLoketQuery implements BR-12: among the lokets that are OPEN and map
// the service, take the one idle the longest. Lokets already holding a CALLED or
// SERVING ticket are busy and must not receive another one.
const idleLongestLoketQuery = `
	SELECT l.id, l.code
	FROM mpp.loket l
	JOIN mpp.loket_layanan ll ON ll.loket_id = l.id AND ll.jenis_layanan_id = $1
	WHERE l.status = 'OPEN' AND l.is_active AND l.deleted_at IS NULL
	  AND NOT EXISTS (
	      SELECT 1 FROM mpp.antrian a
	      WHERE a.loket_id = l.id AND a.status IN ('CALLED', 'SERVING'))
	ORDER BY l.last_idle_at ASC
	LIMIT 1`

// nextWaitingQuery locks the ticket it hands out. SKIP LOCKED is what stops two
// operators pressing "next" at the same instant from being given the same person;
// without it both transactions read the same row and the second overwrites the first.
const nextWaitingQuery = `
	SELECT id, nomor
	FROM mpp.antrian
	WHERE jenis_layanan_id = $1 AND queue_date = CURRENT_DATE AND status = 'WAITING'
	ORDER BY priority DESC, nomor_seq ASC
	LIMIT 1
	FOR UPDATE SKIP LOCKED`

const markCalledQuery = `
	UPDATE mpp.antrian
	SET status = 'CALLED', loket_id = $2, call_count = 1, called_at = NOW(), updated_at = NOW()
	WHERE id = $1
	RETURNING call_count, called_at`

// CallNext assigns the next waiting ticket to a loket, all inside one transaction
// so the pick and the assignment cannot drift apart.
func (r *LoketOpsRepository) CallNext(ctx context.Context, layananID, loketID string) (*domain.Called, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.Error("Failed to begin call transaction", logger.Err(err))
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var loketCode string
	if loketID != "" {
		err = tx.QueryRow(ctx, eligibleLoketQuery, loketID, layananID).Scan(&loketID, &loketCode)
	} else {
		err = tx.QueryRow(ctx, idleLongestLoketQuery, layananID).Scan(&loketID, &loketCode)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrLoketNotEligible
		}
		logger.Error("Failed to resolve loket", logger.Err(err))
		return nil, err
	}

	called := domain.Called{Status: domain.StatusCalled, LoketID: loketID, Loket: loketCode}
	err = tx.QueryRow(ctx, nextWaitingQuery, layananID).Scan(&called.AntrianID, &called.Nomor)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoWaiting
		}
		logger.Error("Failed to pick the next antrian", logger.Err(err))
		return nil, err
	}

	if err := tx.QueryRow(ctx, markCalledQuery, called.AntrianID, loketID).
		Scan(&called.CallCount, &called.CalledAt); err != nil {
		logger.Error("Failed to mark antrian called", logger.Err(err))
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("Failed to commit call transaction", logger.Err(err))
		return nil, err
	}

	return &called, nil
}

// Snapshot is the antrian state a transition needs to validate itself.
type Snapshot struct {
	Status    domain.AntrianStatus
	LoketID   *string
	LoketCode string
	Nomor     string
	CallCount int
}

const lockAntrianQuery = `
	SELECT a.status, a.loket_id, COALESCE(l.code, ''), a.nomor, a.call_count
	FROM mpp.antrian a
	LEFT JOIN mpp.loket l ON l.id = a.loket_id
	WHERE a.id = $1
	FOR UPDATE OF a`

// withLockedAntrian runs fn against a row-locked antrian inside one transaction.
// Every transition guard reads and writes under the same lock, so two operators
// acting on one ticket cannot both pass the check.
func (r *LoketOpsRepository) withLockedAntrian(
	ctx context.Context,
	antrianID string,
	fn func(pgx.Tx, Snapshot) (*domain.Called, error),
) (*domain.Called, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.Error("Failed to begin transition transaction", logger.Err(err))
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var snap Snapshot
	err = tx.QueryRow(ctx, lockAntrianQuery, antrianID).
		Scan(&snap.Status, &snap.LoketID, &snap.LoketCode, &snap.Nomor, &snap.CallCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAntrianNotFound
		}
		logger.Error("Failed to lock antrian", logger.Err(err))
		return nil, err
	}

	result, err := fn(tx, snap)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("Failed to commit transition", logger.Err(err))
		return nil, err
	}

	return result, nil
}

// Recall re-announces a CALLED ticket. guard decides whether the transition is legal.
func (r *LoketOpsRepository) Recall(
	ctx context.Context,
	antrianID string,
	guard func(Snapshot) error,
) (*domain.Called, error) {
	return r.withLockedAntrian(ctx, antrianID, func(tx pgx.Tx, snap Snapshot) (*domain.Called, error) {
		if err := guard(snap); err != nil {
			return nil, err
		}

		called := domain.Called{
			AntrianID: antrianID,
			Nomor:     snap.Nomor,
			Status:    domain.StatusCalled,
			LoketID:   derefLoket(snap.LoketID),
			Loket:     snap.LoketCode,
		}
		err := tx.QueryRow(ctx,
			`UPDATE mpp.antrian
			 SET call_count = call_count + 1, called_at = NOW(), updated_at = NOW()
			 WHERE id = $1
			 RETURNING call_count, called_at`, antrianID).Scan(&called.CallCount, &called.CalledAt)
		if err != nil {
			logger.Error("Failed to recall antrian", logger.Err(err))
			return nil, err
		}

		return &called, nil
	})
}

// Start opens a serving session for a CALLED ticket.
func (r *LoketOpsRepository) Start(
	ctx context.Context,
	antrianID, userID string,
	guard func(Snapshot) error,
) (*domain.Called, error) {
	return r.withLockedAntrian(ctx, antrianID, func(tx pgx.Tx, snap Snapshot) (*domain.Called, error) {
		if err := guard(snap); err != nil {
			return nil, err
		}

		if _, err := tx.Exec(ctx,
			`UPDATE mpp.antrian SET status = 'SERVING', served_at = NOW(), updated_at = NOW()
			 WHERE id = $1`, antrianID); err != nil {
			logger.Error("Failed to start serving", logger.Err(err))
			return nil, err
		}

		if _, err := tx.Exec(ctx,
			`INSERT INTO mpp.serving_session (antrian_id, loket_id, user_id) VALUES ($1, $2, $3)`,
			antrianID, snap.LoketID, nullable(userID)); err != nil {
			logger.Error("Failed to open serving session", logger.Err(err))
			return nil, err
		}

		return &domain.Called{
			AntrianID: antrianID,
			Nomor:     snap.Nomor,
			Status:    domain.StatusServing,
			LoketID:   derefLoket(snap.LoketID),
			Loket:     snap.LoketCode,
			CallCount: snap.CallCount,
		}, nil
	})
}

// Skip closes a no-show and frees the loket, refreshing last_idle_at so the
// idle-longest allocator stays fair (queue-state-machine.md:83).
func (r *LoketOpsRepository) Skip(
	ctx context.Context,
	antrianID string,
	guard func(Snapshot) error,
) (*domain.Called, error) {
	return r.withLockedAntrian(ctx, antrianID, func(tx pgx.Tx, snap Snapshot) (*domain.Called, error) {
		if err := guard(snap); err != nil {
			return nil, err
		}

		if _, err := tx.Exec(ctx,
			`UPDATE mpp.antrian SET status = 'SKIPPED', done_at = NOW(), updated_at = NOW()
			 WHERE id = $1`, antrianID); err != nil {
			logger.Error("Failed to skip antrian", logger.Err(err))
			return nil, err
		}

		if snap.LoketID != nil {
			if _, err := tx.Exec(ctx,
				`UPDATE mpp.loket SET last_idle_at = NOW(), updated_at = NOW() WHERE id = $1`,
				*snap.LoketID); err != nil {
				logger.Error("Failed to free loket", logger.Err(err))
				return nil, err
			}
		}

		return &domain.Called{
			AntrianID: antrianID,
			Nomor:     snap.Nomor,
			Status:    domain.StatusSkipped,
			LoketID:   derefLoket(snap.LoketID),
			Loket:     snap.LoketCode,
			CallCount: snap.CallCount,
		}, nil
	})
}

func derefLoket(id *string) string {
	if id == nil {
		return ""
	}
	return *id
}

func nullable(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
