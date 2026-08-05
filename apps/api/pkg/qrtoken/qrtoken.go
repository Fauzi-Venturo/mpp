// Package qrtoken issues and validates the single-use check-in tokens carried by
// booking QR codes (BR-09). Tokens are crypto-random, bound to the booking day plus
// a configurable window, and refuse a second scan.
//
// It is deliberately storage-free: callers persist Value/ExpiresAt on mpp.booking
// and pass what they read back in.
package qrtoken

import (
	"crypto/subtle"
	"errors"
	"time"

	"github.com/ndollem/mpp/apps/api/pkg/token"
)

var (
	// ErrMismatch means the presented token is not the one on the booking.
	ErrMismatch = errors.New("qr token does not match")
	// ErrUsed means the token was already scanned — a replay attempt.
	ErrUsed = errors.New("qr token already used")
	// ErrExpired means the check-in window has closed.
	ErrExpired = errors.New("qr token expired")
	// ErrWrongDay means check-in was attempted before the booking day started.
	ErrWrongDay = errors.New("qr token not valid on this day")
)

// tokenBytes yields a 64-character hex string, matching pkg/token's other tokens.
const tokenBytes = 32

// Token is a freshly issued check-in token.
type Token struct {
	Value     string
	ExpiresAt time.Time
}

// Stored is the token state as persisted on a booking.
type Stored struct {
	Value       string
	BookingDate time.Time
	ExpiresAt   time.Time
	UsedAt      *time.Time // nil while the token is unused
}

// Issue creates a token for a booking. The window is measured from the start of the
// booking day (UTC), so an afternoon booking still expires with its own day.
func Issue(bookingDate time.Time, window time.Duration) (Token, error) {
	value, err := token.GenerateSecureToken(tokenBytes)
	if err != nil {
		return Token{}, err
	}

	return Token{Value: value, ExpiresAt: startOfDay(bookingDate).Add(window)}, nil
}

// Validate reports why a scan must be refused, or nil when the token may be spent.
// Reuse outranks expiry: a replayed token stays "used" however late it is presented.
func Validate(s Stored, presented string, now time.Time) error {
	if subtle.ConstantTimeCompare([]byte(presented), []byte(s.Value)) != 1 {
		return ErrMismatch
	}
	if s.UsedAt != nil {
		return ErrUsed
	}
	if now.Before(startOfDay(s.BookingDate)) {
		return ErrWrongDay
	}
	if now.After(s.ExpiresAt) {
		return ErrExpired
	}

	return nil
}

func startOfDay(t time.Time) time.Time {
	utc := t.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}
