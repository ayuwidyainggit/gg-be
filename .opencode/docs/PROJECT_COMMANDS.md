# Project Commands

All commands assume shell workflows are `rtk`-prefixed per project policy. RTK is installed at `/opt/homebrew/bin/rtk`; its `go`/`docker`/`make` shims forward to native commands. Use raw `go`/`docker` only when debugging RTK itself; record why in evidence.

## Common posture

- `rtk` is mandatory prefix. Do not run plain `go test` or plain `docker compose` for normal workflows.
- Run Go commands from the target service directory; there is no `go.work`.
- Stop on first failure, capture log under `.opencode/evidence/<task-id>/`.
- Use environment variable substitution from each service's own `.env`; never copy values out.

## Safe (default) commands

| Task | Command |
|---|---|
| Inspect compose status | `rtk docker compose -f docker-compose.yml ps` |
| Bring up stack (host-local DB only) | `rtk docker compose -f docker-compose.yml up -d` |
| Stop stack | `rtk docker compose -f docker-compose.yml stop` |
| Module sync | `cd <service> && rtk go mod download && rtk go mod tidy` |
| Run service tests | `cd <service> && rtk go test ./...` |
| Targeted test | `cd <service> && rtk go test ./pkg/... -run TestName -v` |
| Build service | `cd <service> && rtk go build ./...` |
| Run a single service | `cd <service> && rtk go run main.go` |
| Hot reload (services with `air`) | `cd <service> && rtk make dev-reload` |
| Open Swagger docs (`tms`/`pjp`/`pjp-principle`) | `cd <service> && rtk make doc` |
| DB clone (from remote) | `./scripts/clone_db.sh [service]` |
| DB restore from local dump | `./scripts/restore_db.sh [dump] [db]` |
| Generate env template | `./scripts/generate-env.sh` |
| Install Citus extension locally | `./scripts/install_citus.sh` |

## Migration commands (services with golang-migrate only)

`master`, `system`, `inventory`, `sales`, `finance`, `mobile`, `cronjob`, and `pjp-sales` do not have Makefile migration targets. Run migrations only on services that ship them:

- `tms` (`tms/Makefile`): `cd tms && rtk make migration <name>`, `cd tms && rtk make migrateUp`, `cd tms && rtk make migrateDown`, `cd tms && rtk make migrateForce <v>`, `cd tms && rtk make migrateDrop`.
- `pjp` and `pjp-principle`: identical Makefile targets; migration directory is `database/migrate/`.

## Validation commands

Per-service regression command:

```bash
cd <service> && rtk go mod download && rtk go mod tidy && rtk go test ./... && rtk go build ./...
```

Cross-service sweep (when the slice is repo-wide and compose is up):

```bash
for svc in system master inventory sales finance tms mobile pjp cronjob; do
  ( cd "$svc" && rtk go test ./... ) || exit 1
done
```

Service-specific validation patterns live in `.opencode/docs/QUALITY.md`.

## Commands requiring approval

- `rtk docker compose down -v` — destroys compose volumes including Redis and RabbitMQ persisted data.
- `rtk docker system prune` or `rtk docker volume prune` — global destructive.
- Any `migrateDrop` against a non-local environment.
- `./scripts/sync_remote_to_local.sh --drop --yes` to any host other than `localhost`/`127.0.0.1` or to any DB other than `ggn_scyllax` (requires `--allow-nonlocal-target` / `--allow-custom-target`).
- Any commit, push, or PR command without explicit user instruction.
- Re-seeding or destructive local DB writes outside the active task.

## Destructive and prod-risk commands (forbidden without approval)

- `rtk docker compose down -v`, `rtk docker volume rm`, `rtk docker image rm -f`.
- `migrateDrop`, `pg_restore --clean` against non-local DBs, or any DROP DATABASE outside `ggn_scyllax`.
- `--no-verify`, `--no-gpg-sign`, force-push, `amend`, or any other git bypass flag.
- `./scripts/clone_db.sh` pointed at non-target service directory; reroute to `./scripts/sync_remote_to_local.sh --preflight-only` first.
- Editing or staging `.env`, `*.pem`, `*.key`, or any committed secret material.

## Env and local DB notes

- Default local Postgres is `host.docker.internal:5432` (from compose) with credentials `postgres/postgres` and database `ggn_scyllax`.
- The `clone_db.sh` script reads each service's `.env` for `DB_HOST`/`DB_PORT`/`DB_USER`/`DB_PASS`/`DB_NAME`; default local target is `localhost:5432` with `postgres` user. Override via `LOCAL_DB_HOST`/`LOCAL_DB_PORT`/`LOCAL_DB_USER`/`LOCAL_DB_NAME` when needed.
- `sync_remote_to_local.sh` requires `SOURCE_DB_PASSWORD` and `LOCAL_DB_PASSWORD` from environment; never hard-code.
- Service `.env` files contain real secrets and are not safe to copy, log, or commit.
- Compose-level credentials should be treated as read-only infrastructure data; document the names but never paste values.
