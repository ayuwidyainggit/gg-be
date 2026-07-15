# Validation — Secondary Sales dev → demo

Task id: `20260520-2204-secondary-sales-dev-to-demo`
Tanggal: `2026-05-20`
Worktree: `/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260520-2204/sales`
Branch: `demo-20052026-2204`

## Targeted tests

```bash
rtk go test ./controller/... ./service/... ./repository/... \
  -run 'TestSecondarySales|TestSecondaryReportSales|TestPublishSecondarySalesReport|TestSecondarySalesReport|TestBuildSecondarySalesUnionQuery|TestGetReportSecondarySalesReportOrder|TestSecondarySalesLegacyQuery' \
  -count=1
```

Result: `43 passed in 3 packages`

## Full suite

```bash
rtk go test ./... -count=1
```

Result: `189 passed in 22 packages`

## Route verification

```text
25: // POST /secondary-sales export endpoint. Never bind request body directly into
58: reportRouteV1.Post("/secondary-sales", controller.SecondarySales)
60: reportRouteV1.Get("/secondary-sales/sum-date", controller.SecondaryReportSalesSumMonth)
61: reportRouteV1.Get("/secondary-sales/group", controller.SecondaryReportSalesGroup)
62: reportRouteV1.Get("/secondary-sales/trend-sales", controller.SecondaryReportSalesTrendSales)
71: reportRouteExtract.Post("/secondary-sales", controller.SecondarySalesDashboardExtract)
```

All 4 docs endpoints + extract route present.
