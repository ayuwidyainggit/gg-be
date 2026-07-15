# Quality Gate Remediation — SX-2291

Quality gate returned `PASS_WITH_RISKS` with non-blocking follow-ups. Remediated before final summary:

1. Removed dead/unreachable helper `filterCancelStockBasisForNonNeedReview`.
2. Added `activeDetailAgg` / `ad` join assertions in `TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula`.
3. Re-ran build + targeted tests inside `scylla-sales` container.

## Commands

```bash
gofmt -w sales/service/order_service.go sales/repository/stock_repository_cancel_test.go

docker exec scylla-sales sh -c "cd /app && go build ./..."

docker exec scylla-sales sh -c "cd /app && go test ./service/... -run 'TestBulkUpdateStatus_Cancel|TestValidateCancelTransition' -v"

docker exec scylla-sales sh -c "cd /app && go test ./repository/... -run TestGetCancelStockBasisQuery -v"
```

## Results

- Build: exit `0`
- Service cancel tests: exit `0`
- Repository SQL tests: exit `0`

Service test output:

```text
=== RUN   TestValidateCancelTransition
=== RUN   TestValidateCancelTransition/need_review_to_cancelled
=== RUN   TestValidateCancelTransition/processed_to_cancelled
=== RUN   TestValidateCancelTransition/completed_to_cancelled
=== RUN   TestValidateCancelTransition/already_cancelled
--- PASS: TestValidateCancelTransition (0.00s)
    --- PASS: TestValidateCancelTransition/need_review_to_cancelled (0.00s)
    --- PASS: TestValidateCancelTransition/processed_to_cancelled (0.00s)
    --- PASS: TestValidateCancelTransition/completed_to_cancelled (0.00s)
    --- PASS: TestValidateCancelTransition/already_cancelled (0.00s)
=== RUN   TestBulkUpdateStatus_Cancel_ConsistentBasisShouldApplyReversal
--- PASS: TestBulkUpdateStatus_Cancel_ConsistentBasisShouldApplyReversal (0.00s)
=== RUN   TestBulkUpdateStatus_Cancel_NeedReview_EmptyBasisShouldSkipStockWriteAndUpdateStatus
2026/06/19 09:12:20.297862 order_service.go:5126: [Warn] BulkUpdateStatus cancel skip stock reversal due empty/no-outstanding basis for need review order: SO2603220002
--- PASS: TestBulkUpdateStatus_Cancel_NeedReview_EmptyBasisShouldSkipStockWriteAndUpdateStatus (0.00s)
=== RUN   TestBulkUpdateStatus_Cancel_NeedReview_MissingSourceBasisShouldSkipStockWriteAndUpdateStatus
2026/06/19 09:12:20.297903 order_service.go:5126: [Warn] BulkUpdateStatus cancel skip stock reversal due empty/no-outstanding basis for need review order: SO2603220002B
--- PASS: TestBulkUpdateStatus_Cancel_NeedReview_MissingSourceBasisShouldSkipStockWriteAndUpdateStatus (0.00s)
=== RUN   TestBulkUpdateStatus_Cancel_MissingBasisShouldFailWithoutReversal
--- PASS: TestBulkUpdateStatus_Cancel_MissingBasisShouldFailWithoutReversal (0.00s)
=== RUN   TestBulkUpdateStatus_Cancel_AmbiguousBasisShouldFailWithoutReversal
--- PASS: TestBulkUpdateStatus_Cancel_AmbiguousBasisShouldFailWithoutReversal (0.00s)
=== RUN   TestBulkUpdateStatus_Cancel_InvalidOutstandingShouldFailWithoutReversal
--- PASS: TestBulkUpdateStatus_Cancel_InvalidOutstandingShouldFailWithoutReversal (0.00s)
PASS
ok  	sales/service	0.004s
```

Repository test output:

```text
=== RUN   TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula
--- PASS: TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula (0.00s)
=== RUN   TestGetCancelStockBasisQuery_POQtyFallback
--- PASS: TestGetCancelStockBasisQuery_POQtyFallback (0.00s)
PASS
ok  	sales/repository	0.002s
```

Remaining non-blocking risk:

- Full `go test ./...` in container still fails from unrelated existing tests:
  - missing CSV fixture: `/app/service/docs/test promo integrasi sales order - Request mas angga.csv`
  - container-local Postgres connection refused for `TestSyncFinalOrderFields_UsesNoProformaPromoSyncSQL`
