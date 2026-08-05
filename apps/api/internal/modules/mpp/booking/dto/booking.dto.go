package dto

import (
	"github.com/ndollem/mpp/apps/api/pkg/types"
)

// PemohonRequest is the applicant supplied with a booking. Bookings are public,
// so the applicant is created inline rather than referenced by id.
type PemohonRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=255"`
	Phone string `json:"phone" binding:"omitempty,max=20"`
	Email string `json:"email" binding:"omitempty,email,max=255"`
}

type CreateBookingRequest struct {
	InstansiID     string         `json:"instansi_id" binding:"required,uuid"`
	JenisLayananID string         `json:"jenis_layanan_id" binding:"required,uuid"`
	Tanggal        types.Date     `json:"tanggal" binding:"required"`
	Channel        string         `json:"channel" binding:"omitempty,oneof=WEB WHATSAPP"`
	Pemohon        PemohonRequest `json:"pemohon" binding:"required"`
}
