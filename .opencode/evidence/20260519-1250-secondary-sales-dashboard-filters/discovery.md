# Discovery — Task 118/119 Secondary Sales Dashboard Filters

Task id: `20260519-1250-secondary-sales-dashboard-filters`
Waktu: `2026-05-19T12:50:36+07:00`

## File diperiksa

- `.opencode/docs/index.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `.opencode/docs/AGENT_ROUTING.md`
- `.opencode/docs/SECURITY.md`
- `.opencode/docs/PROMPT_GATES.md`
- `sales/controller/report_controller.go`
- `sales/entity/report.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository_test.go`

## Command / tool dicek

- `rtk docker compose -f docker-compose.yml ps`
  - output hanya warning compose `version` obsolete tertangkap; status service tidak terlihat pada output terpotong tool.
- `grep`/`glob` repo-local untuk endpoint, payload, repository, test, dan pola `cust_id`/`parent_cust_id`.
- `@explorer` read-only untuk pola scope BU dan test.
- `@architect` read-only untuk tenant boundary dan backward compatibility.

## Pola project ditemukan

- Repo multi-module Go; validasi harus dari direktori target service `sales`.
- Layer wajib: Controller → Service → Repository → DB.
- Multi-tenant rule: transaksi pakai `cust_id`; parent master pakai `parent_cust_id`.
- Endpoint target ada di `sales/controller/report_controller.go`:
  - `GET /secondary-sales/sum-date` → `SecondaryReportSalesSumMonth`
  - `GET /secondary-sales/group` → `SecondaryReportSalesGroup`
- Controller target memakai `c.QueryParser(&request)`, lalu service dipanggil dengan `c.Locals("cust_id").(string)`.
- Payload target di `sales/entity/report.go` saat ini hanya punya:
  - `SecondarySalesReportDashboardSumPayload.Month`
  - `SecondarySalesReportDashboardGroupPayload.Month`, `GroupBy`
- Service target hanya meneruskan `req.Month` ke repository.
- Repository target filter `report.fact_orders.cust_id = ? AND dt.month = ?`; belum filter `dt."year"`.
- `SecondarySalesReportTrendSales` sudah memakai `dt."year"` pada query raw, jadi quoted year jadi reuse pattern.

## Reuse candidates

- Reuse route, response builder, response shape, dan `QueryParser` existing.
- Reuse constant `SECONDARY_SALES_GROUP_*` untuk switch group.
- Reuse repository GORM chain style pada `report_repository.go`.
- Reuse test style:
  - `sales/repository/report_repository_test.go` memakai string SQL builder test untuk query helper.
  - `sales/service/report_service_test.go` memakai mock repository dengan function hooks.

## Gap final

1. `year` belum ada di payload, service, repository query.
2. `cust_id` request belum ada di payload.
3. Tidak ada scope validator reusable yang terbukti untuk `cust_id` query override pada report dashboard.
4. Kalau `cust_id` query langsung dipercaya, risiko horizontal tenant data leak.
5. Kalau `year` tetap opsional tanpa fallback jelas, request lama tetap bisa campur data lintas tahun.

## Keputusan yang dibutuhkan

1. `year` required atau optional/backward-compatible.
2. `cust_id` query boleh override auth `cust_id` atau hanya fallback auth sampai scope validator BU tersedia.
3. Jika override wajib, sumber scope BU yang harus dipakai belum ditemukan di `sales` path dan perlu keputusan data model/authorization.

## Risiko

- Cross-tenant leak bila `cust_id` query arbitrary diteruskan ke repository.
- Data multi-year bercampur bila `year` kosong dan query lama dipertahankan.
- Sum orders dan returns bisa tidak sinkron bila hanya salah satu query ditambah `year`.
- Branch `group_by` bisa inkonsisten bila hanya outlet yang diubah.

## Research gate

- Local project discovery: dipakai dan wajib.
- Official docs/context7: tidak dipakai; perubahan hanya repo-local Go/GORM/Fiber pattern, tidak version-sensitive.
- GitHub: tidak dipakai; tidak bergantung upstream source/issues.
- Brave/web: tidak dipakai; Google Docs referensi user sudah memberi requirement, dan akses dokumen eksternal tidak perlu untuk kontrak teknis inti.
- Browser/screenshot: tidak relevan, bukan UI visual.
