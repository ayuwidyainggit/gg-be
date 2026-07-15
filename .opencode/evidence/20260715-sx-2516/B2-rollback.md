# B2 — Rollback Review: `20260715_rollback_sales_update_temp.sql`

## Scope
Drop `import.sales_update_temp` table (reverse of B1).

## File reviewed
`sales/migration/20260715_rollback_sales_update_temp.sql`

## Content
```sql
DROP TABLE IF EXISTS import.sales_update_temp;
```

## Static validation
- `DROP TABLE IF EXISTS` — safe to run even if table does not exist.
- Transaction wrapped in `BEGIN`/`COMMIT` matching local migration convention.
- No cascade — explicit `DROP TABLE` without CASCADE to prevent accidental dependency removal.
- No secrets, tokens, or real data.

## Pair symmetry check
- B2 reverses B1: `CREATE TABLE IF NOT EXISTS` → `DROP TABLE IF EXISTS`. Symmetric.
- B2 does NOT touch any other table or column. Correct.

## No runtime SQL execution
This file was created but NOT applied against any database. B5 owns apply.

## Risks
- If data exists in `import.sales_update_temp`, it will be lost on rollback. Acceptable for staging table.
