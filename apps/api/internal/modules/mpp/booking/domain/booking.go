package domain

import (
	"time"

	"github.com/ndollem/mpp/apps/api/pkg/types"
)

// BookingStatus mirrors the mpp.booking_status enum.
type BookingStatus string

const (
	BookingStatusBooked    BookingStatus = "BOOKED"
	BookingStatusCheckedIn BookingStatus = "CHECKED_IN"
	BookingStatusExpired   BookingStatus = "EXPIRED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
)

// BookingChannel mirrors the mpp.booking_channel enum.
type BookingChannel string

const (
	BookingChannelWeb      BookingChannel = "WEB"
	BookingChannelWhatsApp BookingChannel = "WHATSAPP"
)

// Pemohon is the applicant a booking belongs to.
type Pemohon struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Phone *string `json:"phone,omitempty"`
	Email *string `json:"email,omitempty"`
}

// Booking is a scheduled registration, before check-in.
type Booking struct {
	ID             string         `json:"id"`
	PemohonID      string         `json:"pemohon_id"`
	InstansiID     string         `json:"instansi_id"`
	JenisLayananID string         `json:"jenis_layanan_id"`
	Tanggal        types.Date     `json:"tanggal"`
	Channel        BookingChannel `json:"channel"`
	Status         BookingStatus  `json:"status"`
	// QRToken is the single-use check-in token (BR-09); QRExpiresAt closes its window.
	QRToken     string    `json:"qr_token"`
	QRExpiresAt time.Time `json:"qr_expires_at"`
	CreatedAt   time.Time `json:"created_at"`
	Pemohon     *Pemohon  `json:"pemohon,omitempty"`
}
