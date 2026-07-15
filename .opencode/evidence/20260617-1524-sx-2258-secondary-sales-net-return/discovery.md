# Discovery SX-2258 Secondary Sales net return

## File/dokumen dicek

- `.opencode/docs/index.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `.opencode/docs/SERVICE_MATRIX.md`
- `.opencode/docs/PROJECT_STACK.md` tidak ada.
- `.opencode/docs/PROJECT_COMMANDS.md` tidak ada.
- `.opencode/docs/FRAMEWORK_PLAYBOOK.md` tidak ada.
- `.opencode/plans/20260511-1530-sx-1944-secondary-sales-report.md`
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`
- `sales/entity/report.go`
- `sales/model/report.go`
- `sales/service/report_service_test.go`

## Command dicek

- `rtk docker compose -f docker-compose.yml ps`
  - Hasil: command jalan, compose tidak punya container aktif.
  - Warning: `.rtk/filters.toml` belum trusted, filters tidak diterapkan.
  - Warning: `docker-compose.yml` attribute `version` obsolete.

## Route/runtime path ditemukan

- Dashboard summary card:
  - `GET /v1/reports/secondary-sales/sum-date`
  - `ReportController.SecondaryReportSalesSumMonth`
  - `ReportService.SecondarySalesReportSumReportByMonth`
  - `ReportRepository.SecondarySalesReportSumReportByMonth`
- Dashboard group:
  - `GET /v1/reports/secondary-sales/group`
  - `ReportRepository.SecondarySalesReportGroupOutlet/Salesman/ProductCategory/Product`
- Dashboard trend:
  - `GET /v1/reports/secondary-sales/trend-sales`
  - `ReportRepository.SecondarySalesReportTrendSales`
- Export/list report:
  - `POST /v1/reports/secondary-sales`
  - async RMQ to `SubscribeSecondarySalesReport`
  - `ReportRepository.SecondarySalesUnion`

## Existing formula penting

### Summary card query `sales/repository/report_repository.go:1135`

- `order_summary.qty` sudah hitung qty order dengan konversi unit.
- `return_summary.qty_return` sudah hitung qty return dengan konversi unit.
- Bug sesuai SX-2258:
  - `total_discount_promo` saat ini: `(os.discount_promo + rs.discount_promo)`.
  - response `qty` saat ini: `os.qty AS qty`.
- Formula yang harus berubah:
  - `total_discount_promo`: `COALESCE(os.discount_promo, 0) - COALESCE(rs.discount_promo, 0)`.
  - `qty`: `COALESCE(os.qty, 0) - COALESCE(rs.qty_return, 0)`.
- Query return sudah join via invoice:
  - `JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id`
  - `JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id`
- Query summary belum mendukung filter outlet/salesman/product/date query selain month/year/cust. Existing endpoint payload `SecondarySalesReportDashboardSumPayload` hanya punya `month`, `year`, `cust_id`.

### Trend sales query `buildSecondarySalesReportTrendSalesSQL`

- Bug mirip:
  - `COALESCE(os.discount_promo, 0) + COALESCE(rs.discount_promo, 0) AS total_discount_promo`.
- Return trend masih memakai `r.return_date` untuk month/date filter, bukan `o.invoice_date` seperti referensi QA.
- Tidak ada qty sold di trend response.

### Group query

- Sudah memakai `UNION ALL` order + return dengan return `net_sales * -1`.
- Filter group hanya `cust`, `month`, `year`, `group_by`; belum ada endpoint-level `outlet_ids`, `salesman_ids`, `product_ids`.

## Existing tests

- `sales/repository/report_repository_test.go` sudah punya dry-run SQL tests:
  - `TestSecondarySalesReportSumReportByMonthSQLUsesSourceTablesAndDateRange`
  - `TestSecondarySalesReportTrendSalesSQLUsesSourceTablesAndNetSalesFormula`
  - group query tests
- Test saat ini mengunci formula lama yang salah:
  - `total_discount_promo` summary memakai plus.
  - trend discount promo memakai plus.
  - summary `qty` memakai `os.qty AS qty`.
- `sales/service/report_service_test.go` punya mock repository dan service tests untuk dashboard cust scope/trend.

## Reuse candidates

- Reuse `SecondarySalesReportSumReportByMonth` sebagai titik fix utama untuk summary card.
- Reuse dry-run SQL test pattern di `sales/repository/report_repository_test.go`.
- Reuse service response mapping di `SecondarySalesReportSumReportByMonth`; bila repository alias `qty` sudah net, service tidak perlu hitung manual.

## User decisions after question gate

- Target module: `sales`.
- Trend `total_discount_promo`: ikut fixed.
- Promo/discount formula: ikuti docs/Jira reference (`disc_value_final + promo_final1..5`, bukan `promo_value_final`).

## Constraints/risks

- Repo rule: strict Controller -> Service -> Repository -> DB.
- `sales/` compose-managed service, port default `9004`.
- `pjp-sales/` punya mirror file/query serupa tetapi bukan scope SX-2258; user mengonfirmasi target module `sales`.
- Exact QA result `Number of Product Sold = 134` dan `Discount and Promo = 1.238.740` perlu staging/local dataset access. Planner tidak akses kredensial.
- Jangan ubah `total_ppn` karena Jira hanya fokus `Number of Product Sold` dan `Discount and Promo`; referensi PPN ambigu.
- `promo_final1..5` null-safe individual perlu dipakai bila executor pindah dari `promo_value_final` ke field QA reference.
