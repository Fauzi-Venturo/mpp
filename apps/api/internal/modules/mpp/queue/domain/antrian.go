package domain

import "time"

// AntrianStatus mirrors the mpp.antrian_status enum.
type AntrianStatus string

const (
	StatusWaiting AntrianStatus = "WAITING"
	StatusCalled  AntrianStatus = "CALLED"
	StatusServing AntrianStatus = "SERVING"
	StatusDone    AntrianStatus = "DONE"
	StatusSkipped AntrianStatus = "SKIPPED"
)

// AntrianSource mirrors the mpp.antrian_source enum.
type AntrianSource string

const (
	SourceBooking AntrianSource = "BOOKING"
	SourceWalkIn  AntrianSource = "WALK_IN"
)

// Antrian is one ticket in a service's daily stream.
type Antrian struct {
	ID             string        `json:"id"`
	BookingID      *string       `json:"booking_id,omitempty"`
	PemohonID      string        `json:"pemohon_id"`
	InstansiID     string        `json:"instansi_id"`
	JenisLayananID string        `json:"jenis_layanan_id"`
	Nomor          string        `json:"nomor"`
	NomorSeq       int           `json:"nomor_seq"`
	QueueDate      time.Time     `json:"queue_date"`
	Source         AntrianSource `json:"source"`
	Status         AntrianStatus `json:"status"`
	QueuedAt       time.Time     `json:"queued_at"`
}

// EnqueueParams is everything needed to put one applicant into a stream.
type EnqueueParams struct {
	BookingID      *string
	PemohonID      string
	InstansiID     string
	JenisLayananID string
	QueueDate      time.Time
	Source         AntrianSource
}
