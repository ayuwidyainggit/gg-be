# SX-1472 Patch Plan

## Objective

Patch backend behavior for mobile-source order processing so that no-change process requests only update header status and do not create stock mutations.

## Confirmed diagnosis

- Flow [`orderServiceImpl.Update()`](sales/service/order_service.go:3568) currently behaves like full detail sync for mobile process requests.
- This flow can still:
  - sync detail rows
  - delete detail rows not present in payload
  - build stock update entities
  - call [`StockRepositoryImpl.SalesStockUpdates()`](sales/repository/stock_repository.go:133)
- Remote DB validation confirms processed mobile orders already have stock ledger rows with `tr_no = ro_no`.
- Root cause is most likely a combination of:
  - missing no-change guard in service layer
  - missing delta-zero safety net in stock layer
  - compare basis for mobile not being explicitly normalized around PO fields

## Scope of patch

### In scope

- Add service-level no-change guard for mobile process flow
- Add stock-layer delta-zero guard
- Keep existing transaction boundary intact
- Keep repository/business-layer separation intact
- Preserve current behavior for actual edits:
  - qty edit
  - add product
  - delete product

### Out of scope

- Reworking all order update flows
- Rewriting enhance flow unless needed for parity hardening
- Changing FE payload contract
- Refactoring unrelated promotion or discount logic

## Target files and required patch points

### 1. Primary patch in [`order_service.go`](sales/service/order_service.go)

#### Function to patch
- [`orderServiceImpl.Update()`](sales/service/order_service.go:3568)

#### Required changes
- Add early mobile-specific decision branch after fetching existing order and before detail sync/delete logic.
- Introduce helper calls inside service layer to:
  - load current order details
  - normalize incoming mobile process details
  - diff existing vs incoming using PO semantics
  - detect `no meaningful change`
- If no-change is detected:
  - update header only
  - set `data_status = 2`
  - skip detail sync loop
  - skip `FindDetailByNotInDetailIDs`
  - skip `DeleteDetailNotInIDs`
  - skip stock update entity creation
  - skip call to [`SalesStockUpdates()`](sales/repository/stock_repository.go:133)
  - commit transaction and return success

#### New helper candidates inside service layer
- [`orderServiceImpl.normalizeMobileProcessDetails()`](sales/service/order_service.go)
- [`orderServiceImpl.diffMobileProcessDetails()`](sales/service/order_service.go)
- [`orderServiceImpl.hasMeaningfulMobileProcessChanges()`](sales/service/order_service.go)

These helpers should stay in service layer because they contain business rules.

### 2. Safety net patch in [`stock_repository.go`](sales/repository/stock_repository.go)

#### Function to patch
- [`StockRepositoryImpl.SalesStockUpdates()`](sales/repository/stock_repository.go:133)

#### Required changes
- Before generating stock ledger rows / warehouse updates, compute effective delta.
- If effective delta in smallest unit is zero:
  - skip current item entirely
- If after filtering all items there are no effective deltas:
  - return nil without inserting [`inv.stock`](docs/Sales%20Order%20Enhancement_BE.md)
  - return nil without updating [`inv.warehouse_stock`](docs/Sales%20Order%20Enhancement_BE.md)

This is mandatory as a defense-in-depth guard even if service layer is fixed.

### 3. Repository usage constraints to preserve

Do not move business diff logic into repository.

Repository functions that should remain unchanged unless needed for compatibility:
- [`RepositoryOrderImpl.FindDetailByNotInDetailIDs()`](sales/repository/order_repository.go:205)
- [`RepositoryOrderImpl.DeleteDetailNotInIDs()`](sales/repository/order_repository.go:445)
- [`RepositoryOrderImpl.Update()`](sales/repository/order_repository.go:435)
- [`RepositoryOrderImpl.UpdateDetail()`](sales/repository/order_repository.go:468)
- [`RepositoryOrderImpl.SyncFinalOrderFields()`](sales/repository/order_repository.go:1008)

The implementation goal is to avoid calling delete/sync functions in no-change process mode, not to relocate their logic.

## Business compare rules

### Compare basis for source mobile

Use these fields as the authoritative compare set for mobile process flow:
- `pro_id`
- `unit_id1`
- `unit_id2`
- `unit_id3`
- `qty_po1`
- `qty_po2`
- `qty_po3`

### Normalization rules

- Normalize `null` to `0` for qty fields.
- Normalize missing PO qty from legacy payload using fallback:
  - `qty_po1 ?? qty1 ?? 0`
  - `qty_po2 ?? qty2 ?? 0`
  - `qty_po3 ?? qty3 ?? 0`
