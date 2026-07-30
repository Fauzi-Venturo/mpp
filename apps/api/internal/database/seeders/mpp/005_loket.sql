-- =====================================================
-- MPP SEEDER - LOKET & LOKET_LAYANAN
-- =====================================================
-- Seeder 005: counters per agency and their service mappings.
-- Seeded as OPEN so the queue is demoable out of the box.
-- =====================================================

-- -----------------------------------------------------
-- Lokets
-- -----------------------------------------------------
INSERT INTO mpp.loket (id, instansi_id, code, name, status, is_active) VALUES
-- Dukcapil (A)
('a5000000-0000-0000-0000-000000000001', 'a2000000-0000-0000-0000-000000000001', 'L-A1', 'Loket 1', 'OPEN', TRUE),
('a5000000-0000-0000-0000-000000000002', 'a2000000-0000-0000-0000-000000000001', 'L-A2', 'Loket 2', 'OPEN', TRUE),
('a5000000-0000-0000-0000-000000000003', 'a2000000-0000-0000-0000-000000000001', 'L-A3', 'Loket 3', 'OPEN', TRUE),
-- Imigrasi (B)
('a5000000-0000-0000-0000-000000000011', 'a2000000-0000-0000-0000-000000000002', 'L-B1', 'Loket 1', 'OPEN', TRUE),
('a5000000-0000-0000-0000-000000000012', 'a2000000-0000-0000-0000-000000000002', 'L-B2', 'Loket 2', 'OPEN', TRUE),
-- BPJS Kesehatan (C)
('a5000000-0000-0000-0000-000000000021', 'a2000000-0000-0000-0000-000000000003', 'L-C1', 'Loket 1', 'OPEN', TRUE)
ON CONFLICT DO NOTHING;

-- -----------------------------------------------------
-- Loket <-> Service mapping (each loket serves all services of its agency)
-- -----------------------------------------------------
INSERT INTO mpp.loket_layanan (loket_id, jenis_layanan_id)
SELECT l.id, s.id
FROM mpp.loket l
JOIN mpp.jenis_layanan s ON s.instansi_id = l.instansi_id
WHERE l.instansi_id IN (
    'a2000000-0000-0000-0000-000000000001',
    'a2000000-0000-0000-0000-000000000002',
    'a2000000-0000-0000-0000-000000000003'
)
ON CONFLICT DO NOTHING;
