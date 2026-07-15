# Project Stack

## Scope and detection date

This is a backend-only, multi-module Go monorepo. Each service owns its `go.mod` and `go.sum`; no root `go.work` exists. Run Go commands from the target service directory.

Evidence: root `docker-compose.yml`; service manifests; `AGENTS.md`; `.opencode/docs/SERVICE_MATRIX.md`.

## Services and runtimes

| Service/module | Go version | HTTP/data stack | Notes |
|---|---:|---|---|
| `system` | 1.18 | Fiber, GORM/Postgres, MongoDB, Redis | Compose port 9001 |
| `master` | 1.20 | Fiber, sqlx, pgx/lib/pq, MongoDB, RabbitMQ | Compose port 9002; SQLMock available |
| `inventory` | 1.18 | Fiber, GORM/Postgres, MongoDB | Compose port 9003; Testify available |
| `sales` | 1.23.0, toolchain 1.24.6 | Fiber, GORM/Postgres, MongoDB, RabbitMQ, gocron | Compose port 9004 |
| `finance` | 1.18 | Fiber, GORM/Postgres, MongoDB, decimal | Compose port 9005 |
| `tms` | 1.21, toolchain 1.22.1 | Fiber, GORM/Postgres, Viper, Swagger | Compose port 9006; golang-migrate |
| `mobile` | 1.20 | Fiber, GORM/sqlx/Postgres, Redis, MongoDB | Compose port 9008 |
| `pjp` | 1.21, toolchain 1.22.1 | Gin, GORM/Postgres, Viper, Swagger | Compose port 9010; golang-migrate |
| `cronjob` | 1.23.5 | Fiber, GORM/Postgres, Redis, MongoDB, Resty, gocron | Compose port 9100 |
| `pjp-principle` | 1.21, toolchain 1.22.1 | Gin, GORM/Postgres, Viper, Swagger | golang-migrate |
| `pjp-sales` | 1.23.0, toolchain 1.24.6 | Fiber, GORM/Postgres, MongoDB, RabbitMQ, gocron | parallel sales module |

## Shared architecture and packages

- Default Fiber pattern: Controller → Service → Repository → DB. Controllers do not directly call repositories; repository layer does not hold business logic.
- Fiber response envelope is service-local `pkg/responsebuild/response.go`; preserve the target service’s established response/paging style.
- Fiber JWT group middleware follows `pkg/middleware/jwt_middleware.go`; `master/pkg/middleware/jwt_middleware.go` uses `github.com/gofiber/jwt/v2`.
- PJP modules are Gin-based. Do not copy Fiber handler/middleware patterns into them.
- PostgreSQL is dominant. `master` prefers sqlx/pgx/lib-pq; most other Fiber services use GORM. Inspect target module before selecting query style.

## Runtime and dependencies

- Root compose manages Redis, RabbitMQ, and service containers on `scylla-network`.
- Compose services mount local source and shared Go caches. Main services target host-local Postgres via `host.docker.internal:5432`.
- Local runtime contract: target DB is `ggn_scyllax`; do not default to remote/staging DBs.
- Service `.env` files exist and may contain plaintext secrets. Read only needed variable names; never print, copy, stage, or commit values.
- `rtk` is installed at `/opt/homebrew/bin/rtk` in this environment. Repository policy requires `rtk`-prefixed compose/Go commands; use the direct underlying command only when diagnosing RTK itself and record why.

## Build, test, lint, and CI

- Per-service baseline: `rtk go test ./...` and `rtk go build ./...` from the module directory.
- Targeted tests use standard Go test selection: `rtk go test ./path/to/pkg -run TestName -v`.
- No root `go.work`, root Makefile, justfile, GitHub Actions, GitLab CI file, or golangci-lint configuration was detected during harness initialization.
- Unit-test patterns include `net/http/httptest`, `DATA-DOG/go-sqlmock` (notably `master`), and `testify/require` (notably `inventory`). Reuse local test helpers before adding a library.

## Migrations and generators

- `tms/migrations/` uses `golang-migrate` through `tms/Makefile`.
- `pjp/database/migrate/` and `pjp-principle/database/migrate/` use `golang-migrate` through their Makefiles.
- No migration directory was detected in `master`; do not invent migration tooling for it.
- Swagger generation is available only in `tms`, `pjp`, and `pjp-principle` via `swag init` Makefile target.

## Source strategy and staleness

- Repo-local authority: target service `go.mod`, `Makefile`, `Dockerfile`, `main.go`, `pkg/`, tests, root `docker-compose.yml`, and `.opencode/docs/ARCHITECTURE.md`.
- Several service READMEs are generic GitLab templates; do not treat them as runtime authority without corroboration.
- Go versions vary from 1.18 to 1.23.5, with declared toolchains in selected services. Re-check `go.mod` after toolchain upgrades or module changes.
- For version-sensitive Fiber, Gin, GORM, sqlx, golang-migrate, Swagger, or Go toolchain behavior not settled here, use official docs through `@librarian`/Context7 before implementation.
