# Quality

- Validate in the target service directory, not at repo root.
- Common checks: `rtk go mod download && rtk go mod tidy`, `rtk go test ./...`, and targeted `rtk go test ./path/to/pkg -run TestName`.
- Runtime verification starts with `rtk docker compose -f docker-compose.yml ps`; bring services up with `rtk docker compose -f docker-compose.yml up -d` when needed.
- Migration validation is service-specific: `tms` uses `rtk make migrateUp` from `migrations/`; `pjp` and `pjp-principle` use `database/migrate/`. Do not assume every service has a Makefile.
- Evidence for material work should reference changed files, commands run, and remaining risk.
- Done means repo conventions, tenant rules, transaction rules, and relevant tests/checks all align.

## Per-module validation focus

| Module group | Baseline checks | Extra checks/caveats |
| --- | --- | --- |
| Compose-managed Fiber services (`system`, `master`, `inventory`, `sales`, `finance`, `mobile`, `cronjob`, `tms`) | `rtk go mod download && rtk go mod tidy`; `rtk go test ./...` | Confirm runtime port/env mapping against `docker-compose.yml`; for API smoke checks, Fiber health endpoint is typically `GET /ping`. |
| PJP family (`pjp`, `pjp-principle`) | `rtk go mod download && rtk go mod tidy`; `rtk go test ./...` | Gin-oriented stack and migration flow differs from Fiber modules. |
| Extra module (`pjp-sales`) | `rtk go mod download && rtk go mod tidy`; `rtk go test ./...` | Not in root compose defaults; validate with module-local runtime conventions when touched. |

## Per-service command cheat sheet

| Module | Enter module | Baseline run/test | Hot reload / helper | Migration command |
| --- | --- | --- | --- | --- |
| `system` | `cd system` | `rtk go mod download && rtk go mod tidy && rtk go test ./...` | `rtk go run main.go` or compose | none documented |
| `master` | `cd master` | `rtk go mod download && rtk go mod tidy && rtk go test ./...` | `rtk go run main.go` or compose | none documented |
| `inventory` | `cd inventory` | `rtk go mod download && rtk go mod tidy && rtk go test ./...` | `rtk go run main.go` or compose | none documented |
| `sales` | `cd sales` | `rtk go mod download && rtk go mod tidy && rtk go test ./...` | `rtk go run main.go` or compose | none documented |
| `finance` | `cd finance` | `rtk go mod download && rtk go mod tidy && rtk go test ./...` | `rtk go run main.go` or compose | none documented |
| `tms` | `cd tms` | `rtk go mod download && rtk go mod tidy && rtk go test ./...` | `rtk make dev`, `rtk make dev-reload`, `rtk make doc` | `rtk make migrateUp` |
| `mobile` | `cd mobile` | `rtk go mod download && rtk go mod tidy && rtk go test ./...` | `rtk go run main.go` or compose | none documented |
| `pjp` | `cd pjp` | `rtk go mod download && rtk go mod tidy && rtk go test ./...` | `rtk make dev`, `rtk make dev-reload`, `rtk make doc` | `rtk make migrateUp` |
| `cronjob` | `cd cronjob` | `rtk go mod download && rtk go mod tidy && rtk go test ./...` | `rtk go run main.go` or compose | none documented |
| `pjp-principle` | `cd pjp-principle` | `rtk go mod download && rtk go mod tidy && rtk go test ./...` | `rtk make dev`, `rtk make dev-reload`, `rtk make doc` | `rtk make migrateUp` |
| `pjp-sales` | `cd pjp-sales` | `rtk go mod download && rtk go mod tidy && rtk go test ./...` | `rtk go run main.go` | none documented |

Env caveat for PJP family:
- `pjp` and `pjp-principle` Makefiles reference `development.env`; confirm whether local execution is using `.env`, `development.env`, or compose/env-file wiring before validating runtime behavior.

Compose-first checks from repo root:

```bash
rtk docker compose -f docker-compose.yml ps
rtk docker compose -f docker-compose.yml up -d
```

## Repo-local doc reliability notes

- Treat module READMEs as advisory only; many are template placeholders and may not represent current run/test/migration behavior.
- Prefer this precedence when data conflicts: `docker-compose.yml` → module `go.mod` / `Makefile` / env files → README.
- Port, env, and migration expectations should be re-verified per touched module before declaring done.
