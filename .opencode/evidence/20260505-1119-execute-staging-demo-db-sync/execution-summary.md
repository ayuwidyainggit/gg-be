# Execution Summary — DB Schema Sync Staging to Demo

## Applied target

- Target applied: remote demo database, after explicit user confirmation.
- Backup before apply: `e609250a2611412361752b29cc9739975183ef57feaea29b2596c6a1e0199a76  demo_before_apply_20260505_120711.dump`.
- Migration file: `.opencode/draft/20260505-1119-execute-staging-demo-db-sync/migration/20260505_sync_demo_schema_to_staging_strict_no_dml.sql`.

## Guardrails

- DML/data movement guardrail: 0 matches after final migration generation.
- Destructive DDL guardrail: 0 matches after final migration generation.
- Explicitly excluded: data copy/mutation, `DELETE`, destructive DDL, function/procedure/trigger bodies containing DML keywords, FK `ON UPDATE`/`ON DELETE` actions, NOT NULL enforcement for added columns requiring data backfill.

## Apply result

- Final psql apply result: success (`COMMIT`).
- Earlier failed attempts were inside transaction and rolled back by `ON_ERROR_STOP`/transaction abort; final successful run used adjusted migration.

## Initial catalog diff summary

# Catalog Diff Summary

| object | staging count | demo count | staging-only rows | demo-only rows |
|---|---:|---:|---:|---:|
| schemas | 20 | 19 | 1 | 0 |
| tables | 481 | 455 | 32 | 6 |
| columns | 9217 | 8779 | 1865 | 1427 |
| constraints | 429 | 332 | 127 | 30 |
| indexes | 428 | 341 | 111 | 24 |
| functions | 105 | 62 | 43 | 0 |
| triggers | 3 | 0 | 3 | 0 |
| sequences | 259 | 242 | 23 | 6 |
| extensions | 5 | 4 | 1 | 0 |


## Post-apply catalog diff summary

# Catalog Diff Summary After Migration

| object | staging count | demo-after count | staging-only rows | demo-after-only rows |
|---|---:|---:|---:|---:|
| schemas | 20 | 20 | 0 | 0 |
| tables | 481 | 487 | 0 | 6 |
| columns | 9217 | 9332 | 1380 | 1495 |
| constraints | 429 | 430 | 50 | 51 |
| indexes | 428 | 405 | 47 | 24 |
| functions | 105 | 98 | 7 | 0 |
| triggers | 3 | 0 | 3 | 0 |
| sequences | 259 | 260 | 15 | 16 |
| extensions | 5 | 5 | 0 | 0 |


## Residual drift by design

- Staging-only tables after apply: 0.
- Demo-only tables after apply: 6, retained because destructive DDL is prohibited.
- Staging-only functions after apply: 7, not applied because function bodies contain DML keywords and would fail the no-DML guardrail.
- Staging-only triggers after apply: 3, not applied because trigger definitions include DML event keywords and depend on excluded functions.
- Remaining column definition differences include intentional softening, especially nullable columns where staging has `NOT NULL`; enforcing them would require data backfill/mutation or fail on existing demo data.
- Remaining FK action differences include omitted `ON UPDATE`/`ON DELETE` actions to satisfy strict no-DML/destructive keyword guardrail.
