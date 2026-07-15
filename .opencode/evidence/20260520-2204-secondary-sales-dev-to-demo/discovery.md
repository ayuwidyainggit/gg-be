# Discovery — Secondary Sales Report dev → demo (qa-base)

Task id: `20260520-2204-secondary-sales-dev-to-demo`
Tanggal: `2026-05-20` 22:04 Asia/Jakarta
Service target: `sales`
Source branch: `dev` (commit `8a8a0e6`)
Target base branch: `qa` (commit `e16c0a1`)
Target demo branch (baru): `demo-20052026-2204`

## Endpoint scope (dari `docs/Secondary Sales Report_BE.md`)

1. `GET /sales/v1/reports/secondary-sales/trend-sales?year=...`
2. `GET /sales/v1/reports/secondary-sales/sum-date?month=...`
3. `GET /sales/v1/reports/secondary-sales/group?month=...&group_by=...`
4. `POST /sales/v1/reports/secondary-sales`

Endpoint pendamping yang berbagi kode helper Secondary Sales Report di service `sales` (ikut tersinkron supaya compile/behavior konsisten):

- `POST /sales/v1/extract/secondary-sales` (extract pipeline + helper bersama)

## Routing kode (di `dev`)

`sales/controller/report_controller.go` (route `Route(app)`):

```go
reportRouteV1.Post("/secondary-sales", controller.SecondarySales)
reportRouteV1.Get("/secondary-sales/sum-date", controller.SecondaryReportSalesSumMonth)
reportRouteV1.Get("/secondary-sales/group", controller.SecondaryReportSalesGroup)
reportRouteV1.Get("/secondary-sales/trend-sales", controller.SecondaryReportSalesTrendSales)

reportRouteExtract.Post("/secondary-sales", controller.SecondarySalesDashboardExtract)
```

Handler-handler tersebut memanggil method service:

- `PublishSecondarySalesReport(filter)`
- `SubscribeSecondarySalesReport(filter)` (RMQ consumer untuk export)
- `SecondarySalesReportSumReportByMonth(authCustID, parentCustID, payload)`
- `SecondarySalesReportGroupSales(authCustID, parentCustID, payload)`
- `SecondarySalesReportTrendSales(authCustID, parentCustID, year, requestedCustID)`
- `ExtractReportSecondary(req)` (untuk `/v1/extract/secondary-sales` + cron)

Service memanggil repository:

- `SecondarySalesUnion`, `SecondarySalesUnionPagination`, `SecondarySalesReportSumReportByMonth`, `SecondarySalesReportReturnSumReportByMonth`, `SecondarySalesReportGroupOutlet/Salesman/ProductCategory/Product`, `SecondarySalesReportTrendSales`, `ExistsCustomerInParentScope`, plus extract helpers (`GetReportSecondarySalesReportOrder/Return`, `Save*Dim`, `SaveOrderfact`, `SaveReturnfact`, `ListCustIDReportSecondarySalesReport*`).

## Diff scope `qa..dev` (file yang menyentuh route Secondary Sales Report)

```text
M  controller/report_controller.go
M  controller/so_controller_test.go
M  entity/report.go
M  repository/report_repository.go
M  repository/report_repository_test.go
M  service/report_service.go
M  service/report_service_test.go
```

Stats:

```text
 controller/report_controller.go      | 114 +++++++-
 controller/so_controller_test.go     | 433 +++++++++++++++++++++++++++++
 entity/report.go                     |  52 ++--
 repository/report_repository.go      | 116 +++++---
 repository/report_repository_test.go | 293 ++++++++++++++++++++
 service/report_service.go            | 102 +++++--
 service/report_service_test.go       | 508 ++++++++++++++++++++++++++++++++++-
 7 files changed, 1520 insertions(+), 98 deletions(-)
```

`model/report.go` tidak berubah `qa..dev`. `pkg/config/env/env.go`, `pkg/constant/constant.go`, dan `main.go` juga tidak berubah dalam scope endpoint ini, jadi tidak perlu disalin.

## Worktree existing

```text
/Users/ujang/Projects/Geekgarden/scylla-be/sales                                 8a8a0e6 [dev]
/Users/ujang/Projects/Geekgarden/scylla-be-restore-worktrees-20260505/sales      [demo-05052026]
/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260513/sales              [demo-13052026]
/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260518/sales              [demo-18052026]
/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260518-2026/sales         [demo-18052026-2026]
```

Pola worktree path: `/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-<DDMMYYYY>[-HHMM]/sales`.
Branch demo untuk task ini belum ada:

```text
git branch --list demo-20052026-2204  # empty
```

## Konstanta target

- Branch demo baru: `demo-20052026-2204` (DDMMYYYY-HHMM dari `2026-05-20 22:04` WIB).
- Worktree path baru: `/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260520-2204/sales`.
- Source: `origin/dev` (atau `dev` lokal yang sudah sync dengan origin).
- Target base: `origin/qa`.
- Mode sinkronisasi: copy file scope endpoint dari `dev` ke `demo-20052026-2204` (tanpa cherry-pick, tanpa merge dev).

## Risiko awal

- Worktree path mungkin perlu `git fetch` dulu agar `origin/qa` dan `origin/dev` up-to-date.
- File `repository/report_repository.go` di `dev` punya helper baru (`buildReportSecondarySalesReportOrderQuery`) — copy file utuh aman karena scope endpoint sama.
- Validator package `sales/pkg/validation` di `qa` perlu mendukung tag `alphanum,max=20` (sudah ada di `validator/v10` standar; cek startup tidak bermasalah).
- `report.fact_orders.cust_id` filter di trend memerlukan data Demo. Validasi DB read-only setelah copy.
