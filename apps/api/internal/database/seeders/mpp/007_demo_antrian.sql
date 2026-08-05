-- =====================================================
-- MPP SEEDER - DEMO TRANSACTIONS (optional)
-- =====================================================
-- Seeder 007: a couple of sample applicants and queue items so the
-- system is demoable right after db-setup. Safe to delete for a clean
-- production seed. Dates are relative to CURRENT_DATE.
-- =====================================================

-- Applicants
INSERT INTO mpp.pemohon (id, name, phone, email) VALUES
('a6000000-0000-0000-0000-000000000001', 'Siti Aminah', '+6281200000001', 'siti.demo@example.com'),
('a6000000-0000-0000-0000-000000000002', 'Budi Santoso', '+6281200000002', NULL)
ON CONFLICT DO NOTHING;

-- A booked-then-checked-in registration (Dukcapil - Perpanjangan KTP-el, today)
INSERT INTO mpp.booking (id, pemohon_id, instansi_id, jenis_layanan_id, tanggal, channel, status, checked_in_at) VALUES
(
    'a7000000-0000-0000-0000-000000000001',
    'a6000000-0000-0000-0000-000000000001',
    'a2000000-0000-0000-0000-000000000001',
    'a3000000-0000-0000-0000-000000000002',
    CURRENT_DATE,
    'WEB',
    'CHECKED_IN',
    NOW()
)
ON CONFLICT DO NOTHING;

-- Queue item from the booking (checked-in -> waiting)
INSERT INTO mpp.antrian
    (id, booking_id, pemohon_id, instansi_id, jenis_layanan_id, nomor, nomor_seq, queue_date, source, status, fo_verified) VALUES
(
    'a8000000-0000-0000-0000-000000000001',
    'a7000000-0000-0000-0000-000000000001',
    'a6000000-0000-0000-0000-000000000001',
    'a2000000-0000-0000-0000-000000000001',
    'a3000000-0000-0000-0000-000000000002',
    'A-001',
    1,
    CURRENT_DATE,
    'BOOKING',
    'WAITING',
    NULL
)
ON CONFLICT DO NOTHING;

-- A walk-in queue item (BPJS - Pendaftaran Peserta Baru, no booking)
INSERT INTO mpp.antrian
    (id, booking_id, pemohon_id, instansi_id, jenis_layanan_id, nomor, nomor_seq, queue_date, source, status) VALUES
(
    'a8000000-0000-0000-0000-000000000002',
    NULL,
    'a6000000-0000-0000-0000-000000000002',
    'a2000000-0000-0000-0000-000000000003',
    'a3000000-0000-0000-0000-000000000021',
    'C-001',
    1,
    CURRENT_DATE,
    'WALK_IN',
    'WAITING'
)
ON CONFLICT DO NOTHING;

-- Reflect the consumed booking slot in today's Dukcapil quota
UPDATE mpp.kuota_booking
SET terpakai = terpakai + 1, updated_at = NOW()
WHERE instansi_id = 'a2000000-0000-0000-0000-000000000001'
  AND jenis_layanan_id IS NULL
  AND tanggal = CURRENT_DATE
  AND terpakai = 0;
