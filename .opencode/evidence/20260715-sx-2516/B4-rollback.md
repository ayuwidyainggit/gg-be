# B4 — Rollback Review: `20260715_rollback_alter_import_history.sql`

## Scope
Remove `error_message` and `uploaded_by_name` columns from `import.import_history` (reverse of B3).

## File reviewed
`sales/migration/20260715_rollback_alter_import_history.sql`

## Content
```sql
ALTER TABLE IF EXISTS import.import_history
    DROP COLUMN IF EXISTS error_message,
    DROP COLUMN IF EXISTS uploaded_by_name;
```

## Static validation
- `ALTER TABLE IF EXISTS` + `DROP COLUMN IF EXISTS` — safe to run even if columns do not exist.
- Drops only the two columns added by B3. Does NOT touch `status`, `upload_type`, `cust_id`, or any other existing column.
- Transaction wrapped in `BEGIN`/`COMMIT`.
- No secrets, tokens, or real data.

## Pair symmetry check
- B4 reverses B3: `ADD COLUMN IF NOT EXISTS error_message` → `DROP COLUMN IF EXISTS error_message`. Symmetric.
- B4 reverses B3: `ADD COLUMN IF NOT EXISTS uploaded_by_name` → `DROP COLUMN IF EXISTS uploaded_by_name`. Symmetric.
- B4 does NOT drop any column not added by B3. Correct.

## No runtime SQL execution
This file was created but NOT applied against any database. B5 owns apply.

## Risks
- If data exists in `error_message` or `uploaded_by_name` columns, it will be lost on rollback. Acceptable for migration rollback.
