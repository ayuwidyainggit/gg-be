# Framework Playbook

## Default rule: generator/CLI first

Use target module commands and existing generators before manual framework artifact creation. This applies to existing services too. Manual artifact edits are allowed only when:

1. required tool is unavailable or not permitted;
2. tool failed and output is recorded in evidence;
3. repository intentionally avoids that generator;
4. task customizes an existing generated file; or
5. user explicitly requests manual edits.

For each manual fallback, record attempted/skipped command and reason in `.opencode/evidence/<task-id>/`.

## Go module workflow

1. Identify target service from its own `go.mod`; no root `go.work` exists.
2. Read target service `main.go`, controller/service/repository package, tests, and local `pkg/` helpers before adding code.
3. Use `cd <service> && rtk go mod download && rtk go mod tidy` only if dependencies changed or module resolution needs repair.
4. Keep changes inside one module unless active plan explicitly spans multiple modules.
5. Run `cd <service> && rtk go test ./...` then `rtk go build ./...` before a completion claim.

## HTTP framework placement

### Fiber services

Services: `system`, `master`, `inventory`, `sales`, `finance`, `tms`, `mobile`, `cronjob`, `pjp-sales`.

- Register routes in target controller route/group method; preserve middleware order and prefix behavior from `main.go`.
- Use service-local `pkg/responsebuild` and existing entity pagination/response types instead of inventing envelopes.
- Put request parsing/validation in controllers; business policy in services; persistence/query code in repositories.
- Use target service's existing DB style: `master` uses sqlx/pgx/lib-pq patterns; GORM services use existing GORM patterns.
- Reuse `pkg/middleware/JWTProtected()` pattern for protected route groups. Do not weaken it, log bearer tokens, or substitute route scope without an explicit plan decision.

### Gin services

Services: `pjp`, `pjp-principle`.

- Use Gin router, handler, and utility conventions already in target service.
- Do not bring Fiber `*fiber.Ctx`, Fiber middleware, or response helper APIs into Gin services.
- Preserve GORM/Viper/Swagger conventions from the module.

## Database and transactions

- Apply `.opencode/docs/ARCHITECTURE.md` tenant/schema rules before writing SQL or GORM queries.
- Controller → Service → Repository → DB is mandatory.
- Write operations belong in service-layer transactions. Repository writes must honor transaction context extraction.
- Parameterize all user-controlled values. Dynamic sort/order identifiers must come from closed Go allowlists, never raw request strings.
- Do not add migrations to a module without an existing migration practice unless user-approved design explicitly establishes it.

## Migrations (golang-migrate services only)

`tms`, `pjp`, and `pjp-principle` have Makefile-managed golang-migrate workflows.

```bash
cd tms && rtk make migration <name>
cd tms && rtk make migrateUp
cd tms && rtk make migrateDown
```

- Generated migration pairs are source artifacts. Review both `.up.sql` and `.down.sql`.
- `migrateForce` and `migrateDrop` are destructive; require explicit approval and local target confirmation.
- Never run migrations against remote/staging DB by default.

## Swagger and development helpers

- `tms`, `pjp`, `pjp-principle`: generate Swagger via `rtk make doc` (`swag init`). Re-run only after API-doc-relevant changes and include generated diff in plan boundary.
- Their `rtk make dev`, `rtk make dev-reload`, and `rtk make install` targets are the first option over hand-written shell sequences.
- Other services have no detected framework scaffold/generator. Follow local package patterns; do not fabricate a generator workflow.

## Tests

- Start Red → Green → Refactor for material behavior.
- Controller tests: reuse `net/http/httptest` patterns in target service.
- SQL query behavior in `master`: use `DATA-DOG/go-sqlmock` + `sqlx.NewDb` and assert query arguments/expectations.
- `inventory` already uses Testify; reuse it in that service instead of introducing a competing assertion package.
- Test only target module by default. Cross-service sweep is reserved for explicit multi-service changes.

## Docker and runtime

- Begin compose work at root: `rtk docker compose -f docker-compose.yml ps`.
- Compose expects host-local PostgreSQL at `host.docker.internal:5432`; do not replace it with a remote default.
- Root compose is development-oriented with bind mounts and shared Go cache volumes. Inspect service port/prefix before curl smoke tests.
- Runtime proof is env-gated. Missing DB/token/service state means report `not-ready`, not `passed`.

## Official docs and manual fallback

- Local docs settle repository conventions. For version-sensitive Go/Fiber/Gin/GORM/sqlx/golang-migrate/Swagger behavior, use official docs via `@librarian`/Context7 before coding.
- Do not infer generator flags, GORM migration semantics, Fiber/Gin API behavior, or Go toolchain compatibility from memory.
- If a documented command conflicts with local `go.mod` version or actual command output, stop and record the conflict; do not silently change stack assumptions.
