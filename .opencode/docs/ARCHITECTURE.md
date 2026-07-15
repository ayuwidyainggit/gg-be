# Architecture

- Scylla Backend is a multi-module Go monorepo; each service owns its own `go.mod` and runtime.
- Root `docker-compose.yml` manages `system`, `master`, `inventory`, `sales`, `finance`, `tms`, `mobile`, `pjp`, `cronjob`, and `redis`.
- Compose DB runtime baseline is host local Postgres via Docker Desktop target `host.docker.internal:5432` using `postgres/postgres` and DB `ggn_scyllax` (avoid remote dev DB literals in compose defaults).
- `pjp` is Gin-based; the main non-PJP services are Fiber-oriented. Many Fiber services expose `GET /ping` as the health check pattern.
- Layering is strict: Controller → Service → Repository → DB.
- Controllers do request parsing/validation and response shaping; services own business logic and transactions; repositories own data access only.
- Write operations must run in service-layer transactions and repository writes must use tx-aware DB access.
- Multi-tenant rules matter: keep `cust_id` filters, use `parent_cust_id` for parent-company master data, and `custId` for distributor-specific transactional data.
- Schema prefixes are semantically important: `inv.`, `mst.`, `acf.`, `sys.`, `smc.`, `report.`.
- Prefer `Take()` for single-row fetches when not-found behavior matters.

## Service/module inventory source

For service/module inventory details, use `.opencode/docs/SERVICE_MATRIX.md` as canonical source:

- compose presence
- runtime style + default ports
- env/Makefile/migration footprint
- README authority audit status (`authoritative` / `advisory` / `stale-template` / `missing`)

## Practical caveats

- Several service READMEs are stale/incomplete; use README authority status in `SERVICE_MATRIX.md` and prefer `docker-compose.yml`, module `go.mod`, module env files, and module `Makefile` as operational truth.
- Env key conventions are not fully uniform across modules:
  - Most services use `SERVER_PORT` + `DB_*`.
  - `tms` and `pjp` also use `PORT` and `POSTGRES_*` families.
- Ports above are the current compose defaults and should be treated as runtime defaults, not cross-environment guarantees.
- Migration workflows differ by module; do not assume one migration command or folder pattern applies repo-wide.
