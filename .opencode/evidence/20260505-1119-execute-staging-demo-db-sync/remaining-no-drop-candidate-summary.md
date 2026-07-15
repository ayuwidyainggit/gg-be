# Remaining No-Removal Candidate SQL Summary

File generated:

`.opencode/draft/20260505-1119-execute-staging-demo-db-sync/migration/20260505_remaining_schema_parity_no_drop_candidate.sql`

SHA256:

`27550cfabbc70581399683b3ba6d1b6e999fe557db691dfc429d25963ed0ce30`

Included:

- 7 remaining functions/procedural objects from staging.
- 3 remaining triggers from staging.
- 26 remaining staging-only constraints as `ADD CONSTRAINT` statements.
- 47 remaining staging-only indexes as `CREATE INDEX IF NOT EXISTS` / `CREATE UNIQUE INDEX IF NOT EXISTS` statements, with constraint-backed index statements commented as skipped when `ADD CONSTRAINT` would create them.
- Remaining `SET NOT NULL`: 0 statements detected at generation time.

Excluded by request:

- No table removal statements.
- No index removal statements.
- No constraint removal statements.

Still impossible to reach literal 100% schema identity without removal statements because demo has residual demo-only objects:

- 6 demo-only tables.
- 115 demo-only column keys on those demo-only tables.
- 27 demo-only constraint keys.
- 24 demo-only index keys.
- 1 demo-only sequence key.

Guardrail:

- Forbidden removal keyword scan: 0 matches after comment cleanup.
- Procedural/DML keywords are present as expected in function bodies, trigger events, and FK actions; the user requested this remaining candidate to improve parity except removal operations.

Execution status:

- Not executed.
