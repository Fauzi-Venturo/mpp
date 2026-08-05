# Memory Index — MPP (Sistem Antrian Mal Pelayanan Publik)

Mirror dari `~/.claude/memory/mpp/`. Sumber kebenaran = file in-repo ini.

- [MPP — Overview & status fase](project_mpp_overview.md) — apa produknya, dokumen di `docs/`, Fase 0 selesai & kode `mpp` belum ada
- [MPP — Batasan yang mengikat](project_mpp_constraints.md) — tenancy 1 gedung=1 company, Redis wajib, UTC, no ORM, Yarn 1, FE tanpa auth, TTS pre-rendered, atomisitas nomor/kuota
- [TDD + 3 gate tiap slice](feedback_mpp_tdd_gate.md) — test merah dulu lalu implementasi; plan mode · ≥1 subagent · test hijau → `/grademe` → commit; gate = `make check` (api-test + Vitest + tsc)
- [Engine dev mesin ini](reference_mpp_dev_engine.md) — Postgres compose di port 5433, psql keg-only, Node 25 vs .nvmrc 22, cache Yarn korup, Colima
- [CLAUDE.md di apps/ masih basi](feedback_mpp_stale_claudemd.md) — ikut ter-vendor dari skeleton asal; sumber kebenaran ada di `docs/`
