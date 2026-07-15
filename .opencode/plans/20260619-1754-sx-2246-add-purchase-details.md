# Plan — SX-2246 Add Purchase Details (Need Reviews)

Task ID: `20260619-1754-sx-2246-add-purchase-details`
Parent: `SX-2242`
PIC: `dev.be` (Yogie)
Priority: Medium
Endpoint: `PATCH /sales/v1/orders/enhance/{order_no}`

## Goal

Make `add_purchase_details` from the Purchase Order tab correctly insert a new `sls.order_detail` row per item with UOM-aware stock capping, `original_qty_po*` raw preservation, payload-driven `unit_id*` and `is_product_promotion_po`, while preserving the existing `purchase_order` update flow and atomic transaction behavior for `Need Reviews` orders.

## Non-goals

- No change to the `purchase_order` (existing detail update) semantics except reusing safer guards.
- No change to the `Need Reviews` approval transition; that flow is owned by other tasks.
- No new external API, no schema migration.
- No destructive de-dupe; if duplicate-add policy is unclear, leave rows intact and emit a warning log instead of deleting.

## Scope

- Files: `sales/entity/edit_order_enhance.go`, `sales/service/order_service.go`, `sales/service/order_service_test.go`.
- Service layer only; controller and route stay as-is.
- Transaction boundary stays inside `UpdateEnhance`.

## Requirements

- `add_purchase_details` is already accepted at the controller and aliased to `add_purchase_order` by `normalizeEnhancePromoFlags`; no controller change required.
- For each item in `add_purchase_details`, insert a row into `sls.order_detail` with:
  - `pro_id` from payload.
  - `unit_id1/2/3` from payload, fallback to product master when empty.
  - `sell_price_system1/2/3`, `sell_price_po1/2/3` from payload.
  - `original_qty_po1/2/3` raw from payload.
  - `qty_po1/2/3`, `qty1/2/3`, `qty1_final/2_final/3_final` from UOM-aware stock cap.
  - `qty1_stok/2_stok/3_stok` from current warehouse stock at insert time.
  - `is_product_promotion_po` from payload.
  - `conv_unit2/3`, `unit_id1/2/3` (master fallback), `vat`, `item_type=1`, default `qty_po = qty = qty_final` consistent with existing add path.
- UOM-aware cap rule:
  - Pull current stock from `StockRepository.GetCurrentStock(custId, whId, proId)`.
  - Convert total small-unit stock to qty1/2/3 via `canonicalAPIStockBreakdown`.
  - Compare per-level available vs requested; the persisted value is the available stock, never negative, never above requested. `original_qty_po*` is independent of the cap.
  - If `ro.WhId` is missing, skip cap (return error) because warehouse is required for the add path; same as current add flow.
- Existing `purchase_order` updates must not change behavior.
- Transactional rollback: keep all inserts/updates inside the existing `WithinTransaction` block.

## Acceptance Criteria

1. `PATCH /sales/v1/orders/enhance/{order_no}` accepts `add_purchase_details` payload.
2. Each item inserts a new `sls.order_detail` row.
3. `pro_id` is taken from the payload.
4. `original_qty_po1/2/3` equal raw `qty_po1/2/3` from the payload.
5. `qty_po1/2/3`, `qty1/2/3`, `qty1_final/2_final/3_final` are UOM-aware stock-capped values, never above requested.
6. `is_product_promotion_po` is taken from the payload.
7. `unit_id1/2/3` from payload override product master only when non-empty.
8. `qty1_stok/2_stok/3_stok` reflect current warehouse stock at insert time.
9. Existing `purchase_order` update still works.
10. Rollback occurs if any step errors.
11. Tests cover:
    - no `add_purchase_details` (alias smoke test already exists; reuse).
    - single add product with full payload field assertions.
    - multiple add products.
    - stock 0 → capped to 0.
    - requested qty greater than available stock → capped to available.
    - requested qty considered available by UOM conversion → kept as requested.

## Existing Patterns / Reuse

- `OrderRepository.FindProductByID` for `conv_unit2/3`, `unit_id1/2/3`, `vat`.
- `OrderRepository.StoreDetail` for insert.
- `StockRepository.GetCurrentStock` + `canonicalAPIStockBreakdown` for stock snapshot.
- `calculateNormalizedQty` for total qty.
- `applyTakingOrderDetailFields` as reference for `original_qty_po*` semantics.
- `mockOrderRepository` / `mockStockRepository` / `mockDbtransaction` in `order_service_test.go` for unit tests.

