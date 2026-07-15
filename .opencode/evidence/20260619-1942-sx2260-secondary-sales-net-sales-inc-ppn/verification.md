# Verification SX-2260

## Ringkasan eksekusi
Tujuan: endpoint `GET /sales/v1/reports/secondary-sales/group` mengembalikan `net_sales` include PPN untuk semua `group_by` (outlet, salesman, product_category, product).

## Diff
- `sales/repository/report_repository.go`
  - `buildSecondarySalesReportGroupQuery`:
    - `orderNetSales`: tambah `+ COALESCE(od.vat_value_final, 0)`.
    - `returnNetSales`: tambah `+ COALESCE(rd.vat_value, 0)` sebelum `) * -1 AS net_sales`.
  - Empat method `SecondarySalesReportGroup*` tidak berubah signature; ikut builder.
- `sales/repository/report_repository_test.go`
  - `TestSecondarySalesReportGroupQueriesUseSourceTablesAndDateRange`: tambah dua fragmen check `COALESCE(od.vat_value_final, 0)` dan `COALESCE(rd.vat_value, 0)`.

## Test
- Red: `rtk go test ./repository -run TestSecondarySalesReportGroupQueriesUseSourceTablesAndDateRange -v` gagal di fragmen `vat_value_final`.
- Green: `rtk go test ./repository -run TestSecondarySalesReportGroup -v` -> 8 passed.
- Full: `cd sales && rtk go test ./...` -> 276 passed in 22 packages.
- Build: `cd sales && rtk go build ./...` -> Success.

## Bukti formula
- Order row: `((qty1_final * sell_price1) + (qty2_final * sell_price2) + (qty3_final * sell_price3)) - promo_value_final - disc_value_final + vat_value_final` -> `AS net_sales`.
- Return row: ekspresi sama dengan `vat_value` lalu `* -1` -> nilai PPN return otomatis mengurangi total.
- Outer: `SUM(net_sales)` + `COALESCE(..., 0)` + `GROUP BY id, code, name` + `ORDER BY net_sales DESC` -> tidak berubah.

## Catatan
- Pola `sell_price1/2/3` dipertahankan sesuai existing repo. Prompt doc menyebut `sell_price_final1/2/3`; jika QA evidence pakai `sell_price_final*`, swap nama kolom jadi langkah lanjut.
- `ppn_return` ikut tanda minus lewat multiplier `-1` di akhir ekspresi (sesuai II-3 plan).
- Tidak ada perubahan response contract `id/code/name/net_sales`.
- Tidak ada perubahan endpoint Secondary Sales lain.

## Status
Ready for `@quality-gate` signoff. Claim scope: hanya endpoint `secondary-sales/group` untuk `outlet/salesman/product_category/product`.
