package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/checkin/domain"
	"github.com/ndollem/mpp/apps/api/pkg/logger"
)

var (
	// ErrTokenNotFound means no live booking carries that QR token.
	ErrTokenNotFound = errors.New("qr token not found")
	// ErrAlreadyCheckedIn means another scan won the race to the booking.
	ErrAlreadyCheckedIn = errors.New("booking is no longer checkable in")
)

type CheckinRepository struct {
	db *pgxpool.Pool
}

func NewCheckinRepository(db *pgxpool.Pool) *CheckinRepository {
	return &CheckinRepository{db: db}
}

const findByTokenQuery = `
	SELECT id, qr_token, tanggal, qr_expires_at, checked_in_at, status
	FROM mpp.booking
	WHERE qr_token = $1 AND deleted_at IS NULL`

// markCheckedInQuery is guarded by status = 'BOOKED', which is what actually makes
// a token single-use: two concurrent scans both reach the UPDATE, but only the one
// that sees the row still BOOKED changes it. The read in FindByToken exists to pick
// the right rejection reason, never to decide whether the transition may happen.
const markCheckedInQuery = `
	UPDATE mpp.booking b
	SET status = 'CHECKED_IN', checked_in_at = NOW(), updated_at = NOW()
	WHERE b.qr_token = $1 AND b.status = 'BOOKED' AND b.deleted_at IS NULL
	RETURNING b.id, b.instansi_id, b.jenis_layanan_id, b.status, b.checked_in_at, b.tanggal,
	          (SELECT p.id FROM mpp.pemohon p WHERE p.id = b.pemohon_id),
	          (SELECT p.name FROM mpp.pemohon p WHERE p.id = b.pemohon_id),
	          (SELECT p.phone FROM mpp.pemohon p WHERE p.id = b.pemohon_id)`

// FindByToken returns the persisted QR state, or ErrTokenNotFound.
func (r *CheckinRepository) FindByToken(ctx context.Context, token string) (*domain.BookingToken, error) {
	var bt domain.BookingToken
	err := r.db.QueryRow(ctx, findByTokenQuery, token).Scan(
		&bt.BookingID, &bt.Token, &bt.BookingDate, &bt.ExpiresAt, &bt.CheckedInAt, &bt.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTokenNotFound
		}
		logger.Error("Failed to load booking by qr token", logger.Err(err))
		return nil, err
	}
	return &bt, nil
}

// AfterCheckedIn runs inside the same transaction as the transition, so whatever
// it writes (the queue ticket) lives or dies with the check-in.
type AfterCheckedIn func(ctx context.Context, tx pgx.Tx, c *domain.Checkin) error

// MarkCheckedIn transitions BOOKED -> CHECKED_IN and runs after inside the same
// transaction, or reports ErrAlreadyCheckedIn. Both writes must commit together:
// a booking left CHECKED_IN with no queue ticket strands the citizen, because a
// rescan then earns a 409.
func (r *CheckinRepository) MarkCheckedIn(ctx context.Context, token string, after AfterCheckedIn) (*domain.Checkin, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.Error("Failed to begin check-in transaction", logger.Err(err))
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var (
		result  domain.Checkin
		pemohon domain.Pemohon
	)

	err = tx.QueryRow(ctx, markCheckedInQuery, token).Scan(
		&result.BookingID, &result.InstansiID, &result.JenisLayananID,
		&result.Status, &result.CheckedInAt, &result.Tanggal,
		&pemohon.ID, &pemohon.Name, &pemohon.Phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAlreadyCheckedIn
		}
		logger.Error("Failed to check the booking in", logger.Err(err))
		return nil, err
	}
	result.Pemohon = &pemohon
	result.PemohonID = pemohon.ID

	if after != nil {
		if err := after(ctx, tx, &result); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("Failed to commit check-in transaction", logger.Err(err))
		return nil, err
	}

	return &result, nil
}
