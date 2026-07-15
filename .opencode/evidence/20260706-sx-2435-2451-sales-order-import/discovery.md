# Discovery Notes — SX-2435 / SX-2451 (sales order import follow-up)

## Repo evidence

- `.opencode/docs/ARCHITECTURE.md` confirms layering Controller → Service → Repository → DB and write-in-transaction rule.
- `.opencode/docs/QUALITY.md` mandates `rtk go test ./...` from the target service folder for `sales`.
- Prior plan `.opencode/plans/20260604-1024-sx-2154-order-type.md` already covered persistence of `order_type` and `original_qty_po*` for `sls.order` and `sls.order_detail`.
- Prior plan `.opencode/plans/20260608-1118-sx-2184-order-type-o-stock.md` already covered bypass of stock validation and inventory mutation for `order_type = "O"`.
- Prior plan `.opencode/plans/20260608-1347-sx-2131-2184-taking-order-audit.md` confirmed taking-order path is implemented and tests pass on `sales` at that time.
- `sales/service/order_type_helper.go` exposes the canonical helpers: `NormalizedOrderType`, `IsTakingOrder`, `ShouldValidateStockOnCreate`, `ShouldMutateInventoryOnCreate`, `BuildCreateOrderValidationBypassResponse`, `takingOrderQtySource`, `applyTakingOrderValidationSnapshot`, `applyTakingOrderDetailFields`.
- `sales/controller/order_controller.go` already wires `IsTakingOrder` / `ShouldValidateStockOnCreate` branch and normalizes empty `order_type` to nil.
- `sales/controller/order_controller_test.go` already has regression matrix for `O`, nil, empty, `C`, `SO`.
- `sales/service/order_type_helper_test.go` already has taking-order persistence tests including `original_qty_po*` and stock-mutation skip.

## Gaps discovered

- `SX-2435` and `SX-2451` are not referenced anywhere in the repo (no source comment, no plan, no doc).
- No existing sales order import endpoint or file-upload flow exists in the `sales` module.
- No `OrderImport` / `ImportOrders` symbol exists in controllers, services, repositories, entities, or models of `sales`.
- Existing Excel flow in `sales` is export-only (`sales/service/so_service.go` generates Excel via `excelize/v2` for `GET /sales/v1/download`); there is no reverse ingest path.

## Implication for planning

- The plan cannot be marked `PASS` or `PASS_FOR_SLICE` without a verified contract for SX-2435 / SX-2451.
- Status kept as `blocked` / `BLOCKED` with one required clarification input: exact ticket AC or payload sample.
- Reuse strategy is already pinned to the existing taking-order helpers and prior plans to avoid duplicate branches.

## Required clarification (minimum)

- For SX-2435: endpoint target (if any) + payload shape + acceptance criteria summary.
- For SX-2451: endpoint target (if any) + payload shape + acceptance criteria summary.
- Whether the work extends the create-order taking-order path, the update-enhance purchase-detail path, or a brand new import endpoint.
