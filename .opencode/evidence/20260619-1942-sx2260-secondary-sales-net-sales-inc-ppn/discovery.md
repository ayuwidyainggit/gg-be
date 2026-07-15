# Discovery SX-2260

## File diperiksa
- `AGENTS.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`
- `sales/service/report_service.go`

## Pola project
- Repo multi-module Go; validasi target dari folder `sales`.
- Layer wajib: Controller → Service → Repository → DB.
- Endpoint `GET /sales/v1/reports/secondary-sales/group` route di `sales/controller/report_controller.go`.
- Service switch `group_by` ada di `sales/service/report_service.go:1434-1442`.
- Repository group query sudah memakai `buildSecondarySalesReportGroupQuery(groupBy)` dan source table `sls."order"`, `sls.order_detail`, `sls."return"`, `sls.return_det`, bukan `report.fact_orders`.

## Temuan utama
- Bug SX-2260 masih ada di group query builder: `orderNetSales` dan `returnNetSales` belum menambahkan PPN.
- `orderNetSales` saat ini: gross - promo_value_final - disc_value_final.
- `returnNetSales` saat ini: gross - promo_value - disc_value, lalu `* -1`.
- Kolom PPN tersedia:
  - order detail: `od.vat_value_final`
  - return detail: `rd.vat_value`
- Test existing sudah mengunci source-table query dan date range `invoice_date >= dateFrom AND invoice_date < dateTo`.
- Test existing belum mengunci formula include PPN untuk group query.

## Reuse candidates
- Reuse `buildSecondarySalesReportGroupQuery(groupBy)`; tidak perlu rewrite query besar.
- Reuse existing table source dan join mapping untuk `outlet`, `salesman`, `product_category`, `product`.
- Reuse dry-run SQL assertions di `sales/repository/report_repository_test.go`.
- Reuse validation command repo: `cd sales && rtk go test ./...`.

## Constraints
- Jangan ubah response contract: `id`, `code`, `name`, `net_sales`.
- Jangan hardcode token/credential.
- Jaga `cust_id IN ?`, `data_status IN (6,7)`, date range month/year.
- Jaga PPN return mengurangi total via `* -1`.
- Jangan ubah endpoint lain kecuali ditemukan helper sama dipakai endpoint group.

## Risiko
- `sell_price_final1/2/3` vs `sell_price1/2/3` berbeda antara prompt dan repo. Existing group query memakai `sell_price1/2/3`, sedangkan legacy/query lain punya varian `sell_price_final1/2/3`. Implementasi harus cek schema/test evidence sebelum mengganti nama kolom.
- Test dry-run hanya validasi SQL fragment, bukan nilai DB. Perlu manual DB/evidence compare untuk QA sheet.
- Ada `net_sales_exclude_ppn` di fungsi lain; scope SX-2260 endpoint group saja kecuali acceptance melebar.
