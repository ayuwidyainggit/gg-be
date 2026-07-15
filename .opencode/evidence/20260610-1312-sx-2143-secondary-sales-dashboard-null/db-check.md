# DB Check — SX-2143 Local `ggn_scyllax`

Tanggal: 2026-06-10 Asia/Jakarta
Scope: read-only check untuk `cust_id = C260020001`, Juni 2026.

## Code check: SX-2172 `code` enhancement

Branch lokal sudah memuat enhancement group response `code`:

- `sales/entity/report.go:263-268` memiliki `SecondarySalesReportGroupResp.Code` dengan JSON `code`.
- `sales/model/report.go:435-439` memiliki `SecondarySalesReportGroup.Code` dengan GORM column `code`.
- `sales/service/report_service.go:1451-1458` memetakan `Code: r.Code` ke response.
- `sales/repository/report_repository.go:1201-1219` dan `1231-1249` menghasilkan alias `AS code` untuk branch order/return group.
- `sales/repository/report_repository.go:1277-1279` final select/group memakai `id, code, name`.
- Regression tests sudah mencakup code mapping dan SQL alias, misalnya `sales/service/report_service_test.go:670-671` dan `sales/repository/report_repository_test.go:566, 601, 630`.

Kesimpulan: enhancement group response `code` dari SX-2172 sudah ada di code lokal.

## DB check: reporting facts

Query ke local DB berhasil: `current_database = ggn_scyllax`, `current_user = postgres`.

### `report.dim_dates` Juni 2026

Hasil:

- `dim_date_rows = 2`
- `min_day = 5`
- `max_day = 8`

Catatan: dim date lokal tidak lengkap untuk 1–30 Juni 2026; hanya ada day 5 dan 8.

### `report.fact_orders` untuk `C260020001`, Juni 2026

Hasil:

- `order_rows = 0`
- `total_net_sales_exc_ppn = 0`
- `total_gross_sale = 0`
- `total_qty = 0`
- `last_update = null`

### `report.fact_returns` untuk `C260020001`, Juni 2026

Hasil:

- `return_rows = 0`
- `total_return_net_sales_exc_ppn = 0`
- `total_return_gross_sale = 0`
- `total_return_qty = 0`
- `last_update = null`

### Mismatch cust id check di facts

Query `fo.cust_id LIKE 'C260020%'` untuk Juni 2026 menghasilkan 0 rows.

### Coverage facts all dates untuk `C260020001`

- `report.fact_orders`: 37 rows, semuanya April 2026 (`min_year=2026`, `min_month=4`, `max_year=2026`, `max_month=4`).
- `report.fact_returns`: 0 rows untuk semua tanggal.

Kesimpulan: di local `ggn_scyllax`, `report.fact_orders` dan `report.fact_returns` tidak berisi data Juni 2026 untuk `C260020001`.

## DB check: source `sls.*`

### `sls."order"` untuk `C260020001`, Juni 2026

Hasil:

- `sls_order_rows = 15`
- `min_invoice_date = 2026-06-02`
- `max_invoice_date = 2026-06-10`
- `data_status = 6` sebanyak 15 rows
- `order_detail_rows = 18` untuk orders valid `data_status IN (6,7)`

### `sls."return"` terkait order valid untuk `C260020001`, Juni 2026

Hasil:

- `sls_return_rows = 1`

### Source summary kasar dari `sls.*`

Hasil CTE source-table read-only:

- `orders = 15`
- `order_detail_rows = 18`
- `returns = 1`
- `return_detail_rows = 1`
- `order_gross = 5406000000`
- `return_gross = 650000`
- `summary_gross_sales = 5405350000`
- `summary_discount_and_promo = 1442480.0000`
- `summary_ppn = 540394626.0000`
- `summary_net_sales_exc_ppn = 5403946260`
- `summary_net_sales_inc_ppn = 5944340886`

Catatan: angka ini adalah kalkulasi read-only untuk membuktikan source data lokal ada; implementor tetap harus memastikan formula final sesuai BE docs dan kode target.

## Kesimpulan operasional

- Local facts kosong untuk Juni 2026, tetapi source `sls.*` punya order/return valid untuk `C260020001` pada Juni 2026.
- Jika endpoint `sum-date` sekarang membaca `report.fact_orders`/`report.fact_returns`, maka local API untuk Juni 2026 akan kosong meskipun source order ada.
- Untuk menutup SX-2143 di local ini, implementasi perlu salah satu:
  1. menjalankan/fix extract dashboard agar facts Juni 2026 terisi lengkap, termasuk `report.dim_dates` 1–30 Juni 2026; atau
  2. mengubah `sum-date` agar memakai source-table/date-range query `sls.*` sesuai arahan docs untuk summary order + return.
- Group endpoint yang memakai facts juga akan kosong untuk Juni 2026 sampai facts/extract tersedia, meskipun `code` enhancement sudah ada.
