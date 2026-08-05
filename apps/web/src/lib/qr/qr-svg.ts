import QRCode from 'qrcode';

/**
 * Renders a check-in token as an inline SVG string.
 *
 * SVG (not canvas/PNG) so the code stays crisp when printed and can be handed to a
 * download link without a rasterisation step.
 */
export async function renderQrSvg(token: string): Promise<string> {
  if (!token.trim()) {
    throw new Error('QR token is required');
  }

  return QRCode.toString(token, { type: 'svg', margin: 1, errorCorrectionLevel: 'M' });
}

/** Builds a file name a citizen can recognise in their downloads folder. */
export function qrFileName(reference: string): string {
  const slug = reference
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');

  return `qr-${slug}.svg`;
}
