package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/queue/domain"
	"github.com/ndollem/mpp/apps/api/pkg/logger"
)

type AntrianRepository struct {
	db *pgxpool.Pool
}

func NewAntrianRepository(db *pgxpool.Pool) *AntrianRepository {
	return &AntrianRepository{db: db}
}

// PrefixAndMaxSeqTx returns the instansi letter used in queue numbers (BR-01) and
// the highest number already handed out for that service/day. The maximum is what
// a cold Redis reseeds from — system-architecture.md calls the counter
// "rebuildable from Postgres".
func (r *AntrianRepository) PrefixAndMaxSeqTx(ctx context.Context, tx pgx.Tx, instansiID, layananID string, day time.Time) (string, int, error) {
	var (
		prefix string
		maxSeq int
	)

	err := tx.QueryRow(ctx, `SELECT prefix FROM mpp.instansi WHERE id = $1 AND deleted_at IS NULL`,
		instansiID).Scan(&prefix)
	if err != nil {
		logger.Error("Failed to read instansi prefix", logger.Err(err))
		return "", 0, err
	}

	err = tx.QueryRow(ctx,
		`SELECT COALESCE(MAX(nomor_seq), 0) FROM mpp.antrian
		 WHERE jenis_layanan_id = $1 AND queue_date = $2`,
		layananID, day).Scan(&maxSeq)
	if err != nil {
		logger.Error("Failed to read the highest queue number", logger.Err(err))
		return "", 0, err
	}

	return prefix, maxSeq, nil
}

// InsertTx writes the ticket. antrian_service_day_seq_key is the last line of
// defence: if Redis ever drifts behind Postgres this INSERT fails rather than
// handing two citizens the same number.
func (r *AntrianRepository) InsertTx(ctx context.Context, tx pgx.Tx, p domain.EnqueueParams, nomor string, seq int) (*domain.Antrian, error) {
	antrian := domain.Antrian{}

	err := tx.QueryRow(ctx,
		`INSERT INTO mpp.antrian (booking_id, pemohon_id, instansi_id, jenis_layanan_id,
		                          nomor, nomor_seq, queue_date, source, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'WAITING')
		 RETURNING id, booking_id, pemohon_id, instansi_id, jenis_layanan_id,
		           nomor, nomor_seq, queue_date, source, status, queued_at`,
		p.BookingID, p.PemohonID, p.InstansiID, p.JenisLayananID, nomor, seq, p.QueueDate, p.Source,
	).Scan(&antrian.ID, &antrian.BookingID, &antrian.PemohonID, &antrian.InstansiID,
		&antrian.JenisLayananID, &antrian.Nomor, &antrian.NomorSeq, &antrian.QueueDate,
		&antrian.Source, &antrian.Status, &antrian.QueuedAt)
	if err != nil {
		logger.Error("Failed to insert antrian", logger.Err(err))
		return nil, err
	}

	return &antrian, nil
}

const listWaitingQuery = `
	SELECT id, booking_id, pemohon_id, instansi_id, jenis_layanan_id,
	       nomor, nomor_seq, queue_date, source, status, queued_at
	FROM mpp.antrian
	WHERE jenis_layanan_id = $1 AND queue_date = CURRENT_DATE AND status = 'WAITING'
	ORDER BY nomor_seq`

// ListWaiting returns today's waiting stream for one service, oldest first.
func (r *AntrianRepository) ListWaiting(ctx context.Context, layananID string) ([]domain.Antrian, error) {
	rows, err := r.db.Query(ctx, listWaitingQuery, layananID)
	if err != nil {
		logger.Error("Failed to list the waiting queue", logger.Err(err))
		return nil, err
	}
	defer rows.Close()

	stream := []domain.Antrian{}
	for rows.Next() {
		var a domain.Antrian
		if err := rows.Scan(&a.ID, &a.BookingID, &a.PemohonID, &a.InstansiID, &a.JenisLayananID,
			&a.Nomor, &a.NomorSeq, &a.QueueDate, &a.Source, &a.Status, &a.QueuedAt); err != nil {
			logger.Error("Failed to scan an antrian row", logger.Err(err))
			return nil, err
		}
		stream = append(stream, a)
	}

	return stream, rows.Err()
}
