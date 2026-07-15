verdict: PASS

findings:
- BLOCKER none. The `cancelAgg` widening is correct for the stated legacy shape: the query is still parameterized on `cust_id` and exact `c.tr_no = ?` (`cancelTrNo`), while the added `LIKE '%-CO%'` is a constant SQL literal, not user input. No SQL injection or GORM binding concern is introduced.
- HIGH none. Dropping `od.item_type = 1` from the main basis query is the right fix to surface reward rows for the same order/product. It does not regress single-line orders because their qualifying rows still satisfy the unchanged qty predicate; removing a restrictive filter only adds the previously-missed reward rows.
- MEDIUM `activeDetailAgg` keeping `od.item_type = 1` is the right minimal mitigation from plan R1. It prevents order+reward-of-same-product from inflating `active_detail_count` and tripping `is_ambiguous=true` spuriously, while leaving true multi-order-line ambiguity detection for `item_type = 1` intact.
- MEDIUM The `GREATEST(..., 0)` clamp on `qty_out_smallest` does protect the write path from negative reversals if both a new `CO` row and a legacy `SO ...-CO` row are summed in `cancelAgg`. Note: the raw `qty_outstanding` alias itself is not clamped, but the mutation path uses the clamped smallest quantity, so the dangerous effect is contained.
- LOW The widened `cancelAgg` predicate has a redundant suffix check because `c.tr_no = ?` already binds the row to `<SO>-CO`; `LIKE '%-CO%'` only narrows nothing further for that exact value. This is harmless and keeps the legacy intent explicit.
- LOW Test naming is slightly misleading: `TestGetCancelStockBasisQuery_LegacySORowsExcludedFromCancelAgg` actually asserts inclusion of legacy SO reversal rows in `cancelAgg`, not exclusion. This is a readability issue only, not a logic problem.

recommendations:
- Optional: rename `TestGetCancelStockBasisQuery_LegacySORowsExcludedFromCancelAgg` to an inclusion-oriented name on a later cleanup pass for clarity.
- Optional: if future callers ever start consuming raw `qty_outstanding` for decisions, clamp or validate that field too; today the bounded patch is acceptable because `qty_out_smallest` is already guarded.

confirmation:
- The patch is minimal, bounded, and aligned with plan v4: one `cancelAgg` predicate widening, one main-query filter removal, and targeted dry-run tests plus DB simulation evidence.
- Ready for `@quality-gate`.

evidence:
- Repo-local evidence: `sales/repository/stock_repository.go:287-297`, `:299-319`, `:345-396`; `sales/repository/stock_repository_cancel_test.go:154-255`.
- Supporting evidence: `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/tests-fail.txt`, `tests-pass.txt`, `simulate-cancel-v2.txt`, `legacy-audit.sql.txt`, `legacy-audit-post.txt`.
