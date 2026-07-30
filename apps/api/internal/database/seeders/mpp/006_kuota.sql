-- =====================================================
-- MPP SEEDER - KUOTA BOOKING
-- =====================================================
-- Seeder 006: agency-wide booking quota for the next 7 days
-- (today .. today+6). Idempotent via NOT EXISTS on the
-- agency-level (jenis_layanan_id IS NULL) slots.
-- =====================================================

INSERT INTO mpp.kuota_booking (instansi_id, jenis_layanan_id, tanggal, kuota, terpakai)
SELECT i.id, NULL, d::date, 30, 0
FROM (VALUES
    ('a2000000-0000-0000-0000-000000000001'::uuid),
    ('a2000000-0000-0000-0000-000000000002'::uuid),
    ('a2000000-0000-0000-0000-000000000003'::uuid)
) AS i(id)
CROSS JOIN generate_series(CURRENT_DATE, CURRENT_DATE + INTERVAL '6 day', INTERVAL '1 day') AS d
WHERE NOT EXISTS (
    SELECT 1 FROM mpp.kuota_booking k
    WHERE k.instansi_id = i.id
      AND k.jenis_layanan_id IS NULL
      AND k.tanggal = d::date
);
