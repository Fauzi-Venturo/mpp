-- =====================================================
-- MPP SEEDER - TENANT COMPANY
-- =====================================================
-- Seeder 002: the MPP building as a core.company (tenant).
-- One MPP building = one company (see docs/02-domain/domain-model.md).
--
-- References existing core seed data:
--   owner_id  = 10000000-...-001 (user "owner", seeder core/002_users)
--   client_id = 40000000-...-001 (client "owner", seeder core/003_clients)
-- =====================================================

INSERT INTO core.companies (id, parent_id, name, type, owner_id, client_id, sort, is_active) VALUES
(
    'a1000000-0000-0000-0000-000000000001',
    NULL,
    'Mall Pelayanan Publik',
    'holding',
    '10000000-0000-0000-0000-000000000001',
    '40000000-0000-0000-0000-000000000001',
    100,
    TRUE
)
ON CONFLICT DO NOTHING;
