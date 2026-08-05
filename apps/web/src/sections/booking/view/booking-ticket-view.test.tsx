import type { Booking } from 'src/lib/api/booking';

import { it, expect, describe } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';

import { BookingTicketView } from 'src/sections/booking/view/booking-ticket-view';

const booking: Booking = {
  id: '9f1d0c3a-0000-0000-0000-000000000001',
  pemohon_id: '8e1d0c3a-0000-0000-0000-000000000002',
  instansi_id: 'a2000000-0000-0000-0000-000000000001',
  jenis_layanan_id: 'b3000000-0000-0000-0000-000000000003',
  tanggal: '2026-08-10',
  channel: 'WEB',
  status: 'BOOKED',
  qr_token: 'f1e2d3c4b5a6',
  qr_expires_at: '2026-08-10T16:59:59Z',
  pemohon: { id: '8e1d0c3a-0000-0000-0000-000000000002', name: 'Budi Santoso' },
};

const beforeExpiry = new Date('2026-08-10T02:00:00Z');

describe('BookingTicketView', () => {
  it('shows the QR code for a live booking', async () => {
    render(<BookingTicketView booking={booking} now={beforeExpiry} />);

    await waitFor(() => expect(screen.getByRole('img', { name: /qr/i })).toBeInTheDocument());
  });

  it('names the applicant and the service date', () => {
    render(<BookingTicketView booking={booking} now={beforeExpiry} />);

    expect(screen.getByText(/Budi Santoso/)).toBeInTheDocument();
    // Specific on purpose: the expiry line also carries this date.
    expect(screen.getByText(/Tanggal layanan 10 Aug 2026/)).toBeInTheDocument();
  });

  it('withholds the QR when the booking is no longer active', async () => {
    render(<BookingTicketView booking={{ ...booking, status: 'CANCELLED' }} now={beforeExpiry} />);

    expect(await screen.findByText(/dibatalkan|tidak aktif/i)).toBeInTheDocument();
    expect(screen.queryByRole('img', { name: /qr/i })).not.toBeInTheDocument();
  });
});
