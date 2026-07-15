# Remaining Candidate Apply Summary

Applied file:

`.opencode/draft/20260505-1119-execute-staging-demo-db-sync/migration/20260505_remaining_schema_parity_safe_no_drop_executable.sql`

Result: success (`COMMIT`).

Backup before apply:

`d6983eb558c9dd7799a90e0d450894b473b5cbe28be91ffdb6b1237b099e4e72  demo_before_remaining_candidate_20260505_122729.dump`

Final catalog summary:

# Catalog Diff Summary Final

| object | staging count | demo-final count | staging-only rows | demo-final-only rows |
|---|---:|---:|---:|---:|
| schemas | 20 | 20 | 0 | 0 |
| tables | 481 | 487 | 0 | 6 |
| columns | 9217 | 9332 | 1380 | 1495 |
| constraints | 429 | 430 | 50 | 51 |
| indexes | 428 | 439 | 13 | 24 |
| functions | 105 | 105 | 0 | 0 |
| triggers | 3 | 3 | 0 | 0 |
| sequences | 259 | 260 | 15 | 16 |
| extensions | 5 | 5 | 0 | 0 |


Notes:

- The broader candidate with all constraints failed because existing demo data violated FK `fk_replenishment_order_delivery_type` (`delivery_type = Full` not present in `inv.delivery_type`). It was not applied because transaction aborted.
- The successful executable candidate applied functions, triggers, and safe non-unique indexes only.
- Remaining staging-only indexes are constraint-backed/unique indexes that require adding constraints or resolving data/name risks.
- Remaining exact parity gaps still include demo-only tables/columns/constraints/indexes/sequence that cannot be removed under the no-removal constraint.
