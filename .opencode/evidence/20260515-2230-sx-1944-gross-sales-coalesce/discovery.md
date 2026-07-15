# Discovery Evidence — SX-1944 Gross Sales 0 (COALESCE only)

Scope sempit: hanya bagian #1 dari Jira SX-1944, yaitu GrossSales = 0 karena ada nilai NULL pada query Secondary Sales Report. Bagian #2 (insert `amount_final` / `vat_value_final` di `mobile/v1/orders`) tidak termasuk dalam tugas ini.

## Files inspected

- `sales/controller/report_controller.go` — route `POST /v1/reports/secondary-sales` dan handler.
- `sales/service/report_service.go` — flow publish RMQ + subscribe export workbook + ExtractReportSecondary cron.
- `sales/repository/report_repository.go` — query Secondary Sales (multi-path).
- `sales/repository/report_repository_test.go` — pattern test query builder existing.
- `sales/service/report_service_test.go` — pattern test service/export.
- `mobile/service/order.go` dan `mobile/service/order_canvas.go` — hanya untuk konfirmasi bahwa fix insert order detail bukan di scope tugas ini.

## Project patterns yang relevan

- Endpoint produksi `POST /v1/reports/secondary-sales` di-route via:
  - `sales/controller/report_controller.go` → `controller.SecondarySales`
  - publish ke RMQ via `service.PublishSecondarySalesReport`
  - worker subscriber `processSecondarySalesExportMessage` → `SubscribeSecondarySalesReport`
  - export workbook membaca data via `ReportRepository.SecondarySalesUnion(dataFilter)` yang men-call `buildSecondarySalesUnionQuery(dataFilter, false)`.
- Builder `buildSecondarySalesUnionQuery` adalah single source of truth untuk paginated dan export.
- Test pattern existing untuk SQL builder ada di `sales/repository/report_repository_test.go` (TestBuildSecondarySalesUnionQuery*). Sudah aman untuk diperluas.

## Path query Secondary Sales yang menyentuh `GrossSales` (audit COALESCE)

| # | Lokasi | Dipakai endpoint mana | Status COALESCE pada kalkulasi GrossSales |
| --- | --- | --- | --- |
| 1 | `buildSecondarySalesUnionQuery` (order branch) — `report_repository.go:193-195` | Path produksi `POST /v1/reports/secondary-sales` (paginated + export) | Sudah `COALESCE(qty1_final,0) * COALESCE(sell_price1,0)` dst. Aman terhadap NULL. |
| 2 | `buildSecondarySalesUnionQuery` (return branch) — `report_repository.go:228-230` | Path produksi (return) | Sudah `COALESCE(rd.qty*,0) * COALESCE(rd.sell_price*,0)`. Aman terhadap NULL. |
| 3 | `RepositoryReportImpl.SecondarySales` — `report_repository.go:305-417` | Tidak ada caller live (legacy/dead) | NULL-prone: `qty1_final*sell_price1 + qty2_final*sell_price2 + qty3_final*sell_price3`. Tidak dipakai endpoint produksi report. |
| 4 | `GetReportSecondarySalesReportOrder` — `report_repository.go:981-1049` | Cron `ExtractReportSecondary` → tulis `report.fact_orders.gross_sale` (dashboard `sum-date` / `group` / `trend-sales`) | NULL-prone: `(od.qty1_final*od.sell_price1) + (od.qty2_final*od.sell_price2) + (od.qty3_final*od.sell_price3)` tanpa COALESCE. |
| 5 | `GetReportSecondarySalesReportReturn` — `report_repository.go:1051-...` | Cron return path | NULL-prone: `(rd.qty1*rd.sell_price1) + ...` tanpa COALESCE. |

## Konsumsi data agregat fact_orders / fact_returns

