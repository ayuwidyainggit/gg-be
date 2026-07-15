# B3 — Alter Review: `20260715_alter_import_history.sql`

## Scope
Add `error_message` and `uploaded_by_name` columns to `import.import_history`.

## Confirmed source facts (A3)
- `import.import_history` exists with columns: `history_id` (PK), `file_name`, `uploaded_by`, `upload_date`, `successful_data`, `failed_data`, `total_data`, `status` (VARCHAR(50) NOT NULL DEFAULT 'COMPLETE'), `status_reupload`, `upload_type` (VARCHAR(50)), `cust_id` (VARCHAR).
- `error_message` — column absent.
- `uploaded_by_name` — column absent.
- `status` default is `'COMPLETE'` — must NOT be altered.

## File reviewed
`sales/migration/20260715_alter_import_history.sql`

## Content
```sql
ALTER TABLE IF EXISTS import.import_history
    ADD COLUMN IF NOT EXISTS error_message TEXT,
    ADD COLUMN IF NOT EXISTS uploaded_by_name VARCHAR(255);
```

## Static validation
- `ALTER TABLE IF EXISTS` + `ADD COLUMN IF NOT EXISTS` — fully idempotent.
- `error_message TEXT` — nullable, no default (only set on failure).
- `uploaded_by_name VARCHAR(255)` — nullable, no default (set on insert).
- `status` default `'COMPLETE'` is NOT touched. Service will explicitly insert `PROCESSING`.
- Transaction wrapped in `BEGIN`/`COMMIT`.
- No secrets, tokens, or real data.

## Pair symmetry check
- B3 adds only `error_message` and `uploaded_by_name`. Does not alter existing columns or defaults. Correct.

## No runtime SQL execution
This file was created but NOT applied against any database. B5 owns apply.

## Risks
- None. Adding nullable columns to an existing table is safe and non-disruptive.
