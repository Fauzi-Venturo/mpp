'use client';

import type { Booking } from 'src/lib/api/booking';

import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import Alert from '@mui/material/Alert';
import Stack from '@mui/material/Stack';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';

import { fDate } from 'src/utils/format-time';

import { QrTicket } from 'src/sections/booking/qr-ticket';

// Data source: GET /mpp/v1/booking/{id} via src/lib/api/booking.ts, fetched by
// the server component in src/app/booking/[id]/page.tsx.

type Props = {
  booking: Booking;
  /** Injectable clock — production passes nothing. */
  now?: Date;
};

// Only a live booking may be scanned; anything else must not show a code the kiosk
// would refuse anyway (BR-09). The reason matters to the citizen: "already checked
// in" is good news, "cancelled" is not.
const STATUS_NOTICE: Record<string, { severity: 'info' | 'success' | 'warning'; text: string }> = {
  CHECKED_IN: {
    severity: 'info',
    text: 'Anda sudah check-in. Nomor antrean sudah terbit — silakan tunggu dipanggil.',
  },
  SERVING: {
    severity: 'info',
    text: 'Anda sudah check-in dan sedang dilayani di loket.',
  },
  DONE: {
    severity: 'success',
    text: 'Layanan sudah selesai dilayani. Terima kasih.',
  },
  EXPIRED: {
    severity: 'warning',
    text: 'Booking ini kedaluwarsa karena tidak hadir pada tanggal layanan. Silakan daftar ulang.',
  },
  CANCELLED: {
    severity: 'warning',
    text: 'Booking ini sudah dibatalkan, jadi QR check-in tidak diterbitkan.',
  },
};

const FALLBACK_NOTICE = {
  severity: 'warning' as const,
  text: 'Booking ini tidak aktif, jadi QR check-in tidak diterbitkan.',
};

export function BookingTicketView({ booking, now }: Props) {
  const active = booking.status === 'BOOKED';
  const notice = STATUS_NOTICE[booking.status] ?? FALLBACK_NOTICE;

  return (
    <Container sx={{ py: 5 }}>
      <Stack spacing={3} sx={{ maxWidth: 480, mx: 'auto' }}>
        <Box>
          <Typography variant="h4">Tiket antrean</Typography>
          <Typography variant="body2" color="text.secondary">
            Tunjukkan QR ini di kiosk untuk check-in.
          </Typography>
        </Box>

        <Card sx={{ p: 3 }}>
          <Stack spacing={2}>
            <Box>
              <Typography variant="subtitle1">{booking.pemohon?.name ?? 'Pemohon'}</Typography>
              <Typography variant="body2" color="text.secondary">
                Tanggal layanan {fDate(booking.tanggal)}
              </Typography>
            </Box>

            {active ? (
              <QrTicket
                token={booking.qr_token}
                reference={booking.id}
                expiresAt={booking.qr_expires_at}
                now={now}
              />
            ) : (
              <Alert severity={notice.severity}>{notice.text}</Alert>
            )}
          </Stack>
        </Card>
      </Stack>
    </Container>
  );
}
