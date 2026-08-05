'use client';

import { useState, useEffect } from 'react';

import Box from '@mui/material/Box';
import Alert from '@mui/material/Alert';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';

import { fDateTime } from 'src/utils/format-time';

import { qrFileName, renderQrSvg } from 'src/lib/qr/qr-svg';

type Props = {
  /** Single-use check-in token issued with the booking (BR-09). */
  token: string;
  /** Human-readable booking reference, used for the download file name. */
  reference: string;
  /** ISO-8601 UTC instant after which the token is refused at the kiosk. */
  expiresAt: string;
  /** Injectable clock — production passes nothing. */
  now?: Date;
};

export function QrTicket({ token, reference, expiresAt, now = new Date() }: Props) {
  const [svg, setSvg] = useState('');

  const expired = new Date(expiresAt).getTime() <= now.getTime();

  useEffect(() => {
    let active = true;

    if (expired) {
      setSvg('');
    } else {
      renderQrSvg(token).then((markup) => {
        if (active) setSvg(markup);
      });
    }

    return () => {
      active = false;
    };
  }, [token, expired]);

  if (expired) {
    return (
      <Alert severity="warning">
        QR check-in sudah kedaluwarsa. Silakan daftar ulang atau lakukan walk-in di kiosk.
      </Alert>
    );
  }

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2 }}>
      {svg && (
        <Box
          role="img"
          aria-label={`QR check-in ${reference}`}
          sx={{ width: 240, '& svg': { width: '100%', height: 'auto' } }}
          dangerouslySetInnerHTML={{ __html: svg }}
        />
      )}

      <Typography variant="body2" color="text.secondary">
        Berlaku sampai {fDateTime(expiresAt)}
      </Typography>

      {svg && (
        <Button
          component="a"
          variant="outlined"
          href={`data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`}
          download={qrFileName(reference)}
        >
          Unduh QR
        </Button>
      )}
    </Box>
  );
}
