package service

import (
	"context"
	"errors"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/booking/domain"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/booking/dto"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/booking/repository"
)

var (
	// ErrTenantRequired means the X-Company-Slug header was missing.
	ErrTenantRequired = errors.New("company slug is required")
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

	return s.bookingRepo.Create(ctx, repository.CreateParams{
		CompanySlug:    companySlug,
		InstansiID:     req.InstansiID,
		JenisLayananID: req.JenisLayananID,
		Tanggal:        req.Tanggal,
		Channel:        channel,
		PemohonName:    req.Pemohon.Name,
		PemohonPhone:   optional(req.Pemohon.Phone),
		PemohonEmail:   optional(req.Pemohon.Email),
	})
}

func optional(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
