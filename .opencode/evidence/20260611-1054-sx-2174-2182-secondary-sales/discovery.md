# Discovery Evidence — SX-2174 + SX-2182 Secondary Sales BE

Task id: `20260611-1054-sx-2174-2182-secondary-sales`
Tanggal: `2026-06-11`
Mode: Maintenance Stability Mode

## Source strategy

Digunakan:
- Repo-local docs: `AGENTS.md`, `.opencode/docs/index.md`, `.opencode/docs/AGENT_ROUTING.md`, `.opencode/docs/ARCHITECTURE.md`, `.opencode/docs/QUALITY.md`.
- Repo-local source dan tests di modul `sales`.
- Artifact rencana/evidence sebelumnya: `.opencode/plans/20260608-1534-sx-2182-secondary-sales-multiselect.md`, `.opencode/plans/20260605-1647-secondary-sales-defects.md`, `.opencode/plans/20260609-1442-sx-2172-secondary-sales-dashboard.md`, dan evidence terkait.
- Detail Jira/GDocs dari prompt user sebagai reference requirement.

Diskip dengan alasan:
- Official docs/context7: tidak dibutuhkan; perubahan bergantung pada kontrak Jira, SQL lokal, dan pola Go/GORM repo.
- GitHub/web search: tidak relevan karena tidak bergantung upstream package behavior.
- Browser evidence: tidak relevan untuk backend API/report; manual API/cURL runtime tetap wajib pada implementasi jika token aman tersedia.
- Jira/GDocs web fetch: tidak dilakukan karena prompt sudah memberi detail teknis cukup dan token/cookie tidak tersedia.

## Files inspected

Repo docs:
- `.opencode/docs/index.md`
- `.opencode/docs/AGENT_ROUTING.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`

