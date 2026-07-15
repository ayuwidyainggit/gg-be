# Tool Usage

## Authority order

1. Repo-local files and executable commands.
2. Project docs in `.opencode/docs/`.
3. Official documentation through `@librarian` or Context7 when library/API behavior is version-sensitive.
4. Upstream source/examples and GitHub search when local and official docs do not resolve the behavior.
5. Broad web search only when earlier sources are insufficient.

Never present an assumption as a repository fact.

## Local inspection

- Read target service `go.mod`, `main.go`, controller/service/repository files, local tests, Makefile, Dockerfile, and target `.env` variable names before changing a service.
- Use `read`, `grep`, and `glob` for source discovery. Use `rtk`-prefixed terminal commands for runtime/test/build commands.
- Treat `.env`, backups, and compose credentials as sensitive. Inspect only names/structure needed for a task; do not echo values.
- Use `rtk docker compose -f docker-compose.yml ps` before starting compose work.

## Tests, runtime, and evidence

- Unit/test/build path: `cd <service> && rtk go test ./... && rtk go build ./...`.
- Targeted test path: `cd <service> && rtk go test ./path/to/pkg -run TestName -v`.
- For API work, add an env-gated curl smoke only after compose/service/local DB are available. Never print bearer tokens.
- Persist material test/build/runtime output under `.opencode/evidence/<task-id>/`.
- Verify runtime claims with a command or curl; verify code claims with file evidence. A passing compile alone is not functional proof.

## Database and migrations

- Read `.opencode/docs/ARCHITECTURE.md` before query or transaction changes.
- Use local Postgres first. Never invoke a remote DB mutation by default.
- Only `tms`, `pjp`, and `pjp-principle` have detected golang-migrate commands; use `rtk make migration`, `migrateUp`, and `migrateDown` there.
- `migrateDrop`, `migrateForce`, restore, clone, and destructive sync commands require explicit user approval and local-target verification.

## MCP routing

- `@librarian`/Context7: current Go/Fiber/Gin/GORM/sqlx/golang-migrate/Swagger docs.
- `@explorer`: broad repository mapping, unknown ownership, cross-service patterns.
- `@backend`: bounded API/data/service implementation after plan/contract is clear.
- `@quality-gate`: material changes, security, auth, tenant isolation, DB safety, docs/config/prompt changes, or final completion claim.
- `browseros`: browser-only flows; normally not applicable to this backend-only repository.
- `semgrep`: targeted static/security pattern scan when auth, SQL, secrets, or unsafe input handling is touched.

## Generator-first rule

- Prefer target Makefile commands and established generators over manual framework artifacts.
- `tms`, `pjp`, `pjp-principle`: `rtk make doc`, `rtk make migration <name>`, and migration targets are authoritative.
- No controller/service/repository generator is detected for Fiber-only services. Manual files are permitted only after reading existing local pattern and recording the absence of a generator.
- Manual fallback after failed/unavailable generator must record command and reason in evidence.