## Constraints

- Controller → Service → Repository → DB layering stays strict.
- Writes use `txCtx`; never `ctx` for repository writes inside the transaction.
- `cust_id` filter must be applied to all read and write paths.
- Follow repo-local `rtk` shell prefix as per project `AGENTS.md`.

## Risks

- Wrong UOM cap interpretation can silently drop qty.
- Duplicate insert on FE retry could cause inflated totals; not in scope to dedupe.
- `ro.WhId` is dereferenced in existing add helper; if nil, must return explicit error instead of panic.
- Stock lookup failure (transient DB) must not break insert; fallback to storing requested qty and log warning.

## Decisions / Assumptions

- Assumption: `add_purchase_details` payload from FE is always new SKUs (no in-place edit). Duplicate prevention deferred until product owner confirms.
- Decision: cap rule is per-level `min(requested, available)` after `canonicalAPIStockBreakdown` conversion. The docs ambiguity (stock 10, qty 12 → keep 12) is resolved by treating the cap as upper-bound, not lower-bound; no `Math.min(qty, stock)` raw without UOM. This satisfies all three rule lines: cap to available when stock < qty, cap to qty when stock ≥ qty, cap to 0 when stock = 0.
- Decision: use payload `unit_id*` only when non-empty; otherwise inherit from product master (matches current `createOrderDetailFromPurchaseOrder`).
- Decision: `is_product_promotion_po` defaults to `false` when nil.

## Execution Source of Truth

Priority order for implementation:
1. Latest explicit user instruction in this task.
2. Safety, security, permission rules.
3. Non-negotiable Implementation Invariants.
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria and Done Criteria.
6. Implementation Steps.

If any conflict appears, executor must follow the higher source and log the conflict in verification evidence.

## Non-negotiable Implementation Invariants

- Do not modify `purchase_order` (existing detail update) behavior beyond reusing the same stock snapshot helper.
- Do not bypass `WithinTransaction`; all writes go through `txCtx`.
- Do not delete or replace existing rows for `add_purchase_details` dedupe.
- Do not commit secrets, tokens, or `.env` files. Do not log the bearer token.
- Do not collapse Controller → Service → Repository boundaries; repository stays data access only.
- Stock cap must be UOM-aware using `canonicalAPIStockBreakdown`, never raw `min(qty, stock)`.
- `original_qty_po*` must be raw payload, not stock-capped.

## Do Not / Reject If

- Insert uses raw `Math.Min(qty, stock)` without UOM conversion → reject.
- Insert drops `original_qty_po*` or sets it equal to capped qty → reject.
- Insert deletes existing rows for SKU dedupe without business rule confirmation → reject.
- Insert runs outside `txCtx` → reject.
- Insert adds `add_purchase_details` to allowlist but skips `add_purchase_order` → reject.
- Plan finalizes without addressing duplicate-resend behavior at least by warning log → reject.

## Diff Boundary

Allowed file groups:
- `sales/entity/edit_order_enhance.go` (add fields to `AddPurchaseOrderDetail`).
- `sales/service/order_service.go` (extend `createOrderDetailFromPurchaseOrder`).
- `sales/service/order_service_test.go` (extend tests).

Off-limits:
- `sales/controller/order_controller.go` (no signature change; payload already accepted).
- `sales/repository/*` (no schema/query change).
- `sales/model/*` (no schema change; existing `OrderDetail` already has all target fields).
- Any other service or `go.mod` / lockfile.

Evidence paths:
- `.opencode/evidence/20260619-1754-sx-2246-add-purchase-details/discovery.md`
- `.opencode/evidence/20260619-1754-sx-2246-add-purchase-details/verification.md` (created by executor)

## TDD / Test Plan

