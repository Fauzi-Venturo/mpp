import { it, expect, describe } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';

import { QrTicket } from 'src/sections/booking/qr-ticket';

// Slice 2: the citizen must be able to SEE the QR and DOWNLOAD it, and must be
// warned when the token is no longer usable (BR-09: single-use + time-bound).

const token = 'f1e2d3c4b5a6';
const reference = 'DUKCAPIL-2026-08-10';
const expiresAt = '2026-08-10T16:59:59Z';

describe('QrTicket', () => {
  it('renders the QR code for the token', async () => {
    render(
      <QrTicket
        token={token}
        reference={reference}
        expiresAt={expiresAt}
        now={new Date('2026-08-10T09:00:00Z')}
      />
    );

    await waitFor(() => expect(screen.getByRole('img', { name: /qr/i })).toBeInTheDocument());
    expect(screen.getByRole('img', { name: /qr/i }).innerHTML).toContain('<svg');
  });

  it('shows when the token expires', async () => {
    render(
      <QrTicket
        token={token}
        reference={reference}
        expiresAt={expiresAt}
        now={new Date('2026-08-10T09:00:00Z')}
      />
    );

    expect(screen.getByText(/10 Aug 2026/)).toBeInTheDocument();
  });

  it('offers a download named after the booking reference', async () => {
    render(
      <QrTicket
        token={token}
        reference={reference}
        expiresAt={expiresAt}
        now={new Date('2026-08-10T09:00:00Z')}
      />
    );

    const download = await screen.findByRole('link', { name: /unduh/i });
    expect(download).toHaveAttribute('download', 'qr-dukcapil-2026-08-10.svg');
  });

  it('warns and withholds the code once the token has expired', async () => {
    render(
      <QrTicket
        token={token}
        reference={reference}
        expiresAt={expiresAt}
        now={new Date('2026-08-11T07:00:00Z')}
      />
    );

    expect(await screen.findByText(/kedaluwarsa/i)).toBeInTheDocument();
    expect(screen.queryByRole('img', { name: /qr/i })).not.toBeInTheDocument();
    expect(screen.queryByRole('link', { name: /unduh/i })).not.toBeInTheDocument();
  });
});
