# Project Detected Tools

Detected via manifest/file inspection during `/init-harness`. Each entry lists the repo evidence used.

## Language/runtime

- Go, multi-module monorepo. Versions vary 1.18–1.23.5 across services; some declare explicit `toolchain` directives (`sales` 1.24.6, `pjp-sales` 1.24.6, `tms`/`pjp`/`pjp-principle` 1.22.1). Evidence: each service `go.mod`.
- No root `go.work`. Evidence: absence check at repo root.

## HTTP frameworks

- Fiber v2 (`gofiber/fiber/v2`): `system`, `master`, `inventory`, `sales`, `finance`, `tms`, `mobile`, `cronjob`, `pjp-sales`. Evidence: respective `go.mod` require blocks.
- Gin v1.9.1: `pjp`, `pjp-principle`. Evidence: respective `go.mod` require blocks.

## Data layer

- GORM + `gorm.io/driver/postgres`: `system`, `inventory`, `sales`, `finance`, `tms`, `mobile`, `pjp`, `pjp-principle`, `cronjob`, `pjp-sales`.
- sqlx + pgx/lib-pq (no ORM): `master`.
- MongoDB driver present in most services alongside Postgres (dual-store pattern); confirm per-service before assuming Mongo is active.
- Redis (`go-redis/redis/v8`): `system`, `mobile`, `cronjob`.
- RabbitMQ (`amqp`): `sales`, `master`, `pjp-sales`.

## Auth

- `golang-jwt/jwt/v4` + `gofiber/jwt/v2`: Fiber services.
- `golang-jwt/jwt/v5`: `tms`.
- `golang-jwt/jwt` v3 with custom Gin decode helpers (`pjp/utils/decode.go`, `pjp-principle/utils/decode.go`): Gin services.

## Migration tooling

- `golang-migrate` CLI via Makefile targets: `tms` (`tms/migrations/`), `pjp` and `pjp-principle` (`database/migrate/`).
- No migration tool detected for `system`, `master`, `inventory`, `sales`, `finance`, `mobile`, `cronjob`, `pjp-sales`.

## API docs / generators

- `swag` (Swagger) via `make doc`: `tms`, `pjp`, `pjp-principle`.
- `air` (hot reload) via `make dev-reload`: `tms`, `pjp`, `pjp-principle`.
- No detected scaffolding CLI (no Artisan/Rails-generator/Nest-CLI equivalent) for Fiber-only services; controller/service/repository files are hand-authored.

## Test tooling

- `DATA-DOG/go-sqlmock`: declared in `master/go.mod` and `pjp/go.mod`; used directly in `master` repository/service tests.
- `testify`: declared and used in `inventory`.
- `net/http/httptest`: used across Fiber controller tests (`master`, `sales`, `finance`, `pjp-sales`, others).
- No project-wide CI runner detected (no `.github/workflows/`, no `.gitlab-ci.yml`).
- No root or per-service `golangci-lint` config detected.

## Infra tooling

- Docker: every service has its own multi-stage `Dockerfile` (Go builder + Alpine runtime), except `tms`, which runs from a shared `golang:1.23.0-alpine` dev image with a bind mount in compose.
- Docker Compose: root `docker-compose.yml` orchestrates all services plus `redis` and `rabbitmq` on network `scylla-network`, with shared Go module/build cache volumes.
- `rtk`: present at `/opt/homebrew/bin/rtk` in this environment; required command prefix per `AGENTS.md`. `.rtk/filters.toml` exists as local RTK config (no secrets).

## Repo automation scripts

- `scripts/clone_db.sh`, `scripts/clone_staging.sh`, `scripts/restore_db.sh`, `scripts/install_citus.sh`, `scripts/sync_remote_to_local.sh`, `scripts/generate-env.sh`. Documented in `scripts/README.md`.

## Multi-tool rules files (harmonization candidates)

- `CLAUDE.md` exists at repo root (Claude Code rules file).
- `.cursorrules` exists at repo root (Cursor rules file).
- No `.codex/`, `.windsurfrules`, or similar detected.
- Harmonization via `rules-source-scanner.py`/`rules-harmonizer.py` was not run in this pass; run it in a follow-up task if `CLAUDE.md`/`.cursorrules` content should be reconciled into `.opencode/docs/SOURCE_RULES.md`.

## Explicit gaps (do not assume presence)

- No frontend/mobile-app UI codebase detected; this is backend-API-only. `DESIGN.md` in this repo therefore documents API/response/error surface conventions, not visual UI.
- No CI pipeline config detected; do not claim CI gates exist.
- No lint config detected; do not assume `golangci-lint` runs anywhere in this repo.
