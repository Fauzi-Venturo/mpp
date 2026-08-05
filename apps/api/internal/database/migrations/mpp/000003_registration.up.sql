-- =====================================================
-- MPP SCHEMA - REGISTRATION
-- =====================================================
-- Migration 003: kuota_booking, pemohon, booking
-- =====================================================

-- -----------------------------------------------------
-- KUOTA_BOOKING (Quota per date, per agency/service)
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS mpp.kuota_booking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instansi_id UUID NOT NULL REFERENCES mpp.instansi(id) ON DELETE CASCADE,
    jenis_layanan_id UUID REFERENCES mpp.jenis_layanan(id) ON DELETE CASCADE, -- NULL = agency-wide

    tanggal DATE NOT NULL,
    kuota INT NOT NULL DEFAULT 0 CHECK (kuota >= 0),
    terpakai INT NOT NULL DEFAULT 0 CHECK (terpakai >= 0),

    -- Audit
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by UUID
);

-- Uniqueness split by nullability of jenis_layanan_id (mirrors core user_roles pattern)
CREATE UNIQUE INDEX kuota_service_date_key ON mpp.kuota_booking(instansi_id, jenis_layanan_id, tanggal)
    WHERE jenis_layanan_id IS NOT NULL;
CREATE UNIQUE INDEX kuota_agency_date_key ON mpp.kuota_booking(instansi_id, tanggal)
    WHERE jenis_layanan_id IS NULL;
CREATE INDEX idx_kuota_instansi_tanggal ON mpp.kuota_booking(instansi_id, tanggal);

-- -----------------------------------------------------
-- PEMOHON (Applicant)
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS mpp.pemohon (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),                   -- WhatsApp / contact
    email VARCHAR(255),
    nik_hash VARCHAR(255),               -- hashed national ID, only when a service requires it

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_pemohon_phone ON mpp.pemohon(phone);

-- -----------------------------------------------------
-- BOOKING (Scheduled registration, pre-check-in)
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS mpp.booking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pemohon_id UUID NOT NULL REFERENCES mpp.pemohon(id) ON DELETE CASCADE,
    instansi_id UUID NOT NULL REFERENCES mpp.instansi(id) ON DELETE RESTRICT,
    jenis_layanan_id UUID NOT NULL REFERENCES mpp.jenis_layanan(id) ON DELETE RESTRICT,

    tanggal DATE NOT NULL,
    channel mpp.booking_channel NOT NULL DEFAULT 'WEB',

    -- QR check-in token (single-use, time-bound)
    qr_token VARCHAR(255),
    qr_expires_at TIMESTAMPTZ,

    status mpp.booking_status NOT NULL DEFAULT 'BOOKED',
    checked_in_at TIMESTAMPTZ,

    -- Audit
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by UUID,
    deleted_at TIMESTAMPTZ,
    deleted_by UUID
);

CREATE UNIQUE INDEX booking_qr_token_key ON mpp.booking(qr_token) WHERE qr_token IS NOT NULL;
CREATE INDEX idx_booking_status_tanggal ON mpp.booking(status, tanggal);
CREATE INDEX idx_booking_pemohon_id ON mpp.booking(pemohon_id);
CREATE INDEX idx_booking_instansi_id ON mpp.booking(instansi_id);

-- =====================================================
-- COMMENTS
-- =====================================================
COMMENT ON TABLE mpp.kuota_booking IS 'Booking quota per date, per agency (or per service when jenis_layanan_id set)';
COMMENT ON COLUMN mpp.kuota_booking.terpakai IS 'Consumed quota; incremented atomically at booking time';
COMMENT ON TABLE mpp.pemohon IS 'Citizen/applicant; PII minimized per privacy policy';
COMMENT ON COLUMN mpp.pemohon.nik_hash IS 'Hashed national ID; only collected when a service requires it';
COMMENT ON TABLE mpp.booking IS 'Scheduled registration before on-site check-in';
COMMENT ON COLUMN mpp.booking.qr_token IS 'Single-use, time-bound check-in token';
