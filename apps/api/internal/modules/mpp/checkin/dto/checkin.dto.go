package dto

// CheckinRequest is the kiosk scan payload. Tokens are the 64-char hex values
// issued by pkg/qrtoken.
type CheckinRequest struct {
	QRToken string `json:"qr_token" binding:"required,len=64,hexadecimal"`
}