- Use smallest-unit conversion for delta calculation.
- Sort or map by stable identity:
  - preferred: `order_detail_id`
  - fallback: composite of `pro_id + unit_id1 + unit_id2 + unit_id3`

### Meaningful change definition

Treat as meaningful change only when one of the following happens:
- PO qty changes
- product identity changes
- product is added
- product is deleted

Do not treat as meaningful change when:
- only `data_status` changes
- FE resends old detail payload unchanged
- qty values are equivalent after normalization

## Pseudocode for implementation

```pseudo
function Update(roNo, request, validationData):
    begin transaction

    order = getOrder(roNo)

    if order.data_source == MOBILE:
        existingDetails = getOrderDetails(roNo)
        incomingDetails = normalizeMobileProcessDetails(request.details.normal)

        if incomingDetails is empty:
            updateOrderHeaderOnly(roNo, data_status = 2)
            commit
            return

        changes = diffMobileProcessDetails(existingDetails, incomingDetails)

        if changes.isNoChange():
            updateOrderHeaderOnly(roNo, data_status = 2)
            commit
            return

    continue existing legacy update flow for real changes
```

```pseudo
function diffMobileProcessDetails(existing, incoming):
    normalize both datasets
    classify into:
      unchanged
      updated
      added
      deleted
    return diff result
```

```pseudo
function SalesStockUpdates(updates):
    filtered = []

    for update in updates:
        delta = smallest(update.qtyOrder) - smallest(update.qtyOrderBefore)
        if delta == 0:
            continue
        filtered.append(update)

    if filtered is empty:
        return nil

    continue existing stock insert/update logic using filtered
```

## Detailed execution plan for code mode

- [ ] Patch [`orderServiceImpl.Update()`](sales/service/order_service.go:3568) to introduce mobile no-change branch before detail sync/delete
- [ ] Add service-layer helper to normalize incoming mobile detail payload into PO-based compare structure
- [ ] Add service-layer helper to diff existing detail vs incoming detail and classify unchanged/updated/added/deleted
- [ ] Ensure no-change branch only updates header via [`RepositoryOrderImpl.Update()`](sales/repository/order_repository.go:435)
- [ ] Ensure no-change branch skips calls to [`FindDetailByNotInDetailIDs()`](sales/repository/order_repository.go:205) and [`DeleteDetailNotInIDs()`](sales/repository/order_repository.go:445)
- [ ] Ensure no-change branch skips building `salesOrderStockUpdateEntities`
- [ ] Patch [`StockRepositoryImpl.SalesStockUpdates()`](sales/repository/stock_repository.go:133) to skip zero-delta updates
- [ ] Add/adjust unit tests for mobile no-change compare and stock delta filtering
- [ ] Add integration-style test coverage for mobile process no-change vs edit scenarios
- [ ] Validate with one reproduced mobile order scenario and confirm no extra rows in [`inv.stock`](docs/Sales%20Order%20Enhancement_BE.md)

## Mandatory backend test cases

### No-change safety
- Mobile order, full old payload resent unchanged
- Mobile order, empty detail payload process-only
- Same process request sent twice

### Real changes
- Mobile order, qty increase
- Mobile order, qty decrease
- Mobile order, add product
- Mobile order, delete product
- Mixed change set with unchanged + changed rows

### Normalization and idempotency
- `null` vs `0` on PO qty fields
- smallest-unit conversion correctness
- zero-delta update list filtered out in stock repository

## Risks and mitigations

### Risk 1
No-change branch accidentally skips legitimate detail edits.

### Mitigation
Keep compare basis narrow and explicit to PO semantic fields only.

### Risk 2
Delete logic still runs for process-only requests.

### Mitigation
Place no-change branch before `FindDetailByNotInDetailIDs` and `DeleteDetailNotInIDs` are called.

### Risk 3
Service fix works but repository still writes stock rows for zero delta.

### Mitigation
Add stock-layer zero-delta safety net.

### Risk 4
Legacy payload shape without `qty_po*` breaks compare.

### Mitigation
Use fallback normalization from `qty*` into PO compare values for mobile requests.

## Recommendation

Recommended implementation order:
1. Patch [`orderServiceImpl.Update()`](sales/service/order_service.go:3568)
2. Patch [`StockRepositoryImpl.SalesStockUpdates()`](sales/repository/stock_repository.go:133)
3. Add tests
4. Validate with one reproduced mobile order request
