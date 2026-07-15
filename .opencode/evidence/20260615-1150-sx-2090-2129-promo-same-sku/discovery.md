# Discovery — SX-2090 / SX-2129 Promo Flag Same SKU

## Files inspected
- `.opencode/docs/index.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `.opencode/docs/AGENT_ROUTING.md`
- `.opencode/plans/20260528-0915-sx-2090-promo-flag.md`
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`
- `sales/entity/edit_order_enhance.go`
- `sales/entity/order_detail.go`
- `sales/model/order_detail.go`
- `sales/repository/order_repository.go`
- `sales/service/order_type_helper.go`

## Project patterns found
- Repo is multi-module Go monorepo. Target module is `sales`.
- Validation must run from `sales` directory: `rtk go test ./service -run ...`, `rtk go test ./...`.
- Layering rule: Controller → Service → Repository → DB.
- `UpdateEnhance` uses service-level transaction and `UpdateDetailPartial` for detail snapshot updates.
- DTOs already use pointer bool for enhance promo flags:
  - `EditSalesOrderDetail.IsProductPromotionSo *bool`
  - `EditFinalOrderDetail.IsProductPromotionFinal *bool`
  - `AddSalesOrderDetail.IsProductPromotionSo *bool`
  - `AddFinalOrderDetail.IsProductPromotionFinal *bool`
- `normalizeEnhancePromoFlags` maps generic `is_product_promotion` to tab-specific pointer fields and detects conflicts.
- `explicitPromoOverrides map[int64]promoFlagOverride` preserves explicit values during recompute.
- Existing tests already cover MR !95 eligibility filter and some explicit override behavior.

## Sanitized latest PM/FE feedback
- Endpoint evidence is `POST /sales/v1/orders` with same product ordered as normal detail and also promotion reward product from promo `test2090`.
- Raw prompt included `Authorization: Bearer ...`; token is intentionally not copied into this artifact or any planned command/test.
- Two sanitized create-order variants matter:
  - `order_type="O"`, `pro_id=8435`, `item_type=1`, qty3=3, expected normal row plus reward row.
  - `order_type="SO"`, `pro_id=8435`, `item_type=1`, qty3=2, expected normal row plus reward row.
- Latest expected for `SO2606100001`, `pro_id=8435`:
  - normal row `item_type=1`: all promo flags false.
  - reward row `item_type=2`: `is_product_promotion=true`, `is_product_promotion_so=true`, `is_product_promotion_final=true`, `is_product_promotion_po=false`.

## Relevant code findings
- `Store` / `prepareCreateOrderPromoState` are now in scope because latest evidence is create-order POST, not only enhance/update.
- `prepareCreateOrderPromoState` uses `aggregatePromoByProduct(consultResp)` and row aggregation by normal details before reward rows are built. This is a likely create-path source for same SKU reward flag leaking into normal row.
- `distributePromoToDetailRowsV2` skips `detail.ItemType == 2`, so reward rows are not part of normal distribution target.
- `aggregatePromoByProductForDetailSnapshot` still returns `map[int]promoAggregateRow` keyed by `pro_id`.
- Same SKU risk is at `distributePromoToDetailRowsV2` lines 1090-1094: for normal item rows, it merges `aggregateRow.Remarks` and ORs `rowAggregate.IsProductPromotion` from product-level aggregate. If `reward_product.pro_id` equals normal `pro_id`, normal row can inherit reward product flag.
- `buildDetailPromoSnapshotUpdates` writes tab-specific promo flags from aggregate keyed by detail id.
- `syncRewardProductState` deletes all `item_type=2` rows and recreates reward rows via `buildCreateOrderRewardDetails` and `buildRewardOrderDetailModels`; reward rows are already semantically separate by `item_type=2`.
- `buildCreateOrderRewardDetails` sets reward rows with `IsProductPromotionSo=true` and `IsProductPromotionFinal=true`; purchase flag is not set.
- `SyncFinalOrderFields` currently syncs qty/amount/promo_value only, not `promo_so1..5`, `promo_remarks_so`, or `is_product_promotion_so` into final fields.
- `createOrderDetailFromSalesOrder` sets `IsProductPromotionFinal` to `false` even when SO flag explicit true/false; recompute later may override, but final sync requirement says final should follow SO when no proforma.
- Model `OrderList` has `IsProformaInv *bool`, but `OrderDetailRead` does not expose `is_proforma_inv`. Repository helper `SyncFinalOrderFields(orderDetailId)` can apply SQL condition against joined/table field if needed.

## Reuse candidates
- Reuse pointer bool DTOs and `normalizeEnhancePromoFlags`.
- Reuse `promoFlagOverride` and `applyExplicitPromoFlagOverride`.
- Reuse `buildCreateOrderRewardDetails` / `buildRewardOrderDetailModels` for reward row creation.
- Reuse existing test mock patterns in `sales/service/order_service_test.go`.
- Reuse MR !95 tests:
  - `TestAggregatePromoByProductForDetailSnapshot_SX2090_RestrictRewardToEligible`
  - `TestRecomputePromoStateForTab_PreservesRowSpecificPromoFlagsForSameProduct`
  - `TestBuildDetailPromoSnapshotUpdates_SalesOrderMapsSnapshotFields`

## Constraints
- Do not copy Jira tokens, bearer tokens, passwords, or production auth into repo files, tests, logs, or docs.
- Keep write logic inside service transaction and repository tx context.
- Keep `cust_id` filters in repository writes.
- Preserve MR !95 eligibility filtering.
- Do not merge normal and reward rows by `pro_id` only.

## Risks
- Changing aggregate key broadly from `int` to struct can touch many call sites; safer fix may be distribution-level item-type awareness plus targeted helper for reward-product semantics.
- Existing helper `hasPersistedPromoSnapshot` treats `false` as no snapshot because response DTO uses non-pointer bool; do not use it as proof explicit false persisted.
- `SyncFinalOrderFields` SQL must not overwrite final fields for rows/orders with `is_proforma_inv IS NOT NULL` unless domain confirms.
- Purchase tab expected reward `is_product_promotion_po=false`; avoid copying SO reward flag into PO.

## Research gate
- Local project discovery: used and sufficient.
- Official docs/context7: skipped; behavior is repo-local Go/GORM and no version-sensitive API decision.
- GitHub/web: skipped; Jira/MR references are private and user provided enough expected behavior.
- Browser/screenshot: skipped; backend-only task.
