# packages/api-contract

Placeholder for the **shared FE↔BE contract** — the single source of truth for the
API surface that both `apps/api` (Go) and `apps/web` (TypeScript) agree on.

Today the coupling is implicit:
- Backend documents endpoints under `apps/api/docs/api-contract/`.
- Frontend declares them in `apps/web/src/lib/api/endpoints.ts` and unwraps the
  standard response envelope `{ data, message, meta, errors }` in
  `apps/web/src/lib/api/client.ts`.

This package is where that contract can be made explicit and type-safe as the MPP
domain grows. Options (decide during the API-design phase — see
`docs/04-api/api-conventions.md`):

- **OpenAPI spec** committed here, with generated Go server stubs and a generated
  TypeScript client.
- **Hand-maintained TS types** mirroring the DTOs, imported by `apps/web`.

Kept intentionally empty for now so the monorepo layout reserves the seam without
prescribing the tooling.
