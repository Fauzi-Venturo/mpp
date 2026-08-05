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

// The building runs on a wall clock, so a "booking day" starts and ends in that
// zone — not at 00:00 UTC, which in WIB is 07:00 the next morning.
var (
	bookingDay = time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	wib        = time.FixedZone("WIB", 7*60*60)
	// The same calendar day expressed as instants: 10 Aug 00:00 WIB .. 11 Aug 00:00 WIB.
	dayStartUTC = time.Date(2026, 8, 9, 17, 0, 0, 0, time.UTC)
	dayEndUTC   = time.Date(2026, 8, 10, 17, 0, 0, 0, time.UTC)
)

const dayWindow = 24 * time.Hour

func TestIssueGivesADifferentTokenEveryTime(t *testing.T) {
	first, err := qrtoken.Issue(bookingDay, dayWindow, wib)
	require.NoError(t, err)
	second, err := qrtoken.Issue(bookingDay, dayWindow, wib)
	require.NoError(t, err)

	assert.NotEqual(t, first.Value, second.Value)
	assert.NotEmpty(t, first.Value)
}

func TestIssueExpiresWhenTheBookingDayEndsLocally(t *testing.T) {
	// A booking timestamped mid-day still expires relative to the start of that day.
	issued, err := qrtoken.Issue(bookingDay.Add(15*time.Hour), dayWindow, wib)
	require.NoError(t, err)

	assert.Equal(t, dayEndUTC, issued.ExpiresAt.UTC(),
		"midnight in the service zone, not in UTC")
}

func TestIssueFallsBackToUTCWithoutAZone(t *testing.T) {
	issued, err := qrtoken.Issue(bookingDay, dayWindow, nil)
	require.NoError(t, err)

	assert.Equal(t, bookingDay.Add(dayWindow), issued.ExpiresAt.UTC())
}

func TestValidateAcceptsAnUnusedTokenOnTheBookingDay(t *testing.T) {
	issued, err := qrtoken.Issue(bookingDay, dayWindow, wib)
	require.NoError(t, err)

	stored := qrtoken.Stored{
		Value:       issued.Value,
		BookingDate: bookingDay,
		ExpiresAt:   issued.ExpiresAt,
		Location:    wib,
	}

	assert.NoError(t, qrtoken.Validate(stored, issued.Value, dayStartUTC.Add(2*time.Hour)))
}

// A citizen arriving at 06:30 WIB on the booking day is early for the office, but
// it IS their day — rejecting them as "wrong day" is the UTC bug in disguise.
func TestValidateAcceptsAnEarlyMorningScanInTheServiceZone(t *testing.T) {
	stored := qrtoken.Stored{
		Value:       "todays-token",
		BookingDate: bookingDay,
		ExpiresAt:   dayEndUTC,
		Location:    wib,
	}

	// 06:30 WIB on the booking day == 23:30 UTC the day before.
	err := qrtoken.Validate(stored, "todays-token", time.Date(2026, 8, 9, 23, 30, 0, 0, time.UTC))

	assert.NoError(t, err)
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