- TDD required for this change. Service-layer logic.
- Existing pattern: `TestCreateOrderDetailFromPurchaseOrder_InheritsUOMFromProductMaster` and `TestUpdateEnhance_NormalizePurchaseDetailsAlias`.
- First failing test: `TestCreateOrderDetailFromPurchaseOrder_PersistsFullAddPurchaseDetailsPayload` (asserts `original_qty_po*`, `is_product_promotion_po`, `unit_id*`, `qty1_stok/2_stok/3_stok`).
- Green step: extend entity + service.
- Refactor step: extract `applyStockCapToAddDetail` helper if duplication emerges.
- Edge case tests:
  - `TestCreateOrderDetailFromPurchaseOrder_CapsQtyToAvailableStock`.
  - `TestCreateOrderDetailFromPurchaseOrder_StoresZeroQtyWhenStockEmpty`.
  - `TestCreateOrderDetailFromPurchaseOrder_KeepsRequestedQtyWhenStockSufficient`.
  - `TestCreateOrderDetailFromPurchaseOrder_FallsBackUnitIdsFromProductMasterWhenPayloadEmpty`.
  - `TestUpdateEnhance_AddPurchaseDetails_InsertsMultipleRowsAtomically` (use mock tx with rollback failure).
- Commands: `cd sales && rtk go test ./service -run 'TestCreateOrderDetailFromPurchaseOrder|TestUpdateEnhance.*AddPurchase' -count=1`.

## Implementation Steps

1. Extend `AddPurchaseOrderDetail` with `UnitId1/2/3`, `IsProductPromotionPo *bool` (FE flag already sent as `is_product_promotion_po`).
2. In `createOrderDetailFromPurchaseOrder`:
   - Resolve `unit_id1/2/3` from payload with product master fallback.
   - Populate `original_qty_po1/2/3` from payload raw qty.
   - Resolve stock via `GetCurrentStock` if `whId` available; if stock lookup errors, log warn and keep requested qty.
   - Convert stock to qty breakdown via `canonicalAPIStockBreakdown`.
   - Cap per-level qty: `qty_po1 = min(qty_po1_requested, available_qty1)` etc. Same value written to `qty1`, `qty1_final`, etc.
   - Set `qty1_stok/2_stok/3_stok` from the breakdown.
   - Set `is_product_promotion_po` from payload (default false).
   - Guard against `ro.WhId == nil` → return error.
3. No change to controller, repository, or transaction wrapper.
4. Add unit tests covering acceptance criteria.

## Expected Files to Change

- `sales/entity/edit_order_enhance.go` — add `UnitId1/2/3`, `IsProductPromotionPo` to `AddPurchaseOrderDetail`.
- `sales/service/order_service.go` — extend `createOrderDetailFromPurchaseOrder` only.
- `sales/service/order_service_test.go` — add new tests.

## Agent / Tool Routing

- Implementer: `@fixer` with `opencode-backend` skill.
- Reviewer: `@quality-gate` for material risk + multi-row insert.
- No `@designer` / `@architect` needed (no UI/UX or platform architecture change).
- No `@oracle` review needed (small bounded change).

## Executor Handoff Prompt

Copyable prompt for `@fixer`:

```
Implement SX-2246 add_purchase_details insert fix per plan:
.opencode/plans/20260619-1754-sx-2246-add-purchase-details.md

Scope:
- sales/entity/edit_order_enhance.go: extend AddPurchaseOrderDetail with UnitId1/2/3 + IsProductPromotionPo
- sales/service/order_service.go: extend createOrderDetailFromPurchaseOrder only
- sales/service/order_service_test.go: add tests

Must preserve:
- existing purchase_order update flow
- transaction wrapper, cust_id filter, txCtx usage
- service-layer contract (no controller/repo change)

Must enforce:
- original_qty_po* = raw payload
- qty_po* / qty* / qty*_final = UOM-aware stock cap via canonicalAPIStockBreakdown
- unit_id* from payload, fallback to product master
- qty1_stok/2_stok/3_stok set at insert
- is_product_promotion_po from payload

Reject if:
- raw min(qty, stock) without UOM
- writes outside txCtx
- destructive dedupe
- any other service file edited

Validate:
- cd sales && rtk go mod download
- cd sales && rtk go mod tidy
- cd sales && rtk go test ./service -run 'TestCreateOrderDetailFromPurchaseOrder|TestUpdateEnhance.*AddPurchase' -count=1
- cd sales && rtk go test ./... -count=1

Return:
- diff summary
- test run output
- any deviation with rationale
```

## Execution-ready Worklist / Handoff Contract

