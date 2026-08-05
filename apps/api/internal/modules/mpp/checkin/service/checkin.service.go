package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/checkin/domain"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/checkin/dto"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/checkin/repository"
	queueDomain "github.com/ndollem/mpp/apps/api/internal/modules/mpp/queue/domain"
	"github.com/ndollem/mpp/apps/api/pkg/qrtoken"
	"github.com/ndollem/mpp/apps/api/pkg/types"
)

var (
	// ErrTokenNotFound maps to 404.
	ErrTokenNotFound = repository.ErrTokenNotFound
	// ErrTokenRejected covers reuse, expiry and wrong-day scans — all 409.
	ErrTokenRejected = errors.New("qr token rejected")
)

// Enqueuer puts a checked-in applicant into the service's daily stream, using the
// caller's transaction so the ticket and the check-in commit together. Implemented
// by the queue module; injected in the router to keep the modules independent.
type Enqueuer interface {
	EnqueueTx(ctx context.Context, tx pgx.Tx, p queueDomain.EnqueueParams) (*queueDomain.Antrian, error)
}

type CheckinService struct {
	checkinRepo *repository.CheckinRepository
	enqueuer    Enqueuer
}

func NewCheckinService(checkinRepo *repository.CheckinRepository) *CheckinService {
	return &CheckinService{checkinRepo: checkinRepo}
}

// SetEnqueuer injects the queue module so check-in also hands out a number.
func (s *CheckinService) SetEnqueuer(e Enqueuer) {
	s.enqueuer = e
}

// Create validates a scanned token and checks the booking in. Token rules live in
// pkg/qrtoken (BR-09); this only maps storage to it and back.
func (s *CheckinService) Create(ctx context.Context, req dto.CheckinRequest) (*domain.Checkin, error) {
	stored, err := s.checkinRepo.FindByToken(ctx, req.QRToken)
	if err != nil {
		return nil, err
	}

	err = qrtoken.Validate(qrtoken.Stored{
		Value:       stored.Token,
		BookingDate: stored.BookingDate,
		ExpiresAt:   stored.ExpiresAt,
		UsedAt:      stored.CheckedInAt, // the schema records "used" as checked_in_at
		Location:    types.ServiceZone, // the booking day is a local day, not a UTC one
	}, req.QRToken, time.Now().UTC())
	if err != nil {
		return nil, errors.Join(ErrTokenRejected, err)
	}

	result, err := s.checkinRepo.MarkCheckedIn(ctx, req.QRToken, s.enqueue)
	if err != nil {
		// A concurrent scan won between the read and the guarded UPDATE.
		if errors.Is(err, repository.ErrAlreadyCheckedIn) {
			return nil, errors.Join(ErrTokenRejected, qrtoken.ErrUsed)
		}
		return nil, err
	}

	return result, nil
}

// enqueue runs inside the check-in transaction: the citizen leaves the kiosk with
// a number, or the check-in never happened.
func (s *CheckinService) enqueue(ctx context.Context, tx pgx.Tx, c *domain.Checkin) error {
	if s.enqueuer == nil {
		return nil
	}

	bookingID := c.BookingID
	antrian, err := s.enqueuer.EnqueueTx(ctx, tx, queueDomain.EnqueueParams{
		BookingID:      &bookingID,
		PemohonID:      c.PemohonID,
		InstansiID:     c.InstansiID,
		JenisLayananID: c.JenisLayananID,
		QueueDate:      c.Tanggal,
		Source:         queueDomain.SourceBooking,
	})
	if err != nil {
		return err
	}

	c.Antrian = antrian
	return nil
}