- `SecondarySalesReportSumReportByMonth` (line 1303) — `SUM(report.fact_orders.gross_sale) AS total_gross_sale`. Tanpa COALESCE.
- `SecondarySalesReportTrendSales` (line 1385) — sudah pakai `COALESCE(SUM(fo.gross_sale),0) AS total_gross_sale`. Aman.
- `SecondarySalesReportGroupOutlet/Salesman/ProductCategory/Product` — agregasi `net_sales_exclude_ppn`, bukan gross. Bukan target perbaikan #1.

Catatan: Bila row `report.fact_orders.gross_sale` dimasukkan dengan nilai NULL/0 oleh extract path NULL-prone, dashboard `sum-date` bisa menampilkan agregasi yang tidak konsisten. Untuk total non-NULL hari yang berbeda, `SUM` PostgreSQL aman terhadap NULL elemen; tetapi untuk bulan kosong, `SUM` menghasilkan NULL, yang akan di-marshal ke `0` melalui `float64` model. Keputusan: tetap tambahkan `COALESCE(..., 0)` pada agregasi `SecondarySalesReportSumReportByMonth` agar konsisten dengan `SecondarySalesReportTrendSales` dan defensif pada bulan tanpa data.

## Test patterns yang dapat di-reuse

- `TestBuildSecondarySalesUnionQueryUsesTransCTEAndDeterministicParentProductJoin` di `sales/repository/report_repository_test.go` adalah pattern verifikasi substring SQL.
- Cocok untuk menambahkan assertion bahwa `gross_sales` di kedua branch sudah memuat string `COALESCE(`.
- Tidak butuh DB integration; cukup `go test ./repository -run SecondarySales` di module `sales/`.

## Reuse candidates

- Reuse `buildSecondarySalesUnionQuery` (jangan duplikasi).
- Reuse `secondarySalesProductSelect` / `secondarySalesProductJoins` helper.
- Reuse pattern `COALESCE(..., 0)` yang sudah dipakai di builder utama (gunakan ekspresi yang sama di GetReport*) supaya behavior konsisten.

## Commands yang dipakai saat discovery

- `rtk docker compose -f docker-compose.yml ps` (output: stack belum up; tidak menghalangi karena fix ini bersifat query-level dan unit-test cukup).
- `rg`/grep: `secondary-sales`, `amount_final`, `vat_value_final`, `gross_sales`, `qty[123]_final|sell_price[123]`.

## Constraints

- Scope module: `sales/`.
- Tidak menyentuh modul `mobile/` untuk insert `amount_final`/`vat_value_final` (di luar tugas).
- Architecture rule: Controller → Service → Repository → DB tetap dijaga.
- Tidak menambah migration atau DDL.

## Risks

1. False fix risk: hanya menambah COALESCE di path #3 (legacy `SecondarySales`) yang tidak dipakai endpoint produksi. Tidak akan menyelesaikan keluhan QA pada endpoint produksi. Mitigasi: konfirmasi path live → builder utama (#1, #2) — sudah aman; perbaikan tambahan ditargetkan ke path cron extract (#4, #5) untuk dashboard; dan tetap tambahkan COALESCE di legacy path untuk konsistensi defensif.
2. Behavior regression: menambahkan `COALESCE(..., 0)` ke perkalian dapat mengubah hasil untuk row yang sebelumnya menghasilkan NULL menjadi `0` secara eksplisit. Untuk metrik gross sales ini perilaku yang diinginkan (menghindari NULL menyebar ke total dan menyebabkan tampilan `0` setelah `SUM`).
3. Risk lain di luar scope #1 (mis. `amount_final - vat_value_final`) tidak ditangani di tugas ini, sesuai pembagian kerja.

## Research gate

- Local discovery: required, completed.
- Official docs / context7: tidak diperlukan; hanya semantik SQL `COALESCE` standar.
- GitHub upstream: tidak diperlukan.
- Brave/web: tidak diperlukan.
- Browser: tidak diperlukan; bug back-end SQL.
