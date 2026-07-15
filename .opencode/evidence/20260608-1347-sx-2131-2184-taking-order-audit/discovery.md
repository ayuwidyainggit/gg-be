# Discovery SX-2131 / SX-2184 — Taking Order audit

Task ID: `20260608-1347-sx-2131-2184-taking-order-audit`
Tanggal: 2026-06-08 Asia/Jakarta
Mode: Maintenance Stability Mode

## Source strategy

- Dipakai: bukti repo lokal, dokumen harness lokal, plan/evidence SX-2184 sebelumnya, grep kode, pembacaan file target, targeted tests, full sales tests.
- Tidak dipakai langsung: Jira/Google Docs karena prompt sudah memberi detail acceptance criteria dan sebagian referensi kemungkinan membutuhkan akses. Klaim eksternal tidak ditambah di luar prompt.
- Tidak dipakai: GitLab MR API karena repo lokal sudah berada di branch `dev` dan perubahan MR terlihat ada di working tree/commit lokal.
- Tidak dipakai: browser/API smoke baru, karena evidence smoke sudah ada di `.opencode/evidence/20260608-1118-sx-2184-order-type-o-stock/implementation.md` dan test suite lokal sudah cukup untuk audit implementasi. Jika release gate butuh bukti terbaru, executor/quality-gate dapat mengulang smoke.

## Files inspected

- `AGENTS.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/AGENT_ROUTING.md`
- `.opencode/docs/QUALITY.md`
- `.opencode/plans/20260608-1118-sx-2184-order-type-o-stock.md`
- `.opencode/evidence/20260608-1118-sx-2184-order-type-o-stock/implementation.md`
- `sales/controller/order_controller.go`
- `sales/controller/order_controller_test.go`
- `sales/service/order_service.go`
- `sales/service/validate_order_service.go`
- `sales/service/order_type_helper.go`
- `sales/service/order_type_helper_test.go`
- `sales/entity/order.go`
- `sales/entity/order_detail.go`
- `sales/model/order.go`
- `sales/model/order_detail.go`
- `sales/migration/sls.order/add_order_type_and_original_qty_po_fields.sql`
- `sales/migration/sls.order/rollback_add_order_type_and_original_qty_po_fields.sql`
- `sales/go.mod`

## Grep / project patterns found

Search pattern:

```text
order_type|IsTakingOrder|TakingOrder|original_qty_po|validate_stok|inv\.stock|warehouse_stock|opr_type
```

Key findings:

- `sales/entity/order.go` has `OrderType *string json:"order_type" validate:"omitempty,oneof=O C SO"`.
- `sales/model/order.go` has `OrderType *string` and nullable `ValidateStokMessage *string`.
- `sales/entity/order_detail.go` and `sales/model/order_detail.go` have `OriginalQtyPo1/2/3`.
- `sales/service/order_type_helper.go` defines `IsTakingOrder`, `ShouldValidateStockOnCreate`, `ShouldMutateInventoryOnCreate`, `BuildCreateOrderValidationBypassResponse`, and `applyTakingOrderDetailFields`.
- `sales/controller/order_controller.go` normalizes empty `order_type` to nil and calls `ValidateOrderWithoutStock` for `O`, while non-`O` calls `ValidateOrder`.
- `sales/service/validate_order_service.go` supports `ValidateOrderWithoutStock`, skipping only warehouse stock validation while still running AR, credit-limit, overdue, and outstanding logic.
- `sales/service/order_service.go` applies taking-order validation snapshot, sets fallback `opr_type = O`, maps PO/original qty, nils sales qty tiers, zeroes `Qty`/`QtyFinal`, and guards inventory mutation through `ShouldMutateInventoryOnCreate`.
- Migration/rollback files already exist for `sls.order.order_type` and `sls.order_detail.original_qty_po1/2/3`.

## Validation commands checked

From `/Users/ujang/Projects/Geekgarden/scylla-be/sales`:

```bash
rtk go test ./service -run 'Test.*SX2184|Test.*OrderType|Test.*TakingOrder'
rtk go test ./controller -run 'Test.*SX2184|Test.*Create.*OrderType'
```

Result:

```text
Go test: 5 passed in 1 packages
Go test: 5 passed in 1 packages
```

Full suite:

```bash
rtk go test ./...
```

Result:

```text
Go test: 219 passed in 22 packages
```

Git status from `sales`:

- `git diff --stat && git diff --name-only` returned no output, so no new source changes were made by this planner.
- Current branch: `dev`.
- Current short commit: `53d68b3`.

## Reuse candidates

- Existing implementation from prior SX-2184 plan should be reused; no new implementation needed from planner.
- Existing test files cover controller routing, service inventory mutation, nullable stock message, and taking-order qty persistence semantics.
- Existing migration files are additive and idempotent via `IF NOT EXISTS`.

## Constraints

- Planner must not edit source, tests, package, lockfile, migration, or env files outside `.opencode/`.
- Repo requires `rtk`-prefixed shell workflows.
- Service layering must remain Controller → Service → Repository → DB.
- Write operations must stay in service transactions.
- No secrets/tokens/env should be copied.

## Risks / gaps

- I did not re-run live API/DB smoke in this audit. Prior evidence claims successful local smoke for `order_type = O`; quality gate may decide whether to refresh it.
- `SO` live smoke was not previously run because existing stock formula can be misleading; regression is covered by unit/controller tests proving non-`O` paths still call stock validation/mutation.
- The repo root is not a git repo; `sales` is the service git repo for branch/hash/status.
- If downstream clients strictly expect `validate_stok_message` to always be a string, nullable read shape should be communicated, but it matches requested SQL NULL behavior.
