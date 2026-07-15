# Discovery — SX-2314 Cancel Order stock reset

Task id: `20260624-0941-sx-2314-cancel-stock-reset`
Jira: `SX-2314` ([Defect][BE] Cancel Order not reset Warehouse Stock and On Cust Order)

## Endpoint
- `PATCH /sales/v1/orders/status` → `controller.OrderController.BulkUpdateStatus` (sales/controller/order_controller.go:895).
- Body `entity.BulkUpdateStatusOrder` → service `orderServiceImpl.BulkUpdateStatus` (sales/service/order_service.go:5131).

## Cancel branch (current)
`sales/service/order_service.go:5158-5234`:
1. `FindByNo` for order header.
2. Idempotent return `nil` if already `CANCELLED`.
3. `validateCancelTransition` allows only `NEED_REVIEW` (1) and `PROCESSED` (2) as source. Otherwise error.
4. `StockRepository.GetCancelStockBasis` reads basis from `inv.stock`/`sls.order_detail` SQL.
5. If `NEED_REVIEW` and `hasOutstandingStock=false` → `skipCancelStockWrite=true` → only update status.
6. Otherwise build commands → `StockRepository.CancelSalesStockUpdates` (sales/repository/stock_repository.go:398).

## Bug analysis
- `buildCancelStockMutations` (sales/repository/stock_repository.go:231) emits TWO stock rows for each basis row: rowA `tr_code='SO' tr_no=<SO> qty_in=qtyOutSO` and rowB `tr_code='SO' tr_no=<SO>-CO qty_out_order=qtyOutSO`. RowA is a duplicate of the original SO write; it should not be re-inserted on cancel. RowB's `tr_code` is `'SO'`, but docs require the reversal row to be `tr_code='CO'`. The cancel-aggregate subquery (`cancelAgg` line 312) also filters by `tr_code='SO'` against `tr_no='<SO>-CO'`; correct filter is `tr_code='CO'`. With both, the dedup math `qty_outstanding = qty_out_so - qty_out_order_cancel` works only on first cancel; subsequent cancels double-reverse.
- For `NEED_REVIEW` orders, the service short-circuits to "update status only" when no SO source ledger exists. Tester scenario with `2 0 2 → 12 Pieces` and a 0 `On Cust Order` baseline still shows a `qty_on_order=2` somewhere; without an `inv.stock tr_code='SO'` row to drive `QtyOutSmallest`, the reversal skip keeps `qty_on_order` stale. Need Review cancel must still release whatever `qty_on_order` the SO held when it was last active.
- `qty_outstanding` derivation does not fall back to detail qty (`final → sales → po`) when the SO source ledger is missing. Docs require the resolver priority to handle that.

## Existing helpers to reuse
- `entity.CancelStockBasis`, `entity.CancelStockWrite` already carry the resolver contract.
- `cancelStockBasisQuery` (stock_repository.go:286) is the right place to add the detail fallback.
- `repository.SalesStockUpdates` and `model.WarehouseStock` upsert path already does `qty += EXCLUDED.qty, qty_on_order += EXCLUDED.qty_on_order`; sign convention is correct.
- `mst.m_product.conv_unit1/2/3` referenced via `sls.order_detail.conv_unit2/conv_unit3` and product `mconv_unit2/3` (see `model.OrderDetailRead`).

## TDD baseline
- `repository/stock_repository_cancel_test.go` pins:
  - `TestBuildCancelStockMutations_SingleSKU` expects `len(stocks)==2`, rowA tr_code=`SO` tr_no=`<SO>`, rowB tr_code=`SO` tr_no=`<SO>-CO`.
  - `TestGetCancelStockBasisQuery_*` asserts SQL fragments including `SUM(c.qty_out_order) AS qty_out_order_cancel`, `qty1_final`/`qty1`/`qty_po1` fallback, `is_missing_source`, `is_ambiguous`.
- `service/order_service_test.go` pins `TestBulkUpdateStatus_Cancel_*` (5 cases) including Need Review skip-on-empty-basis.
- All baseline tests pass with the current code (verified: `go test ./repository/... ./service/... -run "TestBuildCancelStockMutations|TestGetCancelStockBasisQuery|TestBulkUpdateStatus_Cancel"`).

## Constraints
- Per AGENTS.md: each service is its own Go module, `go.mod`/Makefile/`.air.toml` per service. Use `rtk go test ./...` for the target service.
- Multi-tenant rules in `.opencode/docs/ARCHITECTURE.md` apply; `cust_id` is the auth context, no cross-tenant writes.
- Source files are off-limits for the planner. Implementation must be carried out by `@fixer`/`@backend` after handoff.

## Commands
- `rtk go test ./repository/... -run "TestBuildCancelStockMutations|TestGetCancelStockBasisQuery|TestCancelStockBasisFallback" -count=1 -v`
- `rtk go test ./service/... -run "TestBulkUpdateStatus_Cancel" -count=1 -v`
- `rtk go test ./... -count=1` (target service `sales`).
- DB inspection via `PGPASSWORD=postgres psql "host=localhost user=postgres dbname=ggn_scyllax sslmode=disable"` with the queries listed in the task prompt.
