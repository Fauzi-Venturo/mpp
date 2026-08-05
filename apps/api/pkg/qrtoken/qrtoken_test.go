package qrtoken_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ndollem/mpp/apps/api/pkg/qrtoken"
)

// BR-09: QR tokens are single-use and bound to the booking day + a window.
// Reused, expired or wrong-day tokens are rejected (docs/02-domain/business-rules.md).

var bookingDay = time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)

const dayWindow = 24 * time.Hour

func TestIssueGivesADifferentTokenEveryTime(t *testing.T) {
	first, err := qrtoken.Issue(bookingDay, dayWindow)
	require.NoError(t, err)
	second, err := qrtoken.Issue(bookingDay, dayWindow)
	require.NoError(t, err)

	assert.NotEqual(t, first.Value, second.Value)
	assert.NotEmpty(t, first.Value)
}

func TestIssueExpiresAtTheEndOfTheBookingWindow(t *testing.T) {
	// A booking timestamped mid-day still expires relative to the start of that day.
	issued, err := qrtoken.Issue(bookingDay.Add(15*time.Hour), dayWindow)
	require.NoError(t, err)

	assert.Equal(t, bookingDay.Add(dayWindow), issued.ExpiresAt)
}

func TestValidateAcceptsAnUnusedTokenOnTheBookingDay(t *testing.T) {
	issued, err := qrtoken.Issue(bookingDay, dayWindow)
	require.NoError(t, err)

	stored := qrtoken.Stored{
		Value:       issued.Value,
		BookingDate: bookingDay,
		ExpiresAt:   issued.ExpiresAt,
	}

	assert.NoError(t, qrtoken.Validate(stored, issued.Value, bookingDay.Add(9*time.Hour)))
}

func TestValidateRejectsAnUnknownToken(t *testing.T) {
	stored := qrtoken.Stored{
		Value:       "the-real-one",
		BookingDate: bookingDay,
		ExpiresAt:   bookingDay.Add(dayWindow),
	}

	err := qrtoken.Validate(stored, "a-guess", bookingDay.Add(9*time.Hour))

	assert.ErrorIs(t, err, qrtoken.ErrMismatch)
}

func TestValidateRejectsAReusedToken(t *testing.T) {
	usedAt := bookingDay.Add(8 * time.Hour)
	stored := qrtoken.Stored{
		Value:       "already-scanned",
		BookingDate: bookingDay,
		ExpiresAt:   bookingDay.Add(dayWindow),
		UsedAt:      &usedAt,
	}

	err := qrtoken.Validate(stored, "already-scanned", bookingDay.Add(9*time.Hour))

	assert.ErrorIs(t, err, qrtoken.ErrUsed)
}

func TestValidateRejectsAnExpiredToken(t *testing.T) {
	stored := qrtoken.Stored{
		Value:       "yesterdays-token",
		BookingDate: bookingDay,
		ExpiresAt:   bookingDay.Add(dayWindow),
	}

	err := qrtoken.Validate(stored, "yesterdays-token", bookingDay.Add(dayWindow+time.Second))

	assert.ErrorIs(t, err, qrtoken.ErrExpired)
}

func TestValidateRejectsCheckInBeforeTheBookingDay(t *testing.T) {
	stored := qrtoken.Stored{
		Value:       "tomorrows-token",
		BookingDate: bookingDay,
		ExpiresAt:   bookingDay.Add(dayWindow),
	}

	err := qrtoken.Validate(stored, "tomorrows-token", bookingDay.Add(-2*time.Hour))

	assert.ErrorIs(t, err, qrtoken.ErrWrongDay)
}

func TestValidateReportsReuseEvenAfterExpiry(t *testing.T) {
	// A token scanned once must keep reporting "used" — the more specific reason —
	// instead of degrading into "expired" once the day rolls over.
	usedAt := bookingDay.Add(8 * time.Hour)
	stored := qrtoken.Stored{
		Value:       "scanned-then-expired",
		BookingDate: bookingDay,
		ExpiresAt:   bookingDay.Add(dayWindow),
		UsedAt:      &usedAt,
	}

	err := qrtoken.Validate(stored, "scanned-then-expired", bookingDay.Add(48*time.Hour))

	assert.ErrorIs(t, err, qrtoken.ErrUsed)
}
