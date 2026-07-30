-- =====================================================
-- MPP SEEDER - JENIS LAYANAN & SYARAT DOKUMEN
-- =====================================================
-- Seeder 004: services per agency and their document requirements.
-- =====================================================

-- -----------------------------------------------------
-- Services
-- -----------------------------------------------------
INSERT INTO mpp.jenis_layanan (id, instansi_id, name, description, estimasi_durasi_menit, requires_fo_verification, is_active) VALUES
-- Dukcapil (A)
('a3000000-0000-0000-0000-000000000001', 'a2000000-0000-0000-0000-000000000001', 'Perekaman KTP-el', 'Perekaman biometrik KTP elektronik.', 15, FALSE, TRUE),
('a3000000-0000-0000-0000-000000000002', 'a2000000-0000-0000-0000-000000000001', 'Pencetakan / Perpanjangan KTP-el', 'Cetak ulang atau perpanjangan KTP-el.', 10, TRUE, TRUE),
('a3000000-0000-0000-0000-000000000003', 'a2000000-0000-0000-0000-000000000001', 'Kartu Keluarga (KK)', 'Penerbitan / perubahan Kartu Keluarga.', 20, TRUE, TRUE),
-- Imigrasi (B)
('a3000000-0000-0000-0000-000000000011', 'a2000000-0000-0000-0000-000000000002', 'Pembuatan Paspor Baru', 'Permohonan paspor baru.', 30, TRUE, TRUE),
('a3000000-0000-0000-0000-000000000012', 'a2000000-0000-0000-0000-000000000002', 'Perpanjangan Paspor', 'Penggantian / perpanjangan paspor.', 20, FALSE, TRUE),
-- BPJS Kesehatan (C)
('a3000000-0000-0000-0000-000000000021', 'a2000000-0000-0000-0000-000000000003', 'Pendaftaran Peserta Baru', 'Pendaftaran kepesertaan baru.', 15, FALSE, TRUE),
('a3000000-0000-0000-0000-000000000022', 'a2000000-0000-0000-0000-000000000003', 'Perubahan Data Peserta', 'Perubahan data / faskes peserta.', 10, FALSE, TRUE)
ON CONFLICT DO NOTHING;

-- -----------------------------------------------------
-- Document requirements
-- -----------------------------------------------------
INSERT INTO mpp.syarat_dokumen (id, jenis_layanan_id, name, is_required, notes, sort) VALUES
-- Perpanjangan KTP-el
('a4000000-0000-0000-0000-000000000001', 'a3000000-0000-0000-0000-000000000002', 'Kartu Keluarga (asli)', TRUE, NULL, 1),
('a4000000-0000-0000-0000-000000000002', 'a3000000-0000-0000-0000-000000000002', 'KTP-el lama', FALSE, 'Bila masih ada.', 2),
-- Kartu Keluarga
('a4000000-0000-0000-0000-000000000003', 'a3000000-0000-0000-0000-000000000003', 'Surat pengantar RT/RW', TRUE, NULL, 1),
('a4000000-0000-0000-0000-000000000004', 'a3000000-0000-0000-0000-000000000003', 'Dokumen pendukung perubahan', FALSE, 'Buku nikah / akta, sesuai keperluan.', 2),
-- Paspor Baru
('a4000000-0000-0000-0000-000000000011', 'a3000000-0000-0000-0000-000000000011', 'KTP-el', TRUE, NULL, 1),
('a4000000-0000-0000-0000-000000000012', 'a3000000-0000-0000-0000-000000000011', 'Kartu Keluarga', TRUE, NULL, 2),
('a4000000-0000-0000-0000-000000000013', 'a3000000-0000-0000-0000-000000000011', 'Akta lahir / ijazah / buku nikah', TRUE, 'Salah satu sebagai bukti identitas.', 3),
-- Perpanjangan Paspor
('a4000000-0000-0000-0000-000000000014', 'a3000000-0000-0000-0000-000000000012', 'Paspor lama', TRUE, NULL, 1),
('a4000000-0000-0000-0000-000000000015', 'a3000000-0000-0000-0000-000000000012', 'KTP-el', TRUE, NULL, 2),
-- BPJS Pendaftaran Peserta Baru
('a4000000-0000-0000-0000-000000000021', 'a3000000-0000-0000-0000-000000000021', 'KTP-el', TRUE, NULL, 1),
('a4000000-0000-0000-0000-000000000022', 'a3000000-0000-0000-0000-000000000021', 'Kartu Keluarga', TRUE, NULL, 2),
('a4000000-0000-0000-0000-000000000023', 'a3000000-0000-0000-0000-000000000021', 'Pas foto', FALSE, 'Ukuran 3x4.', 3)
ON CONFLICT DO NOTHING;
