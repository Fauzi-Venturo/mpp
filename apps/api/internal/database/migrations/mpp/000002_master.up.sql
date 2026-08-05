-- =====================================================
-- MPP SCHEMA - MASTER DATA
-- =====================================================
-- Migration 002: instansi, jenis_layanan, syarat_dokumen,
--                loket, loket_layanan, loket_session
-- =====================================================

-- -----------------------------------------------------
-- INSTANSI (Agency)
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS mpp.instansi (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenancy: one MPP building = one company (see docs/02-domain/domain-model.md)
    company_id UUID NOT NULL REFERENCES core.companies(id) ON DELETE CASCADE,

    -- Identity
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(120) NOT NULL,
    prefix VARCHAR(4) NOT NULL,          -- queue-number prefix, e.g. 'A'
    description TEXT,
    logo_url VARCHAR(500),

    -- Operations
    operating_hours JSONB DEFAULT '{}',  -- e.g. {"mon":["08:00","16:00"], "holidays":["2026-08-17"]}
    queue_mode mpp.queue_mode NOT NULL DEFAULT 'FIFO',

    -- Status
    is_active BOOLEAN DEFAULT TRUE,

    -- Audit
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by UUID,
    deleted_at TIMESTAMPTZ,
    deleted_by UUID
);

CREATE UNIQUE INDEX instansi_company_prefix_key ON mpp.instansi(company_id, prefix)
    WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX instansi_company_slug_key ON mpp.instansi(company_id, slug)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_instansi_company_id ON mpp.instansi(company_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_instansi_is_active ON mpp.instansi(is_active) WHERE deleted_at IS NULL;

-- -----------------------------------------------------
-- JENIS_LAYANAN (Service type)
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS mpp.jenis_layanan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instansi_id UUID NOT NULL REFERENCES mpp.instansi(id) ON DELETE CASCADE,

    name VARCHAR(255) NOT NULL,
    description TEXT,
    estimasi_durasi_menit INT NOT NULL DEFAULT 10 CHECK (estimasi_durasi_menit > 0),
    requires_fo_verification BOOLEAN DEFAULT FALSE,

    -- Optional default follow-up service for "second service"
    default_second_service_id UUID REFERENCES mpp.jenis_layanan(id) ON DELETE SET NULL,

    is_active BOOLEAN DEFAULT TRUE,

    -- Audit
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by UUID,
    deleted_at TIMESTAMPTZ,
    deleted_by UUID
);

CREATE INDEX idx_jenis_layanan_instansi_id ON mpp.jenis_layanan(instansi_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_jenis_layanan_is_active ON mpp.jenis_layanan(is_active) WHERE deleted_at IS NULL;

-- -----------------------------------------------------
-- SYARAT_DOKUMEN (Document requirement)
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS mpp.syarat_dokumen (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jenis_layanan_id UUID NOT NULL REFERENCES mpp.jenis_layanan(id) ON DELETE CASCADE,

    name VARCHAR(255) NOT NULL,
    is_required BOOLEAN DEFAULT TRUE,
    notes TEXT,
    sort INT DEFAULT 0,

    -- Audit
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by UUID
);

CREATE INDEX idx_syarat_dokumen_jenis_layanan_id ON mpp.syarat_dokumen(jenis_layanan_id);

-- -----------------------------------------------------
-- LOKET (Counter)
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS mpp.loket (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instansi_id UUID NOT NULL REFERENCES mpp.instansi(id) ON DELETE CASCADE,

    code VARCHAR(50) NOT NULL,           -- e.g. 'L-A1'
    name VARCHAR(120),                   -- e.g. 'Loket 1'
    status mpp.loket_status NOT NULL DEFAULT 'CLOSED',
    last_idle_at TIMESTAMPTZ DEFAULT NOW(),  -- drives idle-longest allocation

    is_active BOOLEAN DEFAULT TRUE,

    -- Audit
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by UUID,
    deleted_at TIMESTAMPTZ,
    deleted_by UUID
);

CREATE UNIQUE INDEX loket_instansi_code_key ON mpp.loket(instansi_id, code) WHERE deleted_at IS NULL;
CREATE INDEX idx_loket_instansi_id ON mpp.loket(instansi_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_loket_status ON mpp.loket(status) WHERE deleted_at IS NULL;

-- -----------------------------------------------------
-- LOKET_LAYANAN (Counter <-> Service map)
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS mpp.loket_layanan (
    loket_id UUID NOT NULL REFERENCES mpp.loket(id) ON DELETE CASCADE,
    jenis_layanan_id UUID NOT NULL REFERENCES mpp.jenis_layanan(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (loket_id, jenis_layanan_id)
);

CREATE INDEX idx_loket_layanan_jenis_layanan_id ON mpp.loket_layanan(jenis_layanan_id);

-- -----------------------------------------------------
-- LOKET_SESSION (Operator shift)
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS mpp.loket_session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    loket_id UUID NOT NULL REFERENCES mpp.loket(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES core.users(id) ON DELETE CASCADE,

    opened_at TIMESTAMPTZ DEFAULT NOW(),
    closed_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Only one active (open) session per loket at a time
CREATE UNIQUE INDEX loket_session_active_key ON mpp.loket_session(loket_id)
    WHERE is_active = TRUE AND closed_at IS NULL;
CREATE INDEX idx_loket_session_user_id ON mpp.loket_session(user_id);
CREATE INDEX idx_loket_session_loket_id ON mpp.loket_session(loket_id);

-- =====================================================
-- COMMENTS
-- =====================================================
COMMENT ON TABLE mpp.instansi IS 'Agencies operating in the MPP; owns services and lokets';
COMMENT ON COLUMN mpp.instansi.prefix IS 'Single/short prefix used in queue numbers, e.g. A -> A-014';
COMMENT ON COLUMN mpp.instansi.queue_mode IS 'FIFO or BOOKING_PRIORITY, configurable per agency';
COMMENT ON TABLE mpp.jenis_layanan IS 'Service types under an agency; carry requirements and duration';
COMMENT ON COLUMN mpp.jenis_layanan.default_second_service_id IS 'Suggested follow-up service (second service)';
COMMENT ON TABLE mpp.syarat_dokumen IS 'Document requirements attached to a service';
COMMENT ON TABLE mpp.loket IS 'Physical service counters';
COMMENT ON COLUMN mpp.loket.last_idle_at IS 'Timestamp used by the idle-longest allocation rule';
COMMENT ON TABLE mpp.loket_layanan IS 'Which services a loket can serve (many-to-many)';
COMMENT ON TABLE mpp.loket_session IS 'Binds a logged-in operator to a loket for a shift';
