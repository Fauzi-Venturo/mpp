# ============================================================================
# MPP Monorepo — root orchestration (lightweight meta-repo)
# ----------------------------------------------------------------------------
# Delegates to apps/api (Go) and apps/web (Next.js / Yarn 1). This root does not
# unify the two dependency graphs; each app keeps its own toolchain.
# Run `make help` for the full list.
# ============================================================================

API_DIR := apps/api
WEB_DIR := apps/web

.DEFAULT_GOAL := help

## ── Setup ───────────────────────────────────────────────────────────────────
.PHONY: bootstrap
bootstrap: ## Install deps for both apps (Go modules + Yarn) and copy env files
	@echo "→ backend: go mod download"
	cd $(API_DIR) && go mod download
	@echo "→ frontend: yarn install"
	cd $(WEB_DIR) && yarn install --frozen-lockfile
	@test -f .env || cp .env.example .env
	@test -f $(API_DIR)/.env || cp $(API_DIR)/.env.example $(API_DIR)/.env
	@test -f $(WEB_DIR)/.env || cp $(WEB_DIR)/.env.example $(WEB_DIR)/.env
	@echo "✓ bootstrap complete — edit .env files, then: make up && make db-setup"

## ── Infra (docker compose) ──────────────────────────────────────────────────
.PHONY: up
up: ## Start infra only (postgres + redis)
	docker compose up -d postgres redis

.PHONY: up-full
up-full: ## Start the full stack (postgres + redis + api + web)
	docker compose --profile full up -d --build

.PHONY: down
down: ## Stop all compose services
	docker compose down

.PHONY: logs
logs: ## Tail compose logs
	docker compose logs -f

## ── Backend (apps/api) ──────────────────────────────────────────────────────
.PHONY: api-dev
api-dev: ## Run backend with hot-reload (air)
	cd $(API_DIR) && make dev

.PHONY: api-build
api-build: ## Build backend binary
	cd $(API_DIR) && make build

.PHONY: api-test
api-test: ## Run backend tests
	cd $(API_DIR) && go test ./...

.PHONY: db-setup
db-setup: ## Run backend migrations + seeders (core, then mpp)
	cd $(API_DIR) && make db-setup

.PHONY: db-reset
db-reset: ## Drop, recreate, migrate and seed the database
	cd $(API_DIR) && make db-reset

.PHONY: migrate-up
migrate-up: ## Apply migrations (pass MODULE=mpp for the MPP domain)
	cd $(API_DIR) && make migrate-up MODULE=$(or $(MODULE),core)

## ── Frontend (apps/web) ─────────────────────────────────────────────────────
.PHONY: web-dev
web-dev: ## Run frontend dev server (port 8002)
	cd $(WEB_DIR) && yarn dev

.PHONY: web-build
web-build: ## Production build of the frontend
	cd $(WEB_DIR) && yarn build

.PHONY: web-lint
web-lint: ## Lint the frontend
	cd $(WEB_DIR) && yarn lint

.PHONY: web-check
web-check: ## Type-check the frontend (tsc --noEmit)
	cd $(WEB_DIR) && yarn tsc:check

## ── Meta ────────────────────────────────────────────────────────────────────
.PHONY: check
check: api-test web-check ## Run backend tests + frontend type-check

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'
