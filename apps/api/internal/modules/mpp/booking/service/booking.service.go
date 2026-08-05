package service

import (
	"context"
	"errors"
	"time"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/booking/domain"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/booking/dto"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/booking/repository"
	"github.com/ndollem/mpp/apps/api/pkg/qrtoken"
	"github.com/ndollem/mpp/apps/api/pkg/types"
)

// defaultCheckinWindow keeps a QR token usable for the whole booking day, counted
// from 00:00 of that day (BR-09).
// ponytail: a constant until mpp.system_config actually seeds `checkin_window`
// (migrations/mpp/000005_config.up.sql); read it from there once per-instansi
// windows are configured.
const defaultCheckinWindow = 24 * time.Hour

var (
	// ErrTenantRequired means the X-Company-Slug header was missing.
	ErrTenantRequired = errors.New("company slug is required")
	// ErrBookingNotFound maps to 404 on the detail endpoint.
	ErrBookingNotFound = repository.ErrBookingNotFound
	// ErrServiceNotFound covers unknown tenant, unknown instansi and unknown
	// service alike — a public endpoint must not disclose which one it was.
	ErrServiceNotFound = repository.ErrServiceNotFound
	// ErrQuotaFull maps to 409.
	ErrQuotaFull = repository.ErrQuotaFull
)

type BookingService struct {
	bookingRepo *repository.BookingRepository
}

func NewBookingService(bookingRepo *repository.BookingRepository) *BookingService {
	return &BookingService{bookingRepo: bookingRepo}
}

func (s *BookingService) Create(ctx context.Context, companySlug string, req dto.CreateBookingRequest) (*domain.Booking, error) {
	if companySlug == "" {
		return nil, ErrTenantRequired
	}

	ok, err := s.bookingRepo.ServiceInTenant(ctx, companySlug, req.InstansiID, req.JenisLayananID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrServiceNotFound
	}

	channel := domain.BookingChannel(req.Channel)
	if channel == "" {
		channel = domain.BookingChannelWeb
	}

	qr, err := qrtoken.Issue(req.Tanggal.Time, defaultCheckinWindow, types.ServiceZone)
	if err != nil {
		return nil, err
	}

	return s.bookingRepo.Create(ctx, repository.CreateParams{
		CompanySlug:    companySlug,
		InstansiID:     req.InstansiID,
		JenisLayananID: req.JenisLayananID,
		Tanggal:        req.Tanggal,
		Channel:        channel,
		PemohonName:    req.Pemohon.Name,
		PemohonPhone:   optional(req.Pemohon.Phone),
		PemohonEmail:   optional(req.Pemohon.Email),
		QRToken:        qr.Value,
		QRExpiresAt:    qr.ExpiresAt,
	})
}

// GetByID returns one booking, including its QR token, for the ticket screen.
func (s *BookingService) GetByID(ctx context.Context, id string) (*domain.Booking, error) {
	return s.bookingRepo.GetByID(ctx, id)
}

func optional(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