| id | action | deps | owner | validation | exit | blocking | requires_user_decision | must_preserve | do_not_touch | evidence_update | exit_verification | start_with |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T1 | Add `UnitId1/2/3` and `IsProductPromotionPo *bool` to `AddPurchaseOrderDetail` in `sales/entity/edit_order_enhance.go`. | none | `@fixer` | `cd sales && rtk go build ./entity/...` | entity compiles | ready | no | entity JSON tags match FE payload keys | controller, repo, model | append verification.md entry | go build passes | yes |
| T2 | Add failing test `TestCreateOrderDetailFromPurchaseOrder_PersistsFullAddPurchaseDetailsPayload` asserting `original_qty_po*`, `is_product_promotion_po`, `unit_id*` (payload override), `qty1_stok/2_stok/3_stok`. | T1 | `@fixer` | `cd sales && rtk go test ./service -run TestCreateOrderDetailFromPurchaseOrder_PersistsFullAddPurchaseDetailsPayload -count=1` | test fails (Red) | ready | no | mock repo pattern from existing tests | other tests | append verification.md | test fails as expected | no |
| T3 | Extend `createOrderDetailFromPurchaseOrder` to populate the new fields and call stock cap helper. | T2 | `@fixer` | re-run T2 | test passes (Green) | ready | no | uses `GetCurrentStock` + `canonicalAPIStockBreakdown` + `calculateNormalizedQty` | other service code | append verification.md | T2 passes | no |
| T4 | Add edge tests: stock 0 caps to 0; available < requested caps to available; available ≥ requested keeps requested; unit_id fallback to product master. | T3 | `@fixer` | `cd sales && rtk go test ./service -run TestCreateOrderDetailFromPurchaseOrder -count=1` | all pass | ready | no | each test uses independent mock state | shared mock state | append verification.md | all pass | no |
| T5 | Add `TestUpdateEnhance_AddPurchaseDetails_InsertsMultipleRowsAtomically` exercising both success and rollback on second insert failure. | T3 | `@fixer` | `cd sales && rtk go test ./service -run TestUpdateEnhance_AddPurchaseDetails -count=1` | pass | ready | no | use `mockDbtransaction` | real DB | append verification.md | pass | no |
| T6 | Run full sales test suite to confirm no regression in `purchase_order` update path. | T5 | `@fixer` | `cd sales && rtk go test ./... -count=1` | pass | ready | no | existing tests untouched in semantics | other files | append verification.md | pass | no |
| T7 | Hand off to `@quality-gate` for final review with diff summary and test output. | T6 | `@quality-gate` | review pass | review pass | ready | no | evidence paths intact | none | final verification.md | review pass | no |

`start_with`: T1.

## Validation Commands

```bash
cd /Users/ujang/Projects/Geekgarden/scylla-be/sales
rtk go mod download
rtk go mod tidy
rtk go test ./service -run 'TestCreateOrderDetailFromPurchaseOrder|TestUpdateEnhance.*AddPurchase' -count=1
rtk go test ./... -count=1
```

## Evidence Requirements

- `.opencode/evidence/20260619-1754-sx-2246-add-purchase-details/discovery.md` (already written).
- `.opencode/evidence/20260619-1754-sx-2246-add-purchase-details/verification.md` (executor appends per task with test output, key diff lines, and any deviation).
- Final review note from `@quality-gate`.

## Done Criteria

- All TDD tests pass.
- No regression in existing `purchase_order` update tests.
- Transaction rollback path covered.
- Plan file and evidence file up to date.
- No secrets committed; no bearer token in logs.

## Final Planning Summary

- Artifacts consulted: `AGENTS.md`, `.opencode/docs/ARCHITECTURE.md`, sales source.
- Artifacts created: discovery under evidence, this plan.
- Key decisions: UOM-aware stock cap via existing helpers; preserve existing `purchase_order` semantics; payload `unit_id*` overrides product master only when non-empty; `original_qty_po*` stays raw.
- Assumptions: `add_purchase_details` items are new SKUs per call; dedupe policy deferred.
- Open questions: none blocking. Duplicate-resend policy is a follow-up, not a blocker.
- Readiness: `ready-for-implementation`.
- Cleanup: `.opencode/draft/20260619-1754-sx-2246-add-purchase-details/` may stay empty; remove at finalization if unused.
