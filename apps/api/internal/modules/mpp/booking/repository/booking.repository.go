package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/booking/domain"
	"github.com/ndollem/mpp/apps/api/pkg/logger"
	"github.com/ndollem/mpp/apps/api/pkg/types"
)

var (
	// ErrQuotaFull means the applicable kuota_booking row is exhausted (or absent).
	ErrQuotaFull = errors.New("booking quota is full")
	// ErrServiceNotFound means the service does not exist inside the caller's tenant.
	ErrServiceNotFound = errors.New("service not found for this tenant")
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}

// CreateParams carries everything one booking needs, already validated.
type CreateParams struct {
	CompanySlug    string
	InstansiID     string
	JenisLayananID string
	Tanggal        types.Date
	Channel        domain.BookingChannel
	PemohonName    string
	PemohonPhone   *string
	PemohonEmail   *string
	QRToken        string
	QRExpiresAt    time.Time
}

// ErrBookingNotFound means no live booking carries that id.
var ErrBookingNotFound = errors.New("booking not found")

const bookingByIDQuery = `
	SELECT b.id, b.pemohon_id, b.instansi_id, b.jenis_layanan_id, b.tanggal, b.channel,
	       b.status, b.qr_token, b.qr_expires_at, b.created_at,
	       p.id, p.name, p.phone, p.email
	FROM mpp.booking b
	JOIN mpp.pemohon p ON p.id = b.pemohon_id
	WHERE b.id = $1 AND b.deleted_at IS NULL`

// GetByID returns one booking with its applicant, or ErrBookingNotFound.
func (r *BookingRepository) GetByID(ctx context.Context, id string) (*domain.Booking, error) {
	var (
		booking domain.Booking
		pemohon domain.Pemohon
	)

	err := r.db.QueryRow(ctx, bookingByIDQuery, id).Scan(
		&booking.ID, &booking.PemohonID, &booking.InstansiID, &booking.JenisLayananID,
		&booking.Tanggal, &booking.Channel, &booking.Status, &booking.QRToken,
		&booking.QRExpiresAt, &booking.CreatedAt,
		&pemohon.ID, &pemohon.Name, &pemohon.Phone, &pemohon.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBookingNotFound
		}
		logger.Error("Failed to load booking", logger.Err(err))
		return nil, err
	}

	booking.Pemohon = &pemohon
	return &booking, nil
}

// Tenancy runs through the instansi, since mpp.booking carries no company_id:
// booking -> instansi -> core.companies -> core.clients.slug (the X-Company-Slug value).
const serviceInTenantQuery = `
	SELECT 1
	FROM mpp.jenis_layanan jl
	JOIN mpp.instansi i     ON i.id = jl.instansi_id AND i.deleted_at IS NULL AND i.is_active
	JOIN core.companies c   ON c.id = i.company_id   AND c.deleted_at IS NULL AND c.is_active
	JOIN core.clients cl    ON cl.id = c.client_id   AND cl.deleted_at IS NULL AND cl.slug = $1
	WHERE jl.id = $3 AND jl.instansi_id = $2 AND jl.deleted_at IS NULL AND jl.is_active`

// consumeQuotaQuery increments the applicable quota row and returns no row when it
// is already exhausted. It must stay a SINGLE statement: the UPDATE takes the row
// lock, and under READ COMMITTED Postgres re-evaluates `terpakai < kuota` against
// the freshly committed version, so a concurrent writer cannot overbook. Splitting
// this into SELECT-then-UPDATE reintroduces the race, and the happy path would
// still pass.
const consumeQuotaQuery = `
	UPDATE mpp.kuota_booking k
	SET terpakai = k.terpakai + 1, updated_at = NOW()
	WHERE k.id = (
		SELECT id FROM mpp.kuota_booking
		WHERE instansi_id = $1 AND tanggal = $2
		  AND (jenis_layanan_id = $3 OR jenis_layanan_id IS NULL)
		ORDER BY jenis_layanan_id NULLS LAST
		LIMIT 1)
	  AND k.terpakai < k.kuota
	RETURNING k.id`

// ServiceInTenant reports whether the service belongs to the tenant behind slug.
func (r *BookingRepository) ServiceInTenant(ctx context.Context, slug, instansiID, layananID string) (bool, error) {
	var one int
	err := r.db.QueryRow(ctx, serviceInTenantQuery, slug, instansiID, layananID).Scan(&one)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		logger.Error("Failed to resolve service tenancy", logger.Err(err))
		return false, err
	}
	return true, nil
}

// Create consumes one quota unit and writes the applicant and the booking in a
// single transaction, so a rolled-back booking never leaves quota consumed.
func (r *BookingRepository) Create(ctx context.Context, p CreateParams) (*domain.Booking, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.Error("Failed to begin booking transaction", logger.Err(err))
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var quotaID string
	err = tx.QueryRow(ctx, consumeQuotaQuery, p.InstansiID, p.Tanggal.Time, p.JenisLayananID).Scan(&quotaID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrQuotaFull
		}
		logger.Error("Failed to consume booking quota", logger.Err(err))
		return nil, err
	}

	pemohon := domain.Pemohon{Name: p.PemohonName, Phone: p.PemohonPhone, Email: p.PemohonEmail}
	err = tx.QueryRow(ctx,
		`INSERT INTO mpp.pemohon (name, phone, email) VALUES ($1, $2, $3) RETURNING id`,
		p.PemohonName, p.PemohonPhone, p.PemohonEmail).Scan(&pemohon.ID)
	if err != nil {
		logger.Error("Failed to insert pemohon", logger.Err(err))
		return nil, err
	}

	booking := domain.Booking{Pemohon: &pemohon}
	err = tx.QueryRow(ctx,
		`INSERT INTO mpp.booking (pemohon_id, instansi_id, jenis_layanan_id, tanggal, channel, status,
		                          qr_token, qr_expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, pemohon_id, instansi_id, jenis_layanan_id, tanggal, channel, status,
		           qr_token, qr_expires_at, created_at`,
		pemohon.ID, p.InstansiID, p.JenisLayananID, p.Tanggal.Time, p.Channel, domain.BookingStatusBooked,
		p.QRToken, p.QRExpiresAt,
	).Scan(&booking.ID, &booking.PemohonID, &booking.InstansiID, &booking.JenisLayananID,
		&booking.Tanggal, &booking.Channel, &booking.Status,
		&booking.QRToken, &booking.QRExpiresAt, &booking.CreatedAt)
	if err != nil {
		logger.Error("Failed to insert booking", logger.Err(err))
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("Failed to commit booking transaction", logger.Err(err))
		return nil, err
	}

	return &booking, nil
}
