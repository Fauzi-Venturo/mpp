-- =====================================================
-- MPP SEEDER - INSTANSI (Agencies)
-- =====================================================
-- Seeder 003: sample agencies within the MPP tenant company.
-- =====================================================

INSERT INTO mpp.instansi (id, company_id, name, slug, prefix, description, operating_hours, queue_mode, is_active) VALUES
(
    'a2000000-0000-0000-0000-000000000001',
    'a1000000-0000-0000-0000-000000000001',
    'Dinas Kependudukan dan Pencatatan Sipil',
    'dukcapil',
    'A',
    'Layanan administrasi kependudukan (KTP-el, Kartu Keluarga, dll).',
    '{"mon":["08:00","16:00"],"tue":["08:00","16:00"],"wed":["08:00","16:00"],"thu":["08:00","16:00"],"fri":["08:00","16:00"]}'::jsonb,
    'FIFO',
    TRUE
),
(
    'a2000000-0000-0000-0000-000000000002',
    'a1000000-0000-0000-0000-000000000001',
    'Kantor Imigrasi',
    'imigrasi',
    'B',
    'Layanan keimigrasian (paspor, dll).',
    '{"mon":["08:00","15:00"],"tue":["08:00","15:00"],"wed":["08:00","15:00"],"thu":["08:00","15:00"],"fri":["08:00","15:00"]}'::jsonb,
    'BOOKING_PRIORITY',
    TRUE
),
(
    'a2000000-0000-0000-0000-000000000003',
    'a1000000-0000-0000-0000-000000000001',
    'BPJS Kesehatan',
    'bpjs-kesehatan',
    'C',
    'Layanan kepesertaan jaminan kesehatan.',
    '{"mon":["08:00","15:00"],"tue":["08:00","15:00"],"wed":["08:00","15:00"],"thu":["08:00","15:00"],"fri":["08:00","15:00"]}'::jsonb,
    'FIFO',
    TRUE
)
ON CONFLICT DO NOTHING;
