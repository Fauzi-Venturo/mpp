# Dokumentasi MPP — Sistem Antrian Mall Pelayanan Publik

Dokumentasi requirement lengkap untuk membangun sistem antrian MPP. Dokumen bisnis
ditulis dalam **Bahasa Indonesia**, dokumen teknis dalam **English** (bilingual).

## Urutan baca yang disarankan

1. **Overview** — pahami produk & istilah
   - [`00-overview/product-brief.md`](./00-overview/product-brief.md) — visi, tujuan, ruang lingkup _(ID)_
   - [`00-overview/glossary.md`](./00-overview/glossary.md) — kamus istilah domain _(ID/EN)_
   - [`00-overview/battle-w04-mpp-relay.md`](./00-overview/battle-w04-mpp-relay.md) — aturan & kontrak Battle #1 MPP Relay (w04) _(ID)_
2. **Requirements** — apa yang harus dibangun
   - [`01-requirements/prd.md`](./01-requirements/prd.md) — Product Requirements Document _(ID)_
   - [`01-requirements/functional-requirements.md`](./01-requirements/functional-requirements.md) — FR per modul (18 modul) _(ID)_
   - [`01-requirements/non-functional-requirements.md`](./01-requirements/non-functional-requirements.md) — NFR _(EN)_
   - [`01-requirements/user-roles-personas.md`](./01-requirements/user-roles-personas.md) — peran & persona _(ID)_
3. **Domain** — model data & aturan inti
   - [`02-domain/domain-model.md`](./02-domain/domain-model.md) — entities & relationships _(EN)_
   - [`02-domain/erd.md`](./02-domain/erd.md) — ERD (Mermaid) _(EN)_
   - [`02-domain/queue-state-machine.md`](./02-domain/queue-state-machine.md) — state machine antrian _(EN)_
   - [`02-domain/business-rules.md`](./02-domain/business-rules.md) — aturan bisnis _(ID)_
4. **Architecture** — bagaimana sistem disusun
   - [`03-architecture/system-architecture.md`](./03-architecture/system-architecture.md) _(EN)_
   - [`03-architecture/tech-stack.md`](./03-architecture/tech-stack.md) _(EN)_
   - [`03-architecture/data-flow.md`](./03-architecture/data-flow.md) — sequence diagrams _(EN)_
   - [`03-architecture/deployment.md`](./03-architecture/deployment.md) _(EN)_
5. **API**
   - [`04-api/api-conventions.md`](./04-api/api-conventions.md) _(EN)_
   - [`04-api/rest-endpoints.md`](./04-api/rest-endpoints.md) _(EN)_
   - [`04-api/websocket-events.md`](./04-api/websocket-events.md) _(EN)_
6. **Integrations**
   - [`05-integrations/whatsapp-ai-agent.md`](./05-integrations/whatsapp-ai-agent.md) _(EN)_
   - [`05-integrations/tts-voice-calling.md`](./05-integrations/tts-voice-calling.md) _(EN)_
   - [`05-integrations/kiosk-devices.md`](./05-integrations/kiosk-devices.md) _(EN)_
   - [`05-integrations/notifications.md`](./05-integrations/notifications.md) _(EN)_
7. **Security**
   - [`06-security/rbac-matrix.md`](./06-security/rbac-matrix.md) _(EN)_
   - [`06-security/security-privacy.md`](./06-security/security-privacy.md) _(EN)_
8. **UI/UX**
   - [`07-ui-ux/screens-inventory.md`](./07-ui-ux/screens-inventory.md) _(ID)_
   - [`07-ui-ux/wireframes.md`](./07-ui-ux/wireframes.md) _(ID)_
9. **Roadmap**
   - [`08-roadmap/delivery-plan.md`](./08-roadmap/delivery-plan.md) — fase pengerjaan _(ID)_

## Ringkasan sistem

MPP adalah sistem antrian untuk mal pelayanan publik yang menampung **banyak
instansi** dalam satu gedung. Warga mendaftar layanan melalui **WhatsApp (AI agent)**,
**website**, atau **walk-in** di lokasi, melakukan **check-in via QR** di kiosk,
menunggu dipanggil melalui **display TV dengan suara (TTS Bahasa Indonesia)**, lalu
dilayani di **loket** oleh petugas. **Front office** memverifikasi kelengkapan dokumen,
**supervisor** memantau, dan **admin** mengelola master data & konfigurasi. Semua
pergerakan nomor antrian bersifat **real-time** (WebSocket) dengan backend Go dan
Redis sebagai engine antrian.
