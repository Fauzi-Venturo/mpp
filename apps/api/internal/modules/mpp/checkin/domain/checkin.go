package domain

import (
	"time"

	queueDomain "github.com/ndollem/mpp/apps/api/internal/modules/mpp/queue/domain"
)

// Pemohon is the applicant printed on the kiosk ticket.
type Pemohon struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Phone *string `json:"phone,omitempty"`
}

// Checkin is the outcome of a successful QR scan: the booking turns CHECKED_IN and
// the applicant enters the service's stream with a number (Antrian, WAITING). Both
// happen in one transaction — mpp.antrian.nomor/nomor_seq are NOT NULL with no
// default, so a ticket without a number cannot exist as an intermediate state.
type Checkin struct {
	BookingID      string      `json:"booking_id"`
	InstansiID     string      `json:"instansi_id"`
	JenisLayananID string      `json:"jenis_layanan_id"`
	Status         string      `json:"status"`
	CheckedInAt    time.Time   `json:"checked_in_at"`
	Tanggal        time.Time `json:"-"` // booking day, used as the antrian queue_date
	PemohonID      string    `json:"-"`
	Pemohon        *Pemohon  `json:"pemohon,omitempty"`

	Antrian *queueDomain.Antrian `json:"antrian,omitempty"`
}

// BookingToken is the persisted QR state a scan is judged against.
type BookingToken struct {
	BookingID   string
	Token       string
	BookingDate time.Time
	ExpiresAt   time.Time
	CheckedInAt *time.Time
	Status      string
}
