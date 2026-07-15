# Docs Reference Check — Task 118/119 Secondary Sales Dashboard Filters

Task id: `20260519-1250-secondary-sales-dashboard-filters`
Referensi: `docs/Secondary Sales Report_BE.md`

## Bagian dokumen yang relevan

- `SUM Date`, lines 59-93:
  - Method `GET`
  - URL existing `/sales/v1/reports/secondary-sales/sum-date?month=5`
  - Dokumen menulis `Request Body`, tetapi contoh dan code existing memakai query untuk `GET`.
  - Enhancement: tambah `year` dari `report.dim_dates` dan `cust_id` untuk business unit ke `report.fact_orders.cust_id`.
  - Query existing hanya filter `cust_id` dan `dt.month`; dokumen meminta tambah filter `year`.

- `Secondary Sales Group`, lines 96-123:
  - Method `GET`
  - URL existing `/sales/v1/reports/secondary-sales/group?month=5&group_by=outlet`
  - Enhancement: tambah `year` dan `cust_id` untuk business unit ke `report.fact_orders.cust_id`.
  - Response shape tetap `message`, `data`, `request_id`.

## Kesesuaian plan

Sesuai:

- Plan mempertahankan endpoint `GET` existing.
- Plan memakai query param, bukan request body, karena `GET` + controller existing memakai `QueryParser`.
- Plan menambahkan `year` dan `cust_id` pada payload query.
- Plan menambahkan `dt."year" = ?` pada repository untuk sum-date dan semua branch group.
- Plan memakai effective `cust_id` untuk `report.fact_orders.cust_id`; untuk return summary juga memakai `report.fact_returns.cust_id` supaya summary return tidak mismatch.
- Plan menjaga response shape existing.
- Plan menambahkan scope check karena docs meminta business unit, sementara repo punya tenant isolation rule.

Batas scope:

- Dokumen global menyebut filter wilayah `Region`, `Area`, `Distributor` pada Homepage dan Modal Export.
- Task user saat ini hanya Task 118 `SUM Date` dan Task 119 `Secondary Sales Group`, dengan expected filter `year` dan `cust_id`.
- Plan tidak mencakup Trend Sales dan Export Secondary Sales.
- Plan tidak menambah filter `region`/`area` karena bagian endpoint Task 118/119 pada dokumen hanya menyebut `year` dan `cust_id`.

## Penyesuaian yang perlu masuk plan

- Tambahkan referensi eksplisit ke `docs/Secondary Sales Report_BE.md` pada `Existing Patterns/Reuse` atau `Evidence Requirements`.
- Tambahkan constraint bahwa `region`/`area`/`distributor` dan export endpoint di dokumen berada di luar scope Task 118/119 kecuali user memperluas scope.
- Tambahkan note bahwa `cust_id` query memenuhi kebutuhan Distributor/BU selection, tetapi Region/Area butuh task terpisah bila diwajibkan.
