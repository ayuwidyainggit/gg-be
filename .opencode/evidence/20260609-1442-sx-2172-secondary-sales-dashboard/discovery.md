# Discovery Evidence — SX-2172 Secondary Sales Dashboard

Task ID: `20260609-1442-sx-2172-secondary-sales-dashboard`

## Files inspected
- `AGENTS.md`
- `.opencode/docs/index.md`
- `.opencode/docs/AGENT_ROUTING.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `sales/controller/report_controller.go` disebut user sebagai route owner, tidak perlu dibaca ulang untuk query fix.
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/model/report.go`
- `sales/entity/report.go`
- `sales/repository/report_repository_test.go`
- `sales/service/report_service_test.go`

## Project patterns found
- Repo adalah multi-module Go; target modul adalah `sales`.
- Validasi dilakukan dari direktori `sales`, bukan root.
- Pola layer wajib: Controller → Service → Repository → DB.
- Dashboard group endpoint memakai service `SecondarySalesReportGroupSales` lalu repository branch sesuai `group_by`.
- `SecondarySalesReportGroup` berisi `ID`, `Code`, `Name`, `NetSales`; response `SecondarySalesReportGroupResp` mempertahankan field yang sama.
- Query group sudah memakai raw SQL builder `buildSecondarySalesReportGroupQuery(groupBy string)` dan menggabungkan `fact_orders` plus `fact_returns` via `UNION ALL`, dengan return dikurangi melalui `fr.net_sales_exclude_ppn * -1`.
- Test repository sudah memakai GORM dry-run capture SQL untuk assertion fragment SQL dan vars.
- Test service sudah punya mock branch untuk `outlet`, `salesman`, `product_category`, dan fallback `product`.

## Reuse candidates
- Reuse `buildSecondarySalesReportGroupQuery` sebagai pusat perubahan, jangan tambah query terpisah jika tidak perlu.
- Reuse `newReportRepoDryRunDB`, `latestRecordedQuery`, dan pola `strings.Contains` untuk regression test SQL.
- Reuse `SecondarySalesReportGroup` dan `SecondarySalesReportGroupResp`; tidak perlu field baru jika SQL alias tetap `id`, `code`, `name`, `net_sales`.
- Reuse master-table precedence yang sudah muncul pada query export/extract: `mst.m_product`, `mst.m_product_cat`, dan row-level `cust_id` untuk child distributor data.

## Commands/docs checked
- Dokumentasi lokal dibaca: `ARCHITECTURE.md`, `QUALITY.md`, `AGENT_ROUTING.md`.
- Tidak menjalankan test karena planner tidak mengubah source.
- Docs eksternal/GitHub/web search tidak digunakan; perilaku yang dibutuhkan adalah SQL lokal dan Jira context sudah cukup.

## Constraints
- Jangan mengubah controller/API contract kecuali mapping response terbukti tidak cocok.
- Pertahankan `cust_id IN ?`, `month`, dan quoted `dt."year"` filter.
- Jangan merusak outlet/salesman branch.
- Product/category label harus memprioritaskan master table ketika report dim kosong/tidak valid.
- Setiap perubahan source harus dilakukan oleh implementation lane, bukan planner.

## Risks
- Join ke `mst.m_product` harus memakai row-level `fact_orders.cust_id` / `fact_returns.cust_id`, bukan single auth cust, agar multi-cust dashboard tetap benar.
- Jika master product tidak punya row untuk transaksi lama, query harus fallback ke report dim agar tidak kehilangan data.
- Jika kategori master kosong, dim category boleh fallback, tetapi master tetap prioritas.
- Grouping dengan `COALESCE` harus menghindari penggabungan produk/kategori berbeda karena `id` kosong/0.
- SQL alias harus tetap sesuai `SecondarySalesReportGroup`.
