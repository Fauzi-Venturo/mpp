-- =====================================================
-- MPP SCHEMA - SYSTEM CONFIG
-- =====================================================
-- Migration 005: global & per-agency configuration
-- =====================================================

CREATE TABLE IF NOT EXISTS mpp.system_config (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- NULL instansi_id = global config; otherwise per-agency override
    instansi_id UUID REFERENCES mpp.instansi(id) ON DELETE CASCADE,

    config_key VARCHAR(120) NOT NULL,    -- e.g. 'number_format', 'checkin_window', 'tts_text', 'reminder'
    config_value JSONB NOT NULL DEFAULT '{}',
    description TEXT,

    -- Audit
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by UUID
);

-- One value per key per scope (global vs specific agency)
CREATE UNIQUE INDEX system_config_agency_key ON mpp.system_config(instansi_id, config_key)
    WHERE instansi_id IS NOT NULL;
CREATE UNIQUE INDEX system_config_global_key ON mpp.system_config(config_key)
    WHERE instansi_id IS NULL;
CREATE INDEX idx_system_config_instansi_id ON mpp.system_config(instansi_id);

COMMENT ON TABLE mpp.system_config IS 'Global (instansi_id NULL) and per-agency configuration values';
COMMENT ON COLUMN mpp.system_config.config_key IS 'e.g. number_format, checkin_window, tts_text, reminder, second_service_priority';
