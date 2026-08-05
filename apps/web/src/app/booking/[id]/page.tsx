import type { Metadata } from 'next';

import { notFound } from 'next/navigation';

import { ApiError, getBooking } from 'src/lib/api';

import { BookingTicketView } from 'src/sections/booking/view/booking-ticket-view';

// ----------------------------------------------------------------------
// Always dynamic: a ticket is per-citizen and its QR is time-bound (BR-09), so
// it must never be prerendered, cached, or indexed.

export const dynamic = 'force-dynamic';

export const metadata: Metadata = {
  title: 'Tiket antrean',
  robots: { index: false, follow: false },
};

type Props = {
  params: Promise<{ id: string }>;
};

export default async function Page({ params }: Props) {
  const { id } = await params;

  let booking;
  try {
    booking = await getBooking(id);
  } catch (error) {
    // Only a definitive 404 is "no such booking"; a transient backend failure
    // must throw so the error boundary shows a retryable state.
    if (error instanceof ApiError && error.status === 404) {
      notFound();
    }
    throw error;
  }

  return <BookingTicketView booking={booking} />;
}