Sales source:
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/entity/report.go`
- `sales/model/report.go`

Sales tests:
- `sales/service/report_service_test.go`
- `sales/repository/report_repository_test.go`

Prior artifacts:
- `.opencode/plans/20260608-1534-sx-2182-secondary-sales-multiselect.md`
- `.opencode/plans/20260605-1647-secondary-sales-defects.md`
- `.opencode/plans/20260609-1442-sx-2172-secondary-sales-dashboard.md`
- `.opencode/evidence/20260608-1534-sx-2182-secondary-sales-multiselect/discovery.md`
- `.opencode/evidence/20260605-1647-secondary-sales-defects/discovery.md`

## Project patterns found

- Repo adalah multi-module Go monorepo; validasi wajib dari target service folder.
- Repo-local command policy untuk repo ini memakai prefix `rtk`.
- Arsitektur wajib `Controller -> Service -> Repository -> DB`.
- Tenant rule: transaksi harus difilter `cust_id`, dan principal-child authorization memakai `parent_cust_id`.
- `sales` menggunakan Fiber controller, service resolver, GORM repository raw SQL, RMQ async export, dan Excel generation via `excelize`.
- Test pattern sudah kuat di `sales/service/report_service_test.go` untuk service/mock auth dan di `sales/repository/report_repository_test.go` untuk dry-run SQL assertions.

## Current implementation evidence

### SX-2182 sudah sebagian besar hadir

- `sales/entity/report.go` sudah memiliki `StringListOrScalar` dengan `UnmarshalJSON` yang menerima `null`, string, array string, trimming, comma-splitting, dedupe, dan reject non-alphanumeric `cust_id`.
- `sales/controller/report_controller.go` sudah memakai `rawSecondarySalesExportBody` dengan `RequestedCustIDRaw json.RawMessage`, lalu decode ke `secondarySalesExportBody.RequestedCustIDs`.
- Controller export membangun `entity.SecondarySalesReportQueryFilter` dengan auth fields dari JWT locals: `CustID`, `ParentCustID`, `ExportBy`; body hanya mengisi `RequestedCustID/RequestedCustIDs` dan optional filters.
- `PublishSecondarySalesReport` sudah resolve effective `CustIDs` via `resolveSecondaryDashboardCustIDs`, mempertahankan `report.list.cust_id` sebagai auth owner, dan publish `cust_ids` ke RMQ payload.
- `SubscribeSecondarySalesReport` sudah normalize fallback dari `CustIDs`, `RequestedCustIDs`, `RequestedCustID`, atau `CustID`, lalu meneruskan effective slice ke repository.
- `buildSecondarySalesUnionQuery` sudah memakai `custIDs := dataFilter.CustIDs` fallback ke `dataFilter.CustID`, lalu `od.cust_id IN ?` dan `rd.cust_id IN ?`.
- Export metadata joins sudah memakai row-level `t.cust_id` untuk outlet/salesman/distributor dan product customer join: `cp.cust_id = t.cust_id`; parent product fallback memakai `ParentCustID` sebagai parameter.
- Existing tests sudah mencakup multi-cust resolver, publish owner auth tetap, subscriber multi-cust, and repository SQL multi-cust binding.

### SX-2174 masih ada gap utama

- `SecondarySalesReportSumReportByMonth` sudah memakai source transaction tables `sls."order"`, `sls.order_detail`, `sls.return_det`, `sls."return"`, dan date range `[month start, next month)`.
- Namun `qty` masih `SUM(qty1_final + qty2_final + qty3_final)`, bukan formula satuan terkecil `(qty3 * conv2 * conv3) + (qty2 * conv2) + qty1`.
- `qty_return` masih `SUM(rd.qty1 + rd.qty2 + rd.qty3)`, bukan formula satuan terkecil dengan `conv_unit2/conv_unit3`.
- `net_sales_return` sekarang alias `rs.net_sales_exc_ppn`, padahal requirement SX-2174 meminta `Return Value` / net sales include PPN return.
- `return_rate` dihitung di service sebagai `(QtyReturn / Qty) * 100`, padahal requirement meminta `Return Rate (%)` dari value return/order: `ROUND((net_sales_return / net_sales_order) * 100, 2)` sesuai SQL acuan.
- Model `SumReportByMonthModel` belum punya field `ReturnRate`, sehingga repository belum bisa scan rate dari SQL acuan.
- `SecondarySalesReportReturnSumReportByMonth` masih memakai `report.fact_returns`; service hanya memakai `LastUpdate` dari method itu. Untuk SX-2174 jangan gunakan fact return sebagai source 4 metric merah.
- Existing targeted tests saat discovery masih pass, jadi bug belum tertutup oleh test saat ini.

## Commands checked

Dari repo root:
- `rtk docker compose -f docker-compose.yml ps`
  - Result: compose services termasuk `scylla-sales`, `scylla-master`, `rabbitmq`, `redis` up.

Dari `sales/`:
- `rtk go test ./service -run 'TestSecondarySalesReportSumReportByMonth|TestPublishSecondarySalesReport|TestSubscribeSecondarySalesReport|TestResolveSecondaryDashboardCustIDs'`
  - Result: 15 passed.
- `rtk go test ./repository -run 'TestSecondarySalesReportSumReportByMonth|TestBuildSecondarySalesUnionQuery|TestSecondarySalesReportGroup'`
  - Result: 16 passed.

## Reuse candidates

- Reuse `entity.StringListOrScalar` and `NormalizeStringList`; do not create a second custom `cust_id` parser.
- Reuse `resolveSecondaryDashboardCustIDs` for auth validation; do not move tenant decision to repository.
- Reuse `buildSecondarySalesUnionQuery` for export; fix only if manual regression shows filter all-selected still empty.
- Reuse repository dry-run SQL test helpers: `newReportRepoDryRunDB`, `latestRecordedQuery`, `assertSecondarySalesSummaryDateVars`.
- Reuse `SecondarySalesReportSumReportByMonth` as the single source for sum-date summary; add return-rate/value metrics there rather than relying on separate fact-return query.

## Constraints

- Planner wrote only `.opencode` artifacts; source implementation remains for executor.
- Do not paste or store Authorization token from Jira/cURL.
- Do not hardcode sample `cust_id`, month/year, or expected numeric values into production code.
- Keep optional empty filters as no filter; do not introduce `IN ()`/`ANY('{}')` behavior.
- Keep `report.list.cust_id` auth owner for export list visibility.
- Keep row-level `cust_id` metadata joins for export.

## Risks

- Return rate semantic changes from qty-based to value-based may break existing tests expecting 20 from qty ratio; tests must be updated to new Jira source of truth.
- SQL acuan uses `od.promo_final1..5` and `sell_price_final1..3`, while current repo query/export often uses `promo_value_final`, `sell_price1..3` or `sell_price_final1..3` depending function. Executor must inspect actual schema/current code and choose fields matching DB; if fields differ, document mapped equivalent.
- Multi-cust summary currently returns aggregate across selected `CustIDs`; SQL acuan examples use single `cust_id`. Multi-cust should aggregate safely but manual expected SX-2174 should validate single `C260020001`.
- Product/outlet/salesman optional filters are required by prompt for SX-2174 SQL, but current `sum-date` endpoint payload inspected only has `month`, `year`, `cust_id`. If endpoint lacks these query fields today, adding them changes DTO contract but is backward-compatible; executor should confirm FE sends these filters to sum-date before widening DTO.
- Manual cURL/API validation depends on safe non-committed token; planner cannot run it.

## Open questions

Tidak ada pertanyaan blocking untuk implementation-ready plan. Dua keputusan yang harus dicatat executor:
- Untuk `return_rate`, ikuti SQL acuan Jira: `return net sales include PPN / order net sales include PPN * 100`, rounded 2 decimals; bukan qty ratio.
- Untuk optional filters di `sum-date`, tambahkan DTO/support hanya jika FE/API memang mengirim `outlet_ids`, `salesman_ids`, `product_ids/pro_ids`; jika belum ada contract, repository helper bisa disiapkan tetapi controller tetap tidak menerima filter baru tanpa keputusan eksplisit.
